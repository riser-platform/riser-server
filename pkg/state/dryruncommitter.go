package state

import "github.com/riser-platform/riser-server/pkg/core"

type DryRunCommit struct {
	Message string
	Files   []core.ResourceFile
}

type DryRunComitter struct {
	Commits []DryRunCommit
}

func NewDryRunComitter() *DryRunComitter {
	return &DryRunComitter{
		Commits: []DryRunCommit{},
	}
}

func (committer *DryRunComitter) Commit(message string, files []core.ResourceFile) error {
	committer.Commits = append(committer.Commits, DryRunCommit{Message: message, Files: files})
	return nil
}
