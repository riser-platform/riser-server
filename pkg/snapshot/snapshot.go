package snapshot

import (
	"fmt"
	"os"

	"github.com/riser-platform/riser-server/pkg/state"
)

// ShouldUpdate checks for the UPDATESNAPSHOT env var to determine if snapshot content should be updated
func ShouldUpdate() bool {
	return os.Getenv("UPDATESNAPSHOT") == "true"
}

func CreateCommitter(snapshotPath string) (state.Committer, error) {
	dryRunCommitter := state.NewDryRunCommitter()
	var committer state.Committer
	if ShouldUpdate() {
		fmt.Printf("Updating snapshot for %q", snapshotPath)
		err := os.RemoveAll(snapshotPath)
		if err != nil {
			return nil, err
		}
		committer = state.NewFileCommitter(snapshotPath)
	} else {
		committer = dryRunCommitter
	}
	return committer, nil
}
