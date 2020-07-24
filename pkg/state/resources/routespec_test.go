package resources

import (
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_CreateKNativeRoute_ExposeCluster(t *testing.T) {
	visiblityLabel := "serving.knative.dev/visibility"

	tests := []struct {
		scope          string
		isClusterLocal bool
	}{
		{model.AppExposeScope_Cluster, true},
		{model.AppExposeScope_External, false},
		// Any future app expose options should not expose the app by default
		{"anything-else", true},
	}

	for _, tt := range tests {
		ctx := &core.DeploymentContext{
			DeploymentConfig: &core.DeploymentConfig{
				App: &model.AppConfig{
					Expose: &model.AppConfigExpose{
						Scope: tt.scope,
					},
				},
			},
		}

		result := CreateKNativeRoute(ctx)

		if tt.isClusterLocal {
			assert.Equal(t, result.Labels[visiblityLabel], "cluster-local")
		} else {
			assert.NotContains(t, visiblityLabel, result.Labels)
		}

	}
}
