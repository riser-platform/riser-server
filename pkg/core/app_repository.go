package core

import "github.com/google/uuid"

type AppRepository interface {
	Get(id uuid.UUID) (*App, error)
	GetByName(*NamespacedName) (*App, error)
	Create(app *App) error
	ListApps() ([]App, error)
}

type FakeAppRepository struct {
	GetFn              func(id uuid.UUID) (*App, error)
	GetCallCount       int
	GetByNameFn        func(*NamespacedName) (*App, error)
	GetByNameCallCount int
	CreateFn           func(app *App) error
	ListAppsFn         func() ([]App, error)
}

func (fake *FakeAppRepository) Get(id uuid.UUID) (*App, error) {
	fake.GetCallCount++
	return fake.GetFn(id)
}

func (fake *FakeAppRepository) GetByName(name *NamespacedName) (*App, error) {
	fake.GetByNameCallCount++
	return fake.GetByNameFn(name)
}

func (fake *FakeAppRepository) Create(app *App) error {
	return fake.CreateFn(app)
}

func (fake *FakeAppRepository) ListApps() ([]App, error) {
	return fake.ListAppsFn()
}
