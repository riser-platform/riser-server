package rollout

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/state"
	"github.com/riser-platform/riser-server/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Uses the "golden files" or "snapshot" test pattern to capture deployments resources.
// All test data is in the /testdata folder with a subfolder for each fixture
// Pass the UPDATESNAPSHOT=true env var to "go test" to regenerate the snapshot data

func Test_update_snapshot_rollout(t *testing.T) {
	appId := uuid.New()
	name := core.NewNamespacedName("myapp", "myns")

	deployments := &core.FakeDeploymentRepository{
		GetByNameFn: func(nameArg *core.NamespacedName, stageName string) (*core.Deployment, error) {
			assert.Equal(t, name, nameArg)
			return &core.Deployment{
				DeploymentReservation: core.DeploymentReservation{
					Name:  "myapp",
					AppId: appId,
				},
				DeploymentRecord: core.DeploymentRecord{
					Doc: core.DeploymentDoc{
						Status: &core.DeploymentStatus{
							Revisions: []core.DeploymentRevisionStatus{
								{
									RiserRevision: 1,
								},
							},
						},
					},
				},
			}, nil
		},
	}

	apps := &core.FakeAppRepository{
		GetFn: func(id uuid.UUID) (*core.App, error) {
			assert.Equal(t, appId, id)
			return &core.App{
				Id:   id,
				Name: "myapp",
			}, nil
		},
	}

	traffic := core.TrafficConfig{
		core.TrafficConfigRule{
			RiserRevision: 1,
			Percent:       100,
		},
	}

	svc := service{apps, deployments}

	dryRunCommitter := state.NewDryRunCommitter()
	var committer state.Committer
	snapshotDir, err := filepath.Abs("testdata/snapshots/rollout")
	require.NoError(t, err)

	if util.ShouldUpdateSnapshot() {
		fmt.Printf("Updating snapshot for %q", snapshotDir)
		err = os.RemoveAll(snapshotDir)
		require.NoError(t, err)
		committer = state.NewFileCommitter(snapshotDir)
	} else {
		committer = dryRunCommitter
	}

	err = svc.UpdateTraffic(name, "dev", traffic, committer)

	assert.NoError(t, err)
	if !util.ShouldUpdateSnapshot() {
		require.Len(t, dryRunCommitter.Commits, 1)
		assert.Equal(t, `Updating resources for "myapp.myns" in stage "dev"`, dryRunCommitter.Commits[0].Message)
		util.AssertSnapshot(t, snapshotDir, dryRunCommitter.Commits[0].Files)
	}
}
