package resources

import (
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/util"
	"github.com/stretchr/testify/assert"
)

func Test_createRevisionMeta(t *testing.T) {
	ctx := &core.DeploymentContext{
		Deployment: &core.DeploymentConfig{
			Name:  "myapp",
			Stage: "mystage",
			App: &model.AppConfig{
				Name: "myapp",
			},
		},
		RiserRevision: 1,
	}

	result := createRevisionMeta(ctx)

	assert.Equal(t, "myapp-1", result.Name)
	assert.Equal(t, deploymentLabels(ctx), result.Labels)
	assert.Equal(t, deploymentAnnotations(ctx), result.Annotations)
}

func Test_createRevisionMeta_Autoscale(t *testing.T) {
	ctx := &core.DeploymentContext{
		Deployment: &core.DeploymentConfig{
			Name:  "myapp",
			Stage: "mystage",
			App: &model.AppConfig{
				Name: "myapp",
				Autoscale: &model.AppConfigAutoscale{
					Min: util.PtrInt(1),
					Max: util.PtrInt(2),
				},
			},
		},
		RiserRevision: 1,
	}

	result := createRevisionMeta(ctx)

	assert.Len(t, result.Annotations, 4)
	assert.Equal(t, "1", result.Annotations["autoscaling.knative.dev/minScale"])
	assert.Equal(t, "2", result.Annotations["autoscaling.knative.dev/maxScale"])
	assert.Equal(t, "1", result.Annotations["riser.dev/revision"])
	assert.Equal(t, util.VersionString, result.Annotations["riser.dev/server-version"])
}
