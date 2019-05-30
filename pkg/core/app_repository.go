package core

type AppRepository interface {
	FindByName(string) (*App, error)
	Create(app *App) error
	ListApps() ([]App, error)
}

type FakeAppRepository struct {
	FindByNameFn func(string) (*App, error)
	CreateFn     func(app *App) error
	ListAppsFn   func() ([]App, error)
}

func (fake *FakeAppRepository) FindByName(name string) (*App, error) {
	return fake.FindByNameFn(name)
}

func (fake *FakeAppRepository) Create(app *App) error {
	return fake.CreateFn(app)
}

func (fake *FakeAppRepository) ListApps() ([]App, error) {
	return fake.ListAppsFn()
}
