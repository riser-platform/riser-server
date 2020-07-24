package model

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: See rules_test for more test coverage

// Define the minimum valid config here. Don't reference it directly though. Use createMinAppConfig instead to get a deep clone
var minimumValidAppConfig = &AppConfig{
	Name:      "myapp",
	Namespace: "myns",
	Id:        uuid.New(),
	Image:     "myimage",
	Expose: &AppConfigExpose{
		ContainerPort: 80,
	},
}

func Test_AppConfig_ApplyDefaults(t *testing.T) {
	appConfig := &AppConfig{
		Name: "myapp",
	}

	err := appConfig.ApplyDefaults()

	assert.NoError(t, err)
	assert.EqualValues(t, "myapp", appConfig.Name)
	assert.EqualValues(t, "apps", appConfig.Namespace)
	assert.NotNil(t, appConfig.Expose)
	assert.Equal(t, "http", appConfig.Expose.Protocol)
	assert.Equal(t, "external", appConfig.Expose.Scope)
}

func Test_AppConfig_ApplyDefaults_AllowsNonDefaultValues(t *testing.T) {
	appConfig := &AppConfig{
		Name:      "myapp",
		Namespace: "myns",
		Expose: &AppConfigExpose{
			ContainerPort: 8000,
			Protocol:      "http2",
			Scope:         AppExposeScope_Cluster,
		},
	}

	err := appConfig.ApplyDefaults()

	assert.NoError(t, err)
	assert.EqualValues(t, "myapp", appConfig.Name)
	assert.EqualValues(t, "myns", appConfig.Namespace)
	assert.Equal(t, "http2", appConfig.Expose.Protocol)
	assert.Equal(t, AppExposeScope_Cluster, appConfig.Expose.Scope)
	assert.EqualValues(t, 8000, appConfig.Expose.ContainerPort)
}

func Test_AppConfig_ValidateName(t *testing.T) {
	appConfig := createMinAppConfig()
	appConfig.Name = "5name"
	appConfig.Namespace = "5ns"

	err := appConfig.Validate()
	assert.IsType(t, validation.Errors{}, err)
	validationErrors := err.(validation.Errors)
	require.Len(t, validationErrors, 2)
	assert.Equal(t, "must be lowercase, alphanumeric, and start with a letter", validationErrors["name"].Error())
	assert.Equal(t, "must be lowercase, alphanumeric, and start with a letter", validationErrors["namespace"].Error())
}

func Test_AppConfig_ValidateRequired(t *testing.T) {
	appConfig := AppConfig{}

	err := appConfig.Validate()

	assert.IsType(t, validation.Errors{}, err)
	validationErrors := err.(validation.Errors)
	assert.Len(t, validationErrors, 5)
	assertFieldsRequired(t, validationErrors, "name", "namespace", "id", "image", "expose")
}

func Test_AppConfig_ValidateExposeRequired(t *testing.T) {
	appConfig := createMinAppConfig()
	appConfig.Expose.ContainerPort = 0
	err := appConfig.Validate()

	assert.IsType(t, validation.Errors{}, err)
	validationErrors := err.(validation.Errors)
	assert.Len(t, validationErrors, 1)
}

func Test_AppConfig_ValidateAutoscaleRange(t *testing.T) {
	min := -1
	max := 0
	appConfig := createMinAppConfig()
	appConfig.Autoscale = &AppConfigAutoscale{
		Min: &min,
		Max: &max,
	}
	err := appConfig.Validate()

	assert.IsType(t, validation.Errors{}, err)
	validationErrors := err.(validation.Errors)
	assert.Len(t, validationErrors, 2)
	require.Contains(t, validationErrors, "autoscale.min", validationErrors)
	require.Contains(t, validationErrors, "autoscale.max", validationErrors)
	assert.Equal(t, "must be no less than 0", validationErrors["autoscale.min"].Error())
	assert.Equal(t, "must be no less than 1", validationErrors["autoscale.max"].Error())
}

