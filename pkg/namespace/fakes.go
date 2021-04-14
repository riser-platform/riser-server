package namespace

type FakeService struct {
	ValidateDeployableFn func(string) error
}

func (fake *FakeService) ValidateDeployable(namespaceName string) error {
	return fake.ValidateDeployableFn(namespaceName)
}
func (fake *FakeService) EnsureDefaultNamespace() error {
	panic("NI")
}
func (fake *FakeService) Create(namespaceName string) error {
	panic("NI")
}
