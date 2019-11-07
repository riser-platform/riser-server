package resources

import (
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateVirtualService(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:      "myapp-deployment",
		Namespace: "apps",
		Stage:     "dev",
		App: &model.AppConfig{
			Name: "myapp",
			Expose: &model.AppConfigExpose{
				ContainerPort: 8000,
			},
		},
	}
	stage := &core.StageConfig{
		PublicGatewayHost: "dev.riser.org",
	}

	result, err := CreateVirtualService(&core.DeploymentContext{Deployment: deployment, Stage: stage, RiserGeneration: 3})

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "apps-myapp-deployment-8000", result.Name)
	assert.Equal(t, "apps", result.Namespace)

	assert.Equal(t, 3, len(result.Labels))
	assert.Equal(t, "dev", result.Labels[riserLabel("stage")])
	assert.Equal(t, "myapp", result.Labels[riserLabel("app")])
	assert.Equal(t, "myapp-deployment", result.Labels[riserLabel("deployment")])
	assert.Len(t, result.Annotations, 1)
	assert.Equal(t, "3", result.Annotations[riserLabel("generation")])

	assert.Equal(t, 2, len(result.Spec.Gateways))
	assert.Equal(t, defaultGateway, result.Spec.Gateways[0])
	assert.Equal(t, "mesh", result.Spec.Gateways[1])

	assert.Equal(t, 2, len(result.Spec.Hosts))
	assert.Equal(t, "myapp-deployment.apps.dev.riser.org", result.Spec.Hosts[0])
	assert.Equal(t, "myapp-deployment.apps.svc.cluster.local", result.Spec.Hosts[1])
}
