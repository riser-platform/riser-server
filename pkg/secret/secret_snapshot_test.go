package secret

import (
	"encoding/base64"
	"path/filepath"
	"testing"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/snapshot"
	"github.com/riser-platform/riser-server/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Uses the "golden files" or "snapshot" test pattern to capture deployments resources.
// All test data is in the /testdata folder with a subfolder for each fixture
// Pass the UPDATESNAPSHOT=true env var to "go test" to regenerate the snapshot data

func Test_update_snapshot_sealedsecret(t *testing.T) {
	secretMeta := &core.SecretMeta{
		Name:            "mysecret",
		App:             core.NewNamespacedName("myapp", "apps"),
		Revision:        2,
		EnvironmentName: "dev",
	}

	testCertBytes, _ := base64.StdEncoding.DecodeString(testCert)
	environmentRepository := &core.FakeEnvironmentRepository{
		GetFn: func(envName string) (*core.Environment, error) {
			return &core.Environment{
				Name: "myenv",
				Doc: core.EnvironmentDoc{
					Config: core.EnvironmentConfig{
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

	snapshotPath, err := filepath.Abs("testdata/snapshots/sealedsecret")
	require.NoError(t, err)

	committer, err := snapshot.CreateCommitter(snapshotPath)
	require.NoError(t, err)

	secretService := service{secretMetaRepository, environmentRepository, staticReader{}}

	err = secretService.SealAndSave("mysecretval", secretMeta, committer)

	assert.NoError(t, err)
	if !snapshot.ShouldUpdate() {
		dryRunCommitter := committer.(*state.DryRunCommitter)
		snapshot.AssertCommitter(t, snapshotPath, dryRunCommitter)
		assert.Equal(t, "Updating secret \"myapp-mysecret-2\" in environment \"dev\"", dryRunCommitter.Commits[0].Message)
	}
}

// io.Reader for static random bytes so that we get deterministic results
type staticReader struct {
}

func (staticReader) Read(buf []byte) (n int, err error) {
	return len(buf), nil
}
