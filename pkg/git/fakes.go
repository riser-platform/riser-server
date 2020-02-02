package git

import (
	"github.com/riser-platform/riser-server/pkg/core"
)

type FakeRepo struct {
	CommitFn                 func(message string, files []core.ResourceFile) error
	CommitCallCount          int
	PushFn                   func() error
	PushCallCount            int
	ResetHardRemoteFn        func() error
	ResetHardRemoteCallCount int
}

func (fake *FakeRepo) Commit(message string, files []core.ResourceFile) error {
	fake.CommitCallCount++
	return fake.CommitFn(message, files)
}

func (fake *FakeRepo) Push() error {
	fake.PushCallCount++
	return fake.PushFn()
}

func (fake *FakeRepo) ResetHardRemote() error {
	fake.ResetHardRemoteCallCount++
	return fake.ResetHardRemoteFn()
}
