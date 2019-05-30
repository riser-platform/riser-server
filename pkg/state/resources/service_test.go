package resources

import (
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateService(t *testing.T) {
	deployment := &core.Deployment{
		DeploymentMeta: core.DeploymentMeta{
			Name:      "myapp-deployment",
			Namespace: "apps",
			Stage:     "dev",
		},
		App: &model.AppConfig{
			Name: "myapp",
			Expose: &model.AppConfigExpose{
				ContainerPort: 8000,
			},
		},
	}

	result, err := CreateService(deployment)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, deployment.Name, result.Name)
	assert.Equal(t, "apps", result.Namespace)

	assert.Equal(t, 4, len(result.Labels))
	assert.Equal(t, "dev", result.Labels["stage"])
	assert.Equal(t, "myapp", result.Labels["app"])
	assert.Equal(t, "myapp-deployment", result.Labels["deployment"])
	assert.Equal(t, defaultRiserAppVersion, result.Labels["riser-app"])

	assert.Equal(t, "Service", result.Kind)
	assert.Equal(t, "v1", result.APIVersion)

	assert.Equal(t, 1, len(result.Spec.Ports))
	assert.EqualValues(t, 8000, result.Spec.Ports[0].Port)

	assert.Equal(t, result.Spec.Selector, map[string]string{
		"deployment": "myapp-deployment",
	})
}
