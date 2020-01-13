package state

import (
	"sync"
	"testing"
	"time"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/git"
	"github.com/stretchr/testify/assert"
)

func Test_Commit(t *testing.T) {
	gitProvider := &git.FakeGitProvider{
		ResetHardRemoteFn: func() error {
			return nil
		},
		CommitFn: func(message string, resources []core.ResourceFile) (string, error) {
			assert.Equal(t, "test message", message)
			assert.Len(t, resources, 1)
			assert.Equal(t, "test.yaml", resources[0].Name)
			return "", nil
		},
		PushFn: func() error {
			return nil
		},
	}
	committer := NewGitCommitter(gitProvider)

	resources := []core.ResourceFile{
		core.ResourceFile{
			Name: "test.yaml",
		},
	}

	result := committer.Commit("test message", resources)

	assert.NoError(t, result)
	assert.Equal(t, 1, gitProvider.ResetHardRemoteCallCount)
	assert.Equal(t, 1, gitProvider.CommitCallCount)
	assert.Equal(t, 1, gitProvider.PushCallCount)
}

func Test_Commit_NoChanges_DoesNotPush(t *testing.T) {
	gitProvider := &git.FakeGitProvider{
		ResetHardRemoteFn: func() error {
			return nil
		},
		CommitFn: func(message string, resources []core.ResourceFile) (string, error) {
			return "", git.ErrNoChanges
		},
		PushFn: func() error {
			return nil
		},
	}
	committer := NewGitCommitter(gitProvider)

	resources := []core.ResourceFile{
		core.ResourceFile{
			Name: "test.yaml",
		},
	}

	result := committer.Commit("test message", resources)

	assert.Equal(t, git.ErrNoChanges, result)
	assert.Equal(t, 0, gitProvider.PushCallCount)
}

func Test_Commit_Serialized(t *testing.T) {
	inTransaction := false
	gitProvider := &git.FakeGitProvider{
		ResetHardRemoteFn: func() error {
			assert.False(t, inTransaction, "Must not reset while a transaction is pending")
			inTransaction = true
			time.Sleep(10 * time.Millisecond)
			return nil
		},
		CommitFn: func(message string, resources []core.ResourceFile) (string, error) {
			assert.True(t, inTransaction, "Must not commit while not inside a transaction")
			return "", nil
		},
		PushFn: func() error {
			assert.True(t, inTransaction, "Must not push while not inside a transaction")
			time.Sleep(10 * time.Millisecond)
			inTransaction = false
			return nil
		},
	}

	committer := NewGitCommitter(gitProvider)

	wg := sync.WaitGroup{}

	doCommit := func(committer *GitCommitter) {
		wg.Add(1)
		err := committer.Commit("", []core.ResourceFile{core.ResourceFile{}})
		wg.Done()
		assert.NoError(t, err)
	}

	for i := 0; i < 3; i++ {
		go doCommit(committer)
	}

	wg.Wait()
}