func Test_AppConfig_ValidateAutoscaleMaxGtMin(t *testing.T) {
	min := 2
	max := 1
	appConfig := createMinAppConfig()
	appConfig.Autoscale = &AppConfigAutoscale{
		Min: &min,
		Max: &max,
	}
	err := appConfig.Validate()

	assert.IsType(t, validation.Errors{}, err)
	validationErrors := err.(validation.Errors)
	assert.Len(t, validationErrors, 1)
	require.Contains(t, validationErrors, "autoscale.max", validationErrors)
	assert.Equal(t, "must be greater than or equal to autoscale.min", validationErrors["autoscale.max"].Error())
}

func Test_AppConfig_ValidateAutoscaleMax_NilMin(t *testing.T) {
	max := 0
	appConfig := createMinAppConfig()
	appConfig.Autoscale = &AppConfigAutoscale{
		Max: &max,
	}
	err := appConfig.Validate()

	assert.IsType(t, validation.Errors{}, err)
	validationErrors := err.(validation.Errors)
	assert.Len(t, validationErrors, 1)
	require.Contains(t, validationErrors, "autoscale.max", validationErrors)
	assert.Equal(t, "must be no less than 1", validationErrors["autoscale.max"].Error())
}

// Note: We may not allow registry to be set here - it may be dictated by an admin on a per environment basis instead.
var imageTests = []struct {
	image string
	valid bool
}{
	{"image", true},
	{"my/image", true},
	{"registry.io/image", true},
	{"registry.io/my/image", true},
	{"image:tag", false},
	{"image@sha256:b876dd4c32a96067ab22201e521d4fe3724f6e5af7d48f50b0059ae253359a4c", false},
	{"my/image:tag", false},
	{"registry.io/my/image:tag", false},
}

func Test_AppConfig_ValidateImage(t *testing.T) {
	for _, tt := range imageTests {
		appConfig := createMinAppConfig()
		appConfig.Image = tt.image
		err := appConfig.Validate()

		if tt.valid {
			assert.NoError(t, err, tt.image)
		} else {
			require.IsType(t, validation.Errors{}, err, tt.image)
			validationErrors := err.(validation.Errors)
			assert.Len(t, validationErrors, 1, tt.image)
			require.Contains(t, validationErrors, "image", tt.image)
			assert.Equal(t, "must not contain a tag or digest", validationErrors["image"].Error(), tt.image)
		}
	}
}

var protocolTests = []struct {
	protocol string
	valid    bool
}{
	{"http", true},
	{"http2", true},
	{"", true},
	{"redis", false},
}

func Test_AppConfig_ValidateExposeProtocol(t *testing.T) {
	for _, tt := range protocolTests {
		appConfig := createMinAppConfig()
		appConfig.Expose.Protocol = tt.protocol
		err := appConfig.Validate()

		if tt.valid {
			assert.NoError(t, err, tt.protocol)
		} else {
			require.IsType(t, validation.Errors{}, err, tt.protocol)
			validationErrors := err.(validation.Errors)
			assert.Len(t, validationErrors, 1, tt.protocol)
			require.Contains(t, validationErrors, "expose.protocol", tt.protocol)
			assert.Equal(t, "must be one of: http, http2", validationErrors["expose.protocol"].Error(), tt.protocol)
		}
	}
}

func Test_ApplyOverrides_NoOverrides(t *testing.T) {
	appConfig := &AppConfigWithOverrides{
		AppConfig: AppConfig{
			Name: "myapp",
		},
	}

	result, err := appConfig.ApplyOverrides("test")

	require.NoError(t, err)
	assert.EqualValues(t, "myapp", result.Name)
}

