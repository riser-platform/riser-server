package secret

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/state"
	"github.com/riser-platform/riser-server/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Uses the "golden files" or "snapshot" test pattern to capture deployments resources.
// All test data is in the /testdata folder with a subfolder for each fixture
// Pass the UPDATESNAPSHOT=true env var to "go test" to regenerate the snapshot data

func Test_update_snapshot_sealedsecret(t *testing.T) {
	secretMeta := &core.SecretMeta{
		Name:      "mysecret",
		App:       core.NewNamespacedName("myapp", "apps"),
		Revision:  2,
		StageName: "dev",
	}

	testCertBytes, _ := base64.StdEncoding.DecodeString(testCert)
	stageRepository := &core.FakeStageRepository{
		GetFn: func(stageName string) (*core.Stage, error) {
			return &core.Stage{
				Name: "mystage",
				Doc: core.StageDoc{
					Config: core.StageConfig{
						SealedSecretCert: testCertBytes,
					},
				},
			}, nil
		},
	}

	secretMetaRepository := &core.FakeSecretMetaRepository{
		CommitFn: func(secretMeta *core.SecretMeta) error {
			return nil
		},
		SaveFn: func(secretMeta *core.SecretMeta) (int64, error) {
			return secretMeta.Revision, nil
		},
	}

	dryRunCommitter := state.NewDryRunCommitter()
	var committer state.Committer
	snapshotDir, err := filepath.Abs("testdata/snapshots/sealedsecret")
	require.NoError(t, err)
	if util.ShouldUpdateSnapshot() {
		fmt.Printf("Updating snapshot for %q", snapshotDir)
		err = os.RemoveAll(snapshotDir)
		require.NoError(t, err)
		committer = state.NewFileCommitter(snapshotDir)
	} else {
		committer = dryRunCommitter
	}

	secretService := service{secretMetaRepository, stageRepository, staticReader{}}

	err = secretService.SealAndSave("mysecretval", secretMeta, committer)

	assert.NoError(t, err)
	if !util.ShouldUpdateSnapshot() {
		require.Len(t, dryRunCommitter.Commits, 1)
		assert.Equal(t, "Updating secret \"myapp-mysecret-2\" in stage \"dev\"", dryRunCommitter.Commits[0].Message)
		util.AssertSnapshot(t, snapshotDir, dryRunCommitter.Commits[0].Files)
	}
}

// io.Reader for static random bytes so that we get deterministic results
type staticReader struct {
}

func (staticReader) Read(buf []byte) (n int, err error) {
	return len(buf), nil
}
