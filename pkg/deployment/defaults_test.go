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

	assert.Empty(t, deployment.Name)
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
				Protocol: "http2",
			},
		},
	}

	applyDefaults(deployment)

	assert.Equal(t, "apps", deployment.Namespace)
	// Not a valid name but the defaults should not change it.
	assert.Equal(t, "mydeployment", deployment.Name)
	assert.Equal(t, "http2", deployment.App.Expose.Protocol)
}

func Test_applyDefaults_WhenDeploymentNameSpecified_DoesNotAddPrefixIfNamesMatch(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name: "myapp",
		App: &model.AppConfig{
			Name: "myapp",
		},
	}

	applyDefaults(deployment)

	assert.Equal(t, "myapp", deployment.Name)
}
