package model

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/jinzhu/copier"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var minimumValidAppConfig = &AppConfig{
	Name:  "myapp",
	Id:    "myid",
	Image: "myimage",
}

func Test_AppConfig_ValidateRequired(t *testing.T) {
	appConfig := AppConfig{}
	err := appConfig.Validate()

	assert.IsType(t, validation.Errors{}, err)
	validationErrors := err.(validation.Errors)
	assert.Len(t, validationErrors, 3)
	assertFieldsRequired(t, validationErrors, "name", "id", "image")
}

// Note: We may not allow registry to be set here - it may be dictated by an admin on a per stage basis instead.
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
		appConfig := &AppConfig{}
		_ = copier.Copy(appConfig, minimumValidAppConfig)
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
	{"grpc", true},
	{"", true},
	{"redis", false},
}

func Test_AppConfig_ValidateProtocol(t *testing.T) {
	for _, tt := range protocolTests {
		appConfig := &AppConfig{}
		_ = copier.Copy(appConfig, minimumValidAppConfig)
		appConfig.Expose = &AppConfigExpose{Protocol: tt.protocol}
		err := appConfig.Validate()

		if tt.valid {
			assert.NoError(t, err, tt.protocol)
		} else {
			require.IsType(t, validation.Errors{}, err, tt.protocol)
			validationErrors := err.(validation.Errors)
			assert.Len(t, validationErrors, 1, tt.protocol)
			require.Contains(t, validationErrors, "expose.protocol", tt.protocol)
			assert.Equal(t, "must be one of: http, grpc", validationErrors["expose.protocol"].Error(), tt.protocol)
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
	assert.Equal(t, "myapp", result.Name)
}

func Test_ApplyOverrides_NoOverridesForStage(t *testing.T) {
	cpuCores := float32(2)
	cpuCoresDev := float32(0.1)
	appConfig := &AppConfigWithOverrides{
		AppConfig: AppConfig{
			Name: "myapp",
			Resources: &AppConfigResources{
				CpuCores: &cpuCores,
			},
		},
		Overrides: map[string]AppConfig{
			"dev": AppConfig{
				Resources: &AppConfigResources{
					CpuCores: &cpuCoresDev,
				},
			},
		},
	}

	result, err := appConfig.ApplyOverrides("test")

	require.NoError(t, err)
	assert.Equal(t, "myapp", result.Name)
	assert.Equal(t, cpuCores, *result.Resources.CpuCores)
}

func Test_ApplyOverrides_WithOverrides(t *testing.T) {
	cpuCores := float32(2)
	cpuCoresDev := float32(0.1)
	appConfig := &AppConfigWithOverrides{
		AppConfig: AppConfig{
			Name:  "myapp",
			Image: "hashicorp/http-echo",
			Resources: &AppConfigResources{
				CpuCores: &cpuCores,
			},
			HealthCheck: &AppConfigHealthCheck{
				Path: "/health",
			},
			Environment: map[string]intstr.IntOrString{
				"envKey":     intstr.Parse("envVal"),
				"envKeyBase": intstr.Parse("envValBase"),
			},
			Expose: &AppConfigExpose{
				ContainerPort: 1337,
			},
		},
		Overrides: map[string]AppConfig{
			"dev": AppConfig{
				Resources: &AppConfigResources{
					CpuCores: &cpuCoresDev,
				},
				Environment: map[string]intstr.IntOrString{
					"envKey":    intstr.Parse("envValDevOverride"),
					"envKeyDev": intstr.Parse("envValDev"),
				},
				Expose: &AppConfigExpose{
					ContainerPort: 8080,
				},
			},
		},
	}

	result, err := appConfig.ApplyOverrides("dev")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "myapp", result.Name)
	assert.Len(t, result.Environment, 3)
	assert.Equal(t, "envValDevOverride", result.Environment["envKey"].StrVal)
	assert.Equal(t, "envValDev", result.Environment["envKeyDev"].StrVal)
	assert.Equal(t, "envValBase", result.Environment["envKeyBase"].StrVal)
	assert.EqualValues(t, 8080, result.Expose.ContainerPort)
	assert.EqualValues(t, cpuCoresDev, *result.Resources.CpuCores)
	assert.Equal(t, "/health", result.HealthCheck.Path)
}

func assertFieldsRequired(t *testing.T, errors validation.Errors, fieldNames ...string) {
	for _, fieldName := range fieldNames {
		require.Contains(t, errors, fieldName, "missing required field %q", fieldName)
		assert.Equal(t, "cannot be blank", errors[fieldName].Error())
	}
}
