package model

import (
	"github.com/docker/distribution/reference"
	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/google/uuid"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	DefaultAppConfig = &AppConfig{
		Namespace: "apps",
		Expose: &AppConfigExpose{
			Protocol: "http",
		},
	}
)

// TODO: Move outside the API and into a separate module. The AppConfig should version independently of the API via a Version field on the root AppConfig object
// Also, pkg/* should not have a dependency here. However, moving this into pkg/core (for example) would cause a circular module dependency so we
// may need to create a separate module e.g. pkg/core/appconfig

// AppConfigWithOverrides contains an app with environment level overrides
type AppConfigWithOverrides struct {
	AppConfig
	Overrides map[string]OverrideableAppConfig `json:"environmentOverrides,omitempty"`
}

func (cfg *AppConfigWithOverrides) ApplyOverrides(envName string) (*AppConfig, error) {
	app := cfg.AppConfig
	if overrideApp, ok := cfg.Overrides[envName]; ok {
		err := mergo.Merge(&overrideApp, app.OverrideableAppConfig)
		if err != nil {
			return nil, err
		}

		app.OverrideableAppConfig = overrideApp
	}

	return &app, nil

}

// AppConfig is the root of the application config object graph without environment overrides
type AppConfig struct {
	Id                    uuid.UUID             `json:"id"`
	Name                  AppName               `json:"name"`
	Namespace             NamespaceName         `json:"namespace"`
	Expose                *AppConfigExpose      `json:"expose,omitempty"`
	HealthCheck           *AppConfigHealthCheck `json:"healthcheck,omitempty"`
	Image                 string                `json:"image"`
	OverrideableAppConfig `json:",inline"`
}

// OverrideableAppConfig contains properties that are overrideable
type OverrideableAppConfig struct {
	Autoscale   *AppConfigAutoscale           `json:"autoscale,omitempty"`
	Environment map[string]intstr.IntOrString `json:"environment,omitempty"`
	Resources   *AppConfigResources           `json:"resources,omitempty"`
}

type AppConfigAutoscale struct {
	Min *int `json:"min,omitempty"`
	Max *int `json:"max,omitempty"`
}

type AppConfigExpose struct {
	Protocol      string `json:"protocol,omitempty"`
	ContainerPort int32  `json:"containerPort"`
}

// Mode is not yet implemented (httpGet = default)
type AppConfigHealthCheck struct {
	Path string `json:"path,omitempty"`
}

type AppConfigResources struct {
	CpuCores *float32 `json:"cpuCores,omitempty"`
	MemoryMB *int32   `json:"memoryMB,omitempty"`
}

// ApplyDefaults sets any unset values with their defaults
func (appConfig *AppConfig) ApplyDefaults() error {
	return mergo.Merge(appConfig, DefaultAppConfig)
}

func (appConfig AppConfig) Validate() error {
	err := validation.ValidateStruct(&appConfig,
		validation.Field(&appConfig.Name),
		validation.Field(&appConfig.Namespace),
		validation.Field(&appConfig.Id, validation.By(validId)),
		validation.Field(&appConfig.Image, validation.Required, validation.By(validDockerImageWithoutTagOrDigest)),
		validation.Field(&appConfig.Expose, validation.Required),
	)

	// Break out each struct so that we can have better error messages than the default
	// This has the downside of not allowing nested structs to implement their own Validate.

	if appConfig.Expose != nil {
		exposeErr := validation.ValidateStruct(appConfig.Expose,
			validation.Field(&appConfig.Expose.Protocol, validation.In("http", "http2").Error("must be one of: http, http2")),
			validation.Field(&appConfig.Expose.ContainerPort, validation.Required, validation.Min(1), validation.Max(65535)),
		)
		err = mergeValidationErrors(err, exposeErr, "expose")
	}

	if appConfig.Autoscale != nil {
		maxMinRule := validation.Min(1)
		if appConfig.Autoscale.Min != nil {
			maxMinRule = validation.Min(*appConfig.Autoscale.Min).Error("must be greater than or equal to autoscale.min")
		}
		autoscaleErr := validation.ValidateStruct(appConfig.Autoscale,
			validation.Field(&appConfig.Autoscale.Min, validation.Min(0)),
			// We have to customize the NilOrEmpty error to match "Min since "Min" does not get applied to nillable 0 value
			validation.Field(&appConfig.Autoscale.Max, validation.NilOrNotEmpty.Error("must be no less than 1"), maxMinRule),
		)

		err = mergeValidationErrors(err, autoscaleErr, "autoscale")
	}

	return err
}

func validDockerImageWithoutTagOrDigest(value interface{}) error {
	dockerImageURL, _ := value.(string)
	named, err := reference.ParseNormalizedNamed(dockerImageURL)
	if err != nil {
		return errors.Wrap(err, "must be a valid docker image url")
	}

	if !reference.IsNameOnly(named) {
		return errors.New("must not contain a tag or digest")
	}

	return nil
}

func validId(v interface{}) error {
	id, _ := v.(uuid.UUID)
	if id == uuid.Nil {

		return errors.New("cannot be blank")
	}
	return nil
}
