package core

type AppRepository interface {
	Get(name string) (*App, error)
	Create(app *App) error
	ListApps() ([]App, error)
}

type FakeAppRepository struct {
	GetFn        func(string) (*App, error)
	GetCallCount int
	CreateFn     func(app *App) error
	ListAppsFn   func() ([]App, error)
}

func (fake *FakeAppRepository) Get(name string) (*App, error) {
	fake.GetCallCount++
	return fake.GetFn(name)
}

func (fake *FakeAppRepository) Create(app *App) error {
	return fake.CreateFn(app)
}

func (fake *FakeAppRepository) ListApps() ([]App, error) {
	return fake.ListAppsFn()
}
