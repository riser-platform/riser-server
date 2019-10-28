package core

type StageRepository interface {
	Get(name string) (*Stage, error)
	List() ([]Stage, error)
	Save(stage *Stage) error
}

type FakeStageRepository struct {
	GetFn         func(string) (*Stage, error)
	GetCallCount  int
	ListFn        func() ([]Stage, error)
	ListCallCount int
	SaveFn        func(*Stage) error
	SaveCallCount int
}

func (fake *FakeStageRepository) List() ([]Stage, error) {
	fake.ListCallCount++
	return fake.ListFn()
}

func (fake *FakeStageRepository) Save(stage *Stage) error {
	fake.SaveCallCount++
	return fake.SaveFn(stage)
}

func (fake *FakeStageRepository) Get(name string) (*Stage, error) {
	fake.GetCallCount++
	return fake.GetFn(name)
}
