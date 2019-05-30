package deployment

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/pkg/state"

	"github.com/riser-platform/riser-server/api/v1/model"
)

// Uses the "golden files" or "snapshot" test pattern to capture deployments resources.
// All test data is in the /testdata folder with a subfolder for each fixture
// Pass the UPDATESNAPSHOT=true env var to "go test" to regenerate the snapshot data

func Test_update_snapshot_simple(t *testing.T) {
	newDeployment := &core.NewDeployment{
		DeploymentMeta: core.DeploymentMeta{
			Name:      "myapp",
			Namespace: "apps",
			Stage:     "dev",
			Docker: core.DeploymentDocker{
				Tag: "0.0.1",
			},
		},
		App: &model.AppConfigWithOverrides{
			AppConfig: model.AppConfig{
				Name: "myapp",
				Id:   "myid",
				Expose: &model.AppConfigExpose{
					ContainerPort: 8080,
				},
			},
		},
	}

	secretNames := []string{"mysecret"}

	dryRunComitter := state.NewDryRunComitter()
	var committer state.Committer
	snapshotDir, err := filepath.Abs("testdata/snapshots/simple")
	require.NoError(t, err)
	if shouldUpdateSnapshot() {
		fmt.Printf("Updating snapshot for %q", snapshotDir)
		err = os.RemoveAll(snapshotDir)
		require.NoError(t, err)
		committer = state.NewFileCommitter(snapshotDir)
	} else {
		committer = dryRunComitter
	}
	service := service{}
	err = service.update(newDeployment, committer, secretNames, "dev.riser.org")

	assert.NoError(t, err)
	if !shouldUpdateSnapshot() {
		require.Len(t, dryRunComitter.Commits, 1)
		assert.Equal(t, "Updating resources for \"myapp\" in stage \"dev\"", dryRunComitter.Commits[0].Message)
		AssertSnapshot(t, snapshotDir, dryRunComitter.Commits[0].Files)
	}
}

func AssertSnapshot(t *testing.T, snapshotDir string, actualFiles []core.ResourceFile) {
	actualFileMap := map[string][]byte{}
	snapshotFileMap := map[string][]byte{}

	for _, file := range actualFiles {
		actualFileMap[file.Name] = file.Contents
	}

	err := filepath.Walk(snapshotDir, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			bytes, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}
			relPath, err := filepath.Rel(snapshotDir, filePath)
			if err != nil {
				return err
			}
			snapshotFileMap[relPath] = bytes
		}
		return nil
	})

	require.NoError(t, err)

	assert.Len(t, actualFileMap, len(snapshotFileMap))

	for snapshotPath, snapshotContents := range snapshotFileMap {
		if actualFile, ok := actualFileMap[snapshotPath]; ok {
			message := fmt.Sprintf("File: %s", snapshotPath)
			assert.Equal(t, string(snapshotContents), string(actualFile), message)
		} else {
			assert.Fail(t, "Missing expected file", snapshotPath)
		}
	}
}

func shouldUpdateSnapshot() bool {
	return os.Getenv("UPDATESNAPSHOT") == "true"
}
