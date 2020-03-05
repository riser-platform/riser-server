package core

type NamespaceRepository interface {
	Create(namespace *Namespace) error
	Get(namespaceName string) (*Namespace, error)
	List() ([]Namespace, error)
}

type FakeNamespaceRepository struct {
	CreateFn     func(namespace *Namespace) error
	GetFn        func(namespaceName string) (*Namespace, error)
	GetCallCount int
	ListFn       func() ([]Namespace, error)
}

func (fake *FakeNamespaceRepository) Create(namespace *Namespace) error {
	return fake.CreateFn(namespace)
}

func (fake *FakeNamespaceRepository) Get(namespaceName string) (*Namespace, error) {
	fake.GetCallCount++
	return fake.GetFn(namespaceName)
}

func (fake *FakeNamespaceRepository) List() ([]Namespace, error) {
	return fake.ListFn()
}
