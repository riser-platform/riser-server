package state

import "github.com/riser-platform/riser-server/pkg/core"

type DryRunCommit struct {
	Message string
	Files   []core.ResourceFile
}

type DryRunCommitter struct {
	Commits []DryRunCommit
}

func NewDryRunCommitter() *DryRunCommitter {
	return &DryRunCommitter{
		Commits: []DryRunCommit{},
	}
}

func (committer *DryRunCommitter) Commit(message string, files []core.ResourceFile) error {
	committer.Commits = append(committer.Commits, DryRunCommit{Message: message, Files: files})
	return nil
}
