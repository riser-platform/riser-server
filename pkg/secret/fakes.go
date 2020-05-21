package secret

import (
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/state"
)

type FakeService struct {
	SealAndSaveFn        func(plaintextSecret string, secretMeta *core.SecretMeta, committer state.Committer) error
	SealAndSaveCallCount int
}

func (f *FakeService) SealAndSave(plaintextSecret string, secretMeta *core.SecretMeta, committer state.Committer) error {
	f.SealAndSaveCallCount++
	return f.SealAndSaveFn(plaintextSecret, secretMeta, committer)
}
