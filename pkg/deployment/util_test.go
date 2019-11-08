package deployment

import (
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_applyDefaults_Defaults(t *testing.T) {
	deployment := &core.DeploymentConfig{
		App: &model.AppConfig{
			Name: "myapp",
		},
	}

	applyDefaults(deployment)

	assert.Equal(t, "myapp", deployment.Name)
	assert.Equal(t, "apps", deployment.Namespace)
	assert.Equal(t, "http", deployment.App.Expose.Protocol)
}

func Test_applyDefaults_AllowValues(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:      "mydeployment",
		Namespace: "not-yet-supported",
		App: &model.AppConfig{
			Name: "myapp",
			Expose: &model.AppConfigExpose{
				Protocol: "grpc",
			},
		},
	}

	applyDefaults(deployment)

	assert.Equal(t, "apps", deployment.Namespace)
	assert.Equal(t, "myapp-mydeployment", deployment.Name)
	assert.Equal(t, "grpc", deployment.App.Expose.Protocol)
}

func Test_ApplyDefaults_WhenDeploymentNameSpecified_DoesNotAddPrefixIfNamesMatch(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name: "myapp",
		App: &model.AppConfig{
			Name: "myapp",
		},
	}

	applyDefaults(deployment)

	assert.Equal(t, "myapp", deployment.Name)
}
