package deployment

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/stretchr/testify/require"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_Sanitize_DefaultsToAppName(t *testing.T) {
	deployment := &core.NewDeployment{
		App: &model.AppConfigWithOverrides{
			AppConfig: model.AppConfig{
				Name: "myapp",
			},
		},
	}

	result := Sanitize(deployment)

	assert.Equal(t, "myapp", result.Name)
}

func Test_Sanitize_AddsAppNameAsPrefix(t *testing.T) {
	deployment := &core.NewDeployment{
		DeploymentMeta: core.DeploymentMeta{
			Name: "mydeployment",
		},
		App: &model.AppConfigWithOverrides{
			AppConfig: model.AppConfig{
				Name: "myapp",
			},
		},
	}

	result := Sanitize(deployment)

	assert.Equal(t, "myapp-mydeployment", result.Name)
}

func Test_Sanitize_DoesNotAddPrefixIfNamesMatch(t *testing.T) {
	deployment := &core.NewDeployment{
		DeploymentMeta: core.DeploymentMeta{
			Name: "myapp",
		},
		App: &model.AppConfigWithOverrides{
			AppConfig: model.AppConfig{
				Name: "myapp",
			},
		},
	}

	result := Sanitize(deployment)

	assert.Equal(t, "myapp", result.Name)
}

func Test_Sanitize_DefaultNamespace(t *testing.T) {
	deployment := &core.NewDeployment{
		App: &model.AppConfigWithOverrides{},
	}

	result := Sanitize(deployment)

	assert.Equal(t, "apps", result.Namespace)
}

func Test_ApplyOverrides_NoOverrides(t *testing.T) {
	appConfig := model.AppConfig{
		Name: "myapp",
	}
	deployment := &core.NewDeployment{
		DeploymentMeta: core.DeploymentMeta{
			Name:  "myapp",
			Stage: "dev",
		},
		App: &model.AppConfigWithOverrides{
			AppConfig: appConfig,
		},
	}

	result, err := ApplyOverrides(deployment)

	require.NoError(t, err)
	assert.Equal(t, appConfig, *result.App)
}

func Test_ApplyOverrides_WithOverrides(t *testing.T) {
	replicas := int32(3)
	replicasDev := int32(1)
	deployment := &core.NewDeployment{
		DeploymentMeta: core.DeploymentMeta{
			Name:  "myapp-deployment",
			Stage: "dev",
			Docker: core.DeploymentDocker{
				Tag: "myTag",
			},
		},
		App: &model.AppConfigWithOverrides{
			AppConfig: model.AppConfig{
				Name:     "myapp",
				Image:    "hashicorp/http-echo",
				Replicas: &replicas,
				HealthCheck: &model.AppConfigHealthCheck{
					Path: "/health",
				},
				Environment: map[string]intstr.IntOrString{
					"envKey":     intstr.Parse("envVal"),
					"envKeyBase": intstr.Parse("envValBase"),
				},
				Expose: &model.AppConfigExpose{
					ContainerPort: 1337,
				},
			},
			Overrides: map[string]model.AppConfig{
				"dev": model.AppConfig{
					Replicas: &replicasDev,
					Environment: map[string]intstr.IntOrString{
						"envKey":    intstr.Parse("envValDevOverride"),
						"envKeyDev": intstr.Parse("envValDev"),
					},
					Expose: &model.AppConfigExpose{
						ContainerPort: 8080,
					},
				},
			},
		},
	}

	result, err := ApplyOverrides(deployment)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "myapp-deployment", result.Name)
	assert.Equal(t, "myapp", result.App.Name)
	assert.Len(t, result.App.Environment, 3)
	assert.Equal(t, "envValDevOverride", result.App.Environment["envKey"].StrVal)
	assert.Equal(t, "envValDev", result.App.Environment["envKeyDev"].StrVal)
	assert.Equal(t, "envValBase", result.App.Environment["envKeyBase"].StrVal)
	assert.EqualValues(t, 8080, result.App.Expose.ContainerPort)
	assert.EqualValues(t, 1, *result.App.Replicas)
	assert.Equal(t, "/health", result.App.HealthCheck.Path)
}
