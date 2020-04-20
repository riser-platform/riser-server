package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertSnapshot asserts that a snapshot (aka "gold files") are equal to the generated resource files
func AssertSnapshot(t *testing.T, snapshotDir string, actualFiles []core.ResourceFile) {
	actualFileMap := map[string][]byte{}
	snapshotFileMap := map[string][]byte{}

	for _, file := range actualFiles {
		actualFileMap[file.Name] = file.Contents
	}

	err := filepath.Walk(snapshotDir, func(filePath string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
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

// ShouldUpdateSnapshot checks for the UPDATESNAPSHOT env var to determine if snapshot content should be updated
func ShouldUpdateSnapshot() bool {
	return os.Getenv("UPDATESNAPSHOT") == "true"
}
