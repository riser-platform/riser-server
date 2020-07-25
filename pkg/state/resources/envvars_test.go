package resources

import (
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Test_k8sEnvVars(t *testing.T) {
	secrets := []core.SecretMeta{
		{Name: "secret1", Revision: 5},
		{Name: "secret2", Revision: 1},
	}
	deployment := &core.DeploymentConfig{
		Name:            "myapp-mydep",
		Namespace:       "myns",
		EnvironmentName: "myenv",
		App: &model.AppConfig{
			Name: "myapp",
			OverrideableAppConfig: model.OverrideableAppConfig{
				Environment: map[string]intstr.IntOrString{
					"env1": intstr.Parse("env1Val"),
					"env2": intstr.Parse("env2Val"),
				},
			},
		},
	}
	deploymentCtx := &core.DeploymentContext{
		DeploymentConfig: deployment,
		RiserRevision:    2,
		Secrets:          secrets}

	result := k8sEnvVars(deploymentCtx)

	assert.Len(t, result, 9)
	// User defined env
	assert.Equal(t, "ENV1", result[0].Name)
	assert.Equal(t, "env1Val", result[0].Value)
	assert.Equal(t, "ENV2", result[1].Name)
	assert.Equal(t, "env2Val", result[1].Value)
	// Platform env
	assert.Equal(t, "RISER_APP", result[2].Name)
	assert.Equal(t, "myapp", result[2].Value)
	assert.Equal(t, "RISER_DEPLOYMENT", result[3].Name)
	assert.Equal(t, "myapp-mydep", result[3].Value)
	assert.Equal(t, "RISER_DEPLOYMENT_REVISION", result[4].Name)
	assert.Equal(t, "2", result[4].Value)
	assert.Equal(t, "RISER_ENVIRONMENT", result[5].Name)
	assert.Equal(t, "myenv", result[5].Value)
	assert.Equal(t, "RISER_NAMESPACE", result[6].Name)
	assert.Equal(t, "myns", result[6].Value)

	// Secrets
	assert.Equal(t, "SECRET1", result[7].Name)
	assert.Equal(t, "myapp-secret1-5", result[7].ValueFrom.SecretKeyRef.LocalObjectReference.Name)
	assert.Equal(t, "SECRET2", result[8].Name)
	assert.Equal(t, "myapp-secret2-1", result[8].ValueFrom.SecretKeyRef.LocalObjectReference.Name)
}
