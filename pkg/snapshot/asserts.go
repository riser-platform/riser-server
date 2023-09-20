package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertEqual asserts that a snapshot (aka "gold files") are equal to the generated resource files
func AssertEqual(t *testing.T, snapshotDir string, actualFiles []core.ResourceFile) {
	actualFileMap := map[string][]byte{}
	snapshotFileMap := map[string][]byte{}

	for _, file := range actualFiles {
		actualFileMap[file.Name] = file.Contents
	}

	err := filepath.Walk(snapshotDir, func(filePath string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			bytes, err := os.ReadFile(filePath)
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

// AssertCommitter asserts that the snapshotPath matches the state in a DryRunCommitter.
// This assert always passes if the env var UPDATESNAPSHOT == true
func AssertCommitter(t *testing.T, snapshotPath string, committer *state.DryRunCommitter) {
	if !ShouldUpdate() {
		require.IsType(t, &state.DryRunCommitter{}, committer, "The committer must be a DryRunCommitter when the env var UPDATESNAPSHOT is false")
		require.Len(t, committer.Commits, 1)

		AssertEqual(t, snapshotPath, committer.Commits[0].Files)
	}
}
