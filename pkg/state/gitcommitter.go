package state

import (
	"sync"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/git"

	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Committer interface {
	Commit(message string, files []core.ResourceFile) error
}

type GitCommitter struct {
	git git.GitRepoProvider
	sync.Mutex
}

type KubeResource interface {
	GetName() string
	GetNamespace() string
	GetObjectKind() schema.ObjectKind
}

func NewGitComitter(gitRepo git.GitRepoProvider) *GitCommitter {
	return &GitCommitter{gitRepo, sync.Mutex{}}
}

// Commit commits state changes to the state repo. Commits are authoritative i.e. they represent the absolute desired state.
// No merging takes place for riser managed resources.
func (committer *GitCommitter) Commit(message string, files []core.ResourceFile) error {
	/*
		Commits inside of a riser server instance are atomic as we only keep one instance of the repo in /tmp

		TODO: Timeout this lock and return an error after say 10 seconds.
		The main scenario is that the repo is slow/unresponsive and we don't want a bunch of goroutines hanging here.
	*/
	committer.Lock()
	defer committer.Unlock()

	// Always reset before committing as commits are authoritative
	err := committer.git.ResetHardRemote()
	if err != nil {
		return errors.Wrap(err, "error resetting repo")
	}

	_, err = committer.git.Commit(message, files)
	if err != nil && err != git.ErrNoChanges {
		return errors.Wrap(err, "error committing changes")
	}

	if err == git.ErrNoChanges {
		return err
	} else {
		// TODO: There's a race condition between the reset and push if a push ocurrs outside of riser server. Need to detect this error and retry.
		err = committer.git.Push()

		if err != nil {
			return errors.Wrap(err, "error pushing changes")
		}
	}

	return nil
}
