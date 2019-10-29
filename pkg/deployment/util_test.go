package deployment

import (
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_ApplyDefaults_Defaults(t *testing.T) {
	deployment := &core.Deployment{
		App: &model.AppConfig{
			Name: "myapp",
		},
	}

	result := ApplyDefaults(deployment)

	assert.Equal(t, "myapp", result.Name)
	assert.Equal(t, "apps", result.Namespace)
	assert.Equal(t, "http", result.App.Expose.Protocol)
}

func Test_ApplyDefaults_AllowValues(t *testing.T) {
	deployment := &core.Deployment{
		DeploymentMeta: core.DeploymentMeta{
			Name:      "mydeployment",
			Namespace: "not-yet-supported",
		},
		App: &model.AppConfig{
			Name: "myapp",
			Expose: &model.AppConfigExpose{
				Protocol: "grpc",
			},
		},
	}

	result := ApplyDefaults(deployment)

	assert.Equal(t, "apps", result.Namespace)
	assert.Equal(t, "myapp-mydeployment", result.Name)
	assert.Equal(t, "grpc", result.App.Expose.Protocol)
}

func Test_ApplyDefaults_WhenDeploymentNameSpecified_DoesNotAddPrefixIfNamesMatch(t *testing.T) {
	deployment := &core.Deployment{
		DeploymentMeta: core.DeploymentMeta{
			Name: "myapp",
		},
		App: &model.AppConfig{
			Name: "myapp",
		},
	}

	result := ApplyDefaults(deployment)

	assert.Equal(t, "myapp", result.Name)
}
