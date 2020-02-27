package app

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/pkg/core"
)

func Test_GetByIdOrName_UUID(t *testing.T) {
	appId := uuid.New()
	appRepository := &core.FakeAppRepository{
		GetFn: func(idArg uuid.UUID) (*core.App, error) {
			return &core.App{Id: idArg}, nil
		},
	}

	appService := service{appRepository}

	result, err := appService.GetByIdOrName(appId.String())

	assert.NoError(t, err)
	assert.Equal(t, appId, result.Id)
}

func Test_GetByIdOrName_UUID_WhenError(t *testing.T) {
	appId := uuid.New()
	appRepository := &core.FakeAppRepository{
		GetFn: func(idArg uuid.UUID) (*core.App, error) {
			return nil, errors.New("test")
		},
	}

	appService := service{appRepository}

	result, err := appService.GetByIdOrName(appId.String())

	assert.Equal(t, "test", err.Error())
	assert.Nil(t, result)
}

func Test_GetByIdOrName_Name(t *testing.T) {
	appName := "myapp"
	appRepository := &core.FakeAppRepository{
		GetByNameFn: func(nameArg string) (*core.App, error) {
			return &core.App{Name: nameArg}, nil
		},
	}

	appService := service{appRepository}

	result, err := appService.GetByIdOrName(appName)

	assert.NoError(t, err)
	assert.Equal(t, appName, result.Name)
}

func Test_CreateApp(t *testing.T) {
	var appName string
	var newApp *core.App
	appRepository := &core.FakeAppRepository{
		GetByNameFn: func(nameArg string) (*core.App, error) {
			appName = nameArg
			return nil, core.ErrNotFound
		},
		CreateFn: func(newAppArg *core.App) error {
			newApp = newAppArg
			return nil
		},
	}

	appService := service{appRepository}

	result, err := appService.CreateApp("foo")

	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.Id)
	assert.Equal(t, "foo", appName)
	assert.Equal(t, "foo", newApp.Name)
	assert.Equal(t, newApp.Id, result.Id)
}

func Test_CreateApp_WhenAppExists_ReturnsErr(t *testing.T) {
	appRepository := &core.FakeAppRepository{
		GetByNameFn: func(nameArg string) (*core.App, error) {
			return &core.App{}, nil
		},
	}

	appService := service{
		apps: appRepository,
	}

	result, err := appService.CreateApp("foo")

	assert.Nil(t, result)
	assert.Equal(t, err, ErrAlreadyExists)
}

func Test_CreateApp_WhenErrorCheckingApp_ReturnsErr(t *testing.T) {
	expectedErr := errors.New("error")
	appRepository := &core.FakeAppRepository{
		GetByNameFn: func(nameArg string) (*core.App, error) {
			return &core.App{}, expectedErr
		},
	}

	appService := service{
		apps: appRepository,
	}

	result, err := appService.CreateApp("foo")

	assert.Nil(t, result)
	assert.Equal(t, err.Error(), "unable to validate app: error")
}

func Test_CreateApp_WhenCreateFails_ReturnsErr(t *testing.T) {
	expectedErr := errors.New("error")
	appRepository := &core.FakeAppRepository{
		GetByNameFn: func(nameArg string) (*core.App, error) {
			return nil, core.ErrNotFound
		},
		CreateFn: func(*core.App) error {
			return expectedErr
		},
	}

	appService := service{
		apps: appRepository,
	}

	result, err := appService.CreateApp("foo")

	assert.Equal(t, err, expectedErr)
	require.Nil(t, result)
}

func Test_CheckAppName(t *testing.T) {
	appId := uuid.New()
	var receivedId uuid.UUID
	appRepository := &core.FakeAppRepository{
		GetFn: func(id uuid.UUID) (*core.App, error) {
			receivedId = id
			return &core.App{Id: appId, Name: "myapp"}, nil
		},
	}

	appService := service{
		apps: appRepository,
	}

	err := appService.CheckAppName(appId, "myapp")

	assert.NoError(t, err)
	assert.Equal(t, appId, receivedId)
}

func Test_CheckAppName_WhenAppHasDifferentName_ReturnsErr(t *testing.T) {
	appRepository := &core.FakeAppRepository{
		GetFn: func(id uuid.UUID) (*core.App, error) {
			return &core.App{Id: uuid.New(), Name: "another-name"}, nil
		},
	}

	appService := service{
		apps: appRepository,
	}

	err := appService.CheckAppName(uuid.New(), "myapp")

	assert.Equal(t, ErrInvalidAppName, err)
}

func Test_CheckAppName_WhenAppDoesNotExist_ReturnsErr(t *testing.T) {
	appRepository := &core.FakeAppRepository{
		GetFn: func(uuid.UUID) (*core.App, error) {
			return nil, core.ErrNotFound
		},
	}

	appService := service{
		apps: appRepository,
	}

	err := appService.CheckAppName(uuid.New(), "myapp")

	assert.Equal(t, ErrAppNotFound, err)
}

func Test_CheckAppName_WhenRepositoryError_ReturnsErr(t *testing.T) {
	appRepository := &core.FakeAppRepository{
		GetFn: func(uuid.UUID) (*core.App, error) {
			return nil, errors.New("error")
		},
	}

	appService := service{
		apps: appRepository,
	}

	err := appService.CheckAppName(uuid.New(), "myapp")

	assert.Equal(t, "Error getting app: error", err.Error())
}
