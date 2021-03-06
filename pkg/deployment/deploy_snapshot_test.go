package deployment

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/pkg/snapshot"
	"github.com/riser-platform/riser-server/pkg/state"
	"github.com/riser-platform/riser-server/pkg/util"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/api/v1/model"
)

// Uses the "golden files" or "snapshot" test pattern to capture deployments resources.
// All test data is in the /testdata folder with a subfolder for each fixture
// Pass the UPDATESNAPSHOT=true env var to "go test" to regenerate the snapshot data

func Test_update_snapshot_simple(t *testing.T) {
	newDeployment := &core.DeploymentConfig{
		Name:            "myapp",
		Namespace:       "apps",
		EnvironmentName: "dev",
		Docker: core.DeploymentDocker{
			Tag: "0.0.1",
		},
		App: &model.AppConfig{
			Name:      "myapp",
			Namespace: "apps",
			Id:        uuid.MustParse("2516D5E4-1EC3-46B8-B3CD-C3D72AE38DC0"),
			Expose: &model.AppConfigExpose{
				ContainerPort: 8080,
				Protocol:      "http",
				Scope:         model.AppExposeScope_External,
			},
			HealthCheck: &model.AppConfigHealthCheck{
				Path: "/health",
			},
			OverrideableAppConfig: model.OverrideableAppConfig{
				Autoscale: &model.AppConfigAutoscale{
					Min: util.PtrInt(0),
					Max: util.PtrInt(1),
				},
				Environment: map[string]intstr.IntOrString{
					"myenv": intstr.FromString("myval"),
				},
			},
		},
		Traffic: core.TrafficConfig{
			core.TrafficConfigRule{
				RiserRevision: 1,
				RevisionName:  "myapp-1",
				Percent:       100,
			},
		},
	}

	secrets := []core.SecretMeta{{Name: "mysecret", Revision: 1}}

	snapshotPath, err := filepath.Abs("testdata/snapshots/simple")
	require.NoError(t, err)

	committer, err := snapshot.CreateCommitter(snapshotPath)
	require.NoError(t, err)

	ctx := &core.DeploymentContext{
		DeploymentConfig:  newDeployment,
		EnvironmentConfig: &core.EnvironmentConfig{PublicGatewayHost: "dev.riser.org"},
		RiserRevision:     3,
		Secrets:           secrets,
	}

	err = deploy(ctx, committer)
	assert.NoError(t, err)

	if !snapshot.ShouldUpdate() {
		dryRunCommitter := committer.(*state.DryRunCommitter)
		snapshot.AssertCommitter(t, snapshotPath, dryRunCommitter)
		assert.Equal(t, "Updating resources for \"myapp.apps\" in environment \"dev\"", dryRunCommitter.Commits[0].Message)
	}
}
