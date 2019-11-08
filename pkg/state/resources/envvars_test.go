package resources

import (
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Test_k8sEnvVars(t *testing.T) {
	secretNames := []string{"secret2", "secret1"}
	deployment := &core.DeploymentConfig{
		App: &model.AppConfig{
			Environment: map[string]intstr.IntOrString{
				"env1": intstr.Parse("env1Val"),
				"env2": intstr.Parse("env2Val"),
			},
		},
	}

	result := k8sEnvVars(&core.DeploymentContext{Deployment: deployment, SecretNames: secretNames})

	assert.Len(t, result, 4)
	assert.Equal(t, "ENV1", result[0].Name)
	assert.Equal(t, "ENV2", result[1].Name)
	assert.Equal(t, "SECRET1", result[2].Name)
	assert.Equal(t, "SECRET2", result[3].Name)
}
