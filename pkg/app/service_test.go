package app

import (
	"testing"

	"github.com/riser-platform/riser-server/pkg/namespace"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/pkg/core"
)

func Test_Create(t *testing.T) {
	var newApp *core.App
	appRepository := &core.FakeAppRepository{
		GetByNameFn: func(nameArg *core.NamespacedName) (*core.App, error) {
			assert.Equal(t, "foo", nameArg.Name)
			assert.Equal(t, "myns", nameArg.Namespace)
			return nil, core.ErrNotFound
		},
		CreateFn: func(newAppArg *core.App) error {
			newApp = newAppArg
			return nil
		},
	}

	namespaceService := &namespace.FakeService{
		ValidateDeployableFn: func(nsArg string) error {
			assert.Equal(t, "myns", nsArg)
			return nil
		},
	}

	appService := service{appRepository, namespaceService}

	result, err := appService.Create(core.NewNamespacedName("foo", "myns"))

	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Id)
	assert.Equal(t, "foo", newApp.Name)
	assert.Equal(t, "myns", newApp.Namespace)
	assert.Equal(t, newApp.Id, result.Id)
}

func Test_Create_InvalidNamespace(t *testing.T) {
	appRepository := &core.FakeAppRepository{
		GetByNameFn: func(nameArg *core.NamespacedName) (*core.App, error) {
			return nil, core.ErrNotFound
		},
	}

	namespaceService := &namespace.FakeService{
		ValidateDeployableFn: func(nsArg string) error {
			assert.Equal(t, "myns", nsArg)
			return errors.New("test")
		},
	}

	appService := service{appRepository, namespaceService}

	result, err := appService.Create(core.NewNamespacedName("foo", "myns"))

	assert.Nil(t, result)
	assert.Equal(t, "test", err.Error())
}

func Test_Create_WhenAppExists_ReturnsErr(t *testing.T) {
	appRepository := &core.FakeAppRepository{
		GetByNameFn: func(*core.NamespacedName) (*core.App, error) {
			return &core.App{}, nil
		},
	}

	appService := service{
		apps: appRepository,
	}

	result, err := appService.Create(core.NewNamespacedName("foo", "myns"))

	assert.Nil(t, result)
	assert.Equal(t, err, ErrAlreadyExists)
}

func Test_Create_WhenErrorCheckingApp_ReturnsErr(t *testing.T) {
	expectedErr := errors.New("error")
	appRepository := &core.FakeAppRepository{
		GetByNameFn: func(*core.NamespacedName) (*core.App, error) {
			return &core.App{}, expectedErr
		},
	}

	appService := service{
		apps: appRepository,
	}

	result, err := appService.Create(core.NewNamespacedName("foo", "myns"))

	assert.Nil(t, result)
	assert.Equal(t, err.Error(), "unable to validate app: error")
}

func Test_Create_WhenCreateFails_ReturnsErr(t *testing.T) {
	expectedErr := errors.New("error")
	appRepository := &core.FakeAppRepository{
		GetByNameFn: func(*core.NamespacedName) (*core.App, error) {
			return nil, core.ErrNotFound
		},
		CreateFn: func(*core.App) error {
			return expectedErr
		},
	}

	namespaceService := &namespace.FakeService{
		ValidateDeployableFn: func(nsArg string) error {
			return nil
		},
	}

	appService := service{appRepository, namespaceService}

	result, err := appService.Create(core.NewNamespacedName("foo", "myns"))

	assert.Equal(t, err, expectedErr)
	require.Nil(t, result)
}

func Test_CheckID(t *testing.T) {
	appId := uuid.New()
	var receivedId uuid.UUID
	appRepository := &core.FakeAppRepository{
		GetFn: func(id uuid.UUID) (*core.App, error) {
			receivedId = id
			return &core.App{Id: appId, Name: "myapp", Namespace: "myns"}, nil
		},
	}

	appService := service{
		apps: appRepository,
	}

	err := appService.CheckID(appId, core.NewNamespacedName("myapp", "myns"))

	assert.NoError(t, err)
	assert.Equal(t, appId, receivedId)
}

func Test_CheckID_WhenAppHasDifferentName_ReturnsErr(t *testing.T) {
	appRepository := &core.FakeAppRepository{
		GetFn: func(id uuid.UUID) (*core.App, error) {
			return &core.App{Id: uuid.New(), Name: "another-name", Namespace: "myns"}, nil
		},
	}

	appService := service{
		apps: appRepository,
	}

	err := appService.CheckID(uuid.New(), core.NewNamespacedName("myapp", "myns"))

	assert.Equal(t, ErrInvalidAppName, err)
}

func Test_CheckID_WhenAppHasDifferentNamespace_ReturnsErr(t *testing.T) {
	appRepository := &core.FakeAppRepository{
		GetFn: func(id uuid.UUID) (*core.App, error) {
			return &core.App{Id: uuid.New(), Name: "myapp", Namespace: "another-ns"}, nil
		},
	}

	appService := service{
		apps: appRepository,
	}

	err := appService.CheckID(uuid.New(), core.NewNamespacedName("myapp", "myns"))

	assert.Equal(t, ErrInvalidAppNamespace, err)
}

func Test_CheckID_WhenAppDoesNotExist_ReturnsErr(t *testing.T) {
	appRepository := &core.FakeAppRepository{
		GetFn: func(uuid.UUID) (*core.App, error) {
			return nil, core.ErrNotFound
		},
	}

	appService := service{
		apps: appRepository,
	}

	err := appService.CheckID(uuid.New(), core.NewNamespacedName("myapp", "myns"))

	assert.Equal(t, ErrAppNotFound, err)
}

func Test_CheckID_WhenRepositoryError_ReturnsErr(t *testing.T) {
	appRepository := &core.FakeAppRepository{
		GetFn: func(uuid.UUID) (*core.App, error) {
			return nil, errors.New("error")
		},
	}

	appService := service{
		apps: appRepository,
	}

	err := appService.CheckID(uuid.New(), core.NewNamespacedName("myapp", "myns"))

	assert.Equal(t, "Error getting app: error", err.Error())
}

func Test_GetAppByName(t *testing.T) {
	app := &core.App{}
	appRepository := &core.FakeAppRepository{
		GetByNameFn: func(name *core.NamespacedName) (*core.App, error) {
			assert.Equal(t, "myapp", name.Name)
			assert.Equal(t, "myns", name.Namespace)
			return app, nil
		},
	}

	appService := service{
		apps: appRepository,
	}

	result, err := appService.GetByName(core.NewNamespacedName("myapp", "myns"))

	assert.NoError(t, err)
	assert.Equal(t, app, result)
}

func Test_GetAppByName_NotFound(t *testing.T) {
	appRepository := &core.FakeAppRepository{
		GetByNameFn: func(name *core.NamespacedName) (*core.App, error) {
			return nil, core.ErrNotFound
		},
	}

	appService := service{
		apps: appRepository,
	}

	result, err := appService.GetByName(core.NewNamespacedName("myapp", "myns"))

	assert.Equal(t, ErrAppNotFound, err)
	assert.Nil(t, result)
}
