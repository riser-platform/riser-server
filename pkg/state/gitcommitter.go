package state

import (
	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/git"

	"github.com/pkg/errors"
)

type Committer interface {
	Commit(message string, files []core.ResourceFile) error
}

type GitCommitter struct {
	git git.Repo
}

func NewGitCommitter(gitRepo git.Repo) *GitCommitter {
	return &GitCommitter{gitRepo}
}

// Commit commits state changes to the state repo. Commits are authoritative i.e. they represent the absolute desired state.
// No merging takes place for riser managed resources.
func (committer *GitCommitter) Commit(message string, files []core.ResourceFile) error {
	/*
		Commits inside of a riser server instance are atomic as we only keep one instance of the repo in /tmp

		TODO: Timeout this lock and return an error after say 10 seconds.
		The main scenario is that the repo is slow/unresponsive and we don't want a bunch of goroutines hanging here.
	*/
	committer.git.Lock()
	defer committer.git.Unlock()

	// Always reset before committing as commits are authoritative
	err := committer.git.ResetHardRemote()
	if err != nil {
		return errors.Wrap(err, "error resetting repo")
	}

	err = committer.git.Commit(message, files)
	if err != nil && err != git.ErrNoChanges {
		return errors.Wrap(err, "error committing changes")
	}

	if err == git.ErrNoChanges {
		return err
	} else {
		// TODO: There's a race condition between the reset and push if a change occurs in the repo. Need to detect this error and retry.
		err = committer.git.Push()

		if err != nil {
			return errors.Wrap(err, "error pushing changes")
		}
	}

	return nil
}
