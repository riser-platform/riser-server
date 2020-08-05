package namespace

import "github.com/riser-platform/riser-server/pkg/state"

type FakeService struct {
	ValidateDeployableFn func(string) error
}

func (fake *FakeService) ValidateDeployable(namespaceName string) error {
	return fake.ValidateDeployableFn(namespaceName)
}
func (fake *FakeService) EnsureDefaultNamespace(committer state.Committer) error {
	panic("NI")
}
func (fake *FakeService) Create(namespaceName string, committer state.Committer) error {
	panic("NI")
}
