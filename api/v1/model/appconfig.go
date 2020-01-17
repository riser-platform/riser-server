package model

import (
	"github.com/docker/distribution/reference"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// TODO: Move outside the API and into a separate module. The AppConfig should version independently of the API via a Version field on the root AppConfig object
// Also, pkg/* should not have a dependency here. However, moving this into pkg/core (for example) would cause a circular module dependency so we
// may need to create a separate module e.g. pkg/core/appconfig

// AppConfigWithOverrides contains an app with stage level overrides
type AppConfigWithOverrides struct {
	AppConfig
	Overrides map[string]AppConfig `json:"stages,omitempty"`
}

func (cfg *AppConfigWithOverrides) ApplyOverrides(stageName string) (*AppConfig, error) {
	app := &cfg.AppConfig
	if overrideApp, ok := cfg.Overrides[stageName]; ok {
		err := mergo.Merge(&overrideApp, app)
		if err != nil {
			return nil, err
		}

		app = &overrideApp
	}

	return app, nil

}

// AppConfig is the root of the application config object graph without stage overrides
type AppConfig struct {
	Name        string                        `json:"name"`
	Environment map[string]intstr.IntOrString `json:"environment,omitempty"`
	Expose      *AppConfigExpose              `json:"expose,omitempty"`
	HealthCheck *AppConfigHealthCheck         `json:"healthcheck,omitempty"`
	// Id is a random id used to prevent collisions (two apps with the same name and namespace)
	Id        string              `json:"id"`
	Image     string              `json:"image"`
	Resources *AppConfigResources `json:"resources,omitempty"`
}

type AppConfigExpose struct {
	Protocol      string `json:"protocol,omitempty"`
	ContainerPort int32  `json:"containerPort"`
}

// Mode is not yet implemented (httpGet = default)
type AppConfigHealthCheck struct {
	Path string `json:"path,omitempty"`
	Port *int32 `json:"port,omitempty"`
}

type AppConfigResources struct {
	CpuCores *float32 `json:"cpuCores,omitempty"`
	MemoryMB *int32   `json:"memoryMB,omitempty"`
}

func (appConfig *AppConfig) Validate() error {
	err := validation.ValidateStruct(appConfig,
		validation.Field(&appConfig.Name, RulesAppName()...),
		validation.Field(&appConfig.Id, validation.Required),
		validation.Field(&appConfig.Image, validation.Required, validation.By(validDockerImageWithoutTagOrDigest)),
	)

	if appConfig.Expose != nil {
		exposeError := validation.ValidateStruct(appConfig.Expose,
			validation.Field(&appConfig.Expose.Protocol, validation.In("http", "grpc").Error("must be one of: http, grpc")))
		err = mergeValidationErrors(err, exposeError, "expose")
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