func Test_ApplyOverrides_NoOverridesForEnvironment(t *testing.T) {
	cpuCores := float32(2)
	cpuCoresDev := float32(0.1)
	appConfig := &AppConfigWithOverrides{
		AppConfig: AppConfig{
			Name: "myapp",
			OverrideableAppConfig: OverrideableAppConfig{
				Resources: &AppConfigResources{
					CpuCores: &cpuCores,
				},
			},
		},
		Overrides: map[string]OverrideableAppConfig{
			"dev": {
				Resources: &AppConfigResources{
					CpuCores: &cpuCoresDev,
				},
			},
		},
	}

	result, err := appConfig.ApplyOverrides("test")

	require.NoError(t, err)
	assert.EqualValues(t, "myapp", result.Name)
	assert.Equal(t, cpuCores, *result.Resources.CpuCores)
}

func Test_ApplyOverrides_WithOverrides(t *testing.T) {
	appId := uuid.MustParse("AEEA1A7A-70FE-4B3A-8436-3F8A197279DC")
	cpuCores := float32(2)
	cpuCoresDev := float32(0.1)
	autoscaleMin := int(1)
	// Setting the override to zero covers an important case where mergo does not override 0 by default even for an *int
	autoscaleMinOverride := int(0)
	appConfig := &AppConfigWithOverrides{
		AppConfig: AppConfig{
			Id:    appId,
			Name:  "myapp",
			Image: "hashicorp/http-echo",
			HealthCheck: &AppConfigHealthCheck{
				Path: "/health",
			},
			Expose: &AppConfigExpose{
				ContainerPort: 1337,
			},
			OverrideableAppConfig: OverrideableAppConfig{
				Autoscale: &AppConfigAutoscale{
					Min: &autoscaleMin,
				},
				Resources: &AppConfigResources{
					CpuCores: &cpuCores,
				},
				Environment: map[string]intstr.IntOrString{
					"envKey":     intstr.Parse("envVal"),
					"envKeyBase": intstr.Parse("envValBase"),
				},
			},
		},
		Overrides: map[string]OverrideableAppConfig{
			"dev": {
				Autoscale: &AppConfigAutoscale{
					Min: &autoscaleMinOverride,
				},
				Resources: &AppConfigResources{
					CpuCores: &cpuCoresDev,
				},
				Environment: map[string]intstr.IntOrString{
					"envKey":    intstr.Parse("envValDevOverride"),
					"envKeyDev": intstr.Parse("envValDev"),
				},
			},
		},
	}

	result, err := appConfig.ApplyOverrides("dev")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, appId, result.Id)
	assert.EqualValues(t, "myapp", result.Name)
	assert.Len(t, result.Environment, 3)
	assert.Equal(t, "envValDevOverride", result.Environment["envKey"].StrVal)
	assert.Equal(t, "envValDev", result.Environment["envKeyDev"].StrVal)
	assert.Equal(t, "envValBase", result.Environment["envKeyBase"].StrVal)
	assert.EqualValues(t, cpuCoresDev, *result.Resources.CpuCores)
	assert.Equal(t, "/health", result.HealthCheck.Path)
	assert.Equal(t, 0, *result.Autoscale.Min)
	// Ensure that we don't mutate the original config
	assert.Equal(t, appConfig.Resources.CpuCores, &cpuCores)
}

func Test_AppConfig_ValidateExposeScope(t *testing.T) {
	var tests = []struct {
		scope string
		valid bool
	}{
		{AppExposeScope_External, true},
		{AppExposeScope_Cluster, true},
		{"", true},
		{"nope", false},
	}

	for _, tt := range tests {
		appConfig := createMinAppConfig()
		appConfig.Expose.Scope = tt.scope
		err := appConfig.Validate()

		if tt.valid {
			assert.NoError(t, err, tt.scope)
		} else {
			require.IsType(t, validation.Errors{}, err, tt.scope)
			validationErrors := err.(validation.Errors)
			assert.Len(t, validationErrors, 1, tt.scope)
			require.Contains(t, validationErrors, "expose.scope", tt.scope)
			assert.Equal(t, "must be one of: external, cluster", validationErrors["expose.scope"].Error(), tt.scope)
		}
	}
}

func createMinAppConfig() *AppConfig {
	appConfig := &AppConfig{}
	_ = copier.Copy(appConfig, minimumValidAppConfig)
	return appConfig
}
