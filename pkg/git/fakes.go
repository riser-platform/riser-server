package git

import (
	"github.com/riser-platform/riser-server/pkg/core"
)

type FakeGitProvider struct {
	CommitFn                 func(message string, files []core.ResourceFile) (string, error)
	CommitCallCount          int
	PushFn                   func() error
	PushCallCount            int
	ResetHardRemoteFn        func() error
	ResetHardRemoteCallCount int
}

func (fake *FakeGitProvider) Commit(message string, files []core.ResourceFile) (string, error) {
	fake.CommitCallCount++
	return fake.CommitFn(message, files)
}

func (fake *FakeGitProvider) Push() error {
	fake.PushCallCount++
	return fake.PushFn()
}

func (fake *FakeGitProvider) ResetHardRemote() error {
	fake.ResetHardRemoteCallCount++
	return fake.ResetHardRemoteFn()
}
