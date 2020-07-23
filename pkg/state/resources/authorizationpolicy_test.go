package resources

import (
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_createHealthcheckDenyPolicy(t *testing.T) {
	ctx := &core.DeploymentContext{
		DeploymentConfig: &core.DeploymentConfig{
			Name:            "myapp-dep",
			EnvironmentName: "myenv",
			App: &model.AppConfig{
				Name: "myapp",
				HealthCheck: &model.AppConfigHealthCheck{
					Path: "/health",
				},
			},
		},
	}

	result := CreateHealthcheckDenyPolicy(ctx)

	assert.Equal(t, "myapp-dep-healthcheck-deny", result.Name)
	assert.Equal(t, deploymentLabels(ctx), result.Labels)
	assert.Equal(t, deploymentAnnotations(ctx), result.Annotations)
	assert.Equal(t, "AuthorizationPolicy", result.TypeMeta.Kind)
	assert.Equal(t, "security.istio.io/v1beta1", result.TypeMeta.APIVersion)
	assert.Equal(t, "myapp-dep", result.Spec.Selector.MatchLabels["riser.dev/deployment"])
	assert.Equal(t, "DENY", result.Spec.Action.String())
	assert.Equal(t, "/health", result.Spec.Rules[0].To[0].Operation.Paths[0])
}

func Test_createHealthcheckDenyPolicy_NoHealthcheckReturnsNil(t *testing.T) {
	ctx := &core.DeploymentContext{
		DeploymentConfig: &core.DeploymentConfig{
			Name:            "myapp-dep",
			EnvironmentName: "myenv",
			App: &model.AppConfig{
				Name: "myapp",
			},
		},
	}

	result := CreateHealthcheckDenyPolicy(ctx)

	assert.Nil(t, result)
}
