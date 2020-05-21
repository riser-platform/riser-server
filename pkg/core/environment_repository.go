package core

type EnvironmentRepository interface {
	Get(name string) (*Environment, error)
	List() ([]Environment, error)
	Save(environment *Environment) error
}

type FakeEnvironmentRepository struct {
	GetFn         func(string) (*Environment, error)
	GetCallCount  int
	ListFn        func() ([]Environment, error)
	ListCallCount int
	SaveFn        func(*Environment) error
	SaveCallCount int
}

func (fake *FakeEnvironmentRepository) List() ([]Environment, error) {
	fake.ListCallCount++
	return fake.ListFn()
}

func (fake *FakeEnvironmentRepository) Save(environment *Environment) error {
	fake.SaveCallCount++
	return fake.SaveFn(environment)
}

func (fake *FakeEnvironmentRepository) Get(name string) (*Environment, error) {
	fake.GetCallCount++
	return fake.GetFn(name)
}
