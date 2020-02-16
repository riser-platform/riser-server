package secret

import (
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/state"
)

type FakeService struct {
	SealAndSaveFn        func(plaintextSecret string, secretMeta *core.SecretMeta, namespace string, committer state.Committer) error
	SealAndSaveCallCount int
}

func (f *FakeService) SealAndSave(plaintextSecret string, secretMeta *core.SecretMeta, namespace string, committer state.Committer) error {
	f.SealAndSaveCallCount++
	return f.SealAndSaveFn(plaintextSecret, secretMeta, namespace, committer)
}

func (f *FakeService) FindByStage(appName, stageName string) ([]core.SecretMeta, error) {
	panic("NI")
}

func (f *FakeService) FindNamesByStage(appName, stageName string) ([]string, error) {
	panic("NI")
}
