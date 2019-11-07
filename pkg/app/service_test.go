package app

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/pkg/core"
)

func Test_CreateApp(t *testing.T) {
	var appName string
	var newApp *core.App
	appRepository := &core.FakeAppRepository{
		GetFn: func(nameArg string) (*core.App, error) {
			appName = nameArg
			return nil, core.ErrNotFound
		},
		CreateFn: func(newAppArg *core.App) error {
			newApp = newAppArg
			return nil
		},
	}

	appService := service{
		apps: appRepository,
	}

	result, err := appService.CreateApp("foo")

	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Regexp(t, "[a-f0-9]", result.Hashid)
	assert.Equal(t, "foo", appName)
	assert.Equal(t, "foo", newApp.Name)
	assert.Equal(t, newApp.Hashid, result.Hashid)
}

func Test_CreateApp_WhenAppExists_ReturnsErr(t *testing.T) {
	appRepository := &core.FakeAppRepository{
		GetFn: func(nameArg string) (*core.App, error) {
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
		GetFn: func(nameArg string) (*core.App, error) {
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
		GetFn: func(nameArg string) (*core.App, error) {
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

func Test_CheckAppId(t *testing.T) {
	appId := createAppId()
	var receivedName string
	appRepository := &core.FakeAppRepository{
		GetFn: func(name string) (*core.App, error) {
			receivedName = name
			return &core.App{Hashid: appId}, nil
		},
	}

	appService := service{
		apps: appRepository,
	}

	err := appService.CheckAppId("myapp", appId)

	assert.NoError(t, err)
	assert.Equal(t, "myapp", receivedName)
}

func Test_CheckAppId_WhenInvalidAppId_ReturnsErr(t *testing.T) {
	appRepository := &core.FakeAppRepository{
		GetFn: func(name string) (*core.App, error) {
			return &core.App{Hashid: createAppId()}, nil
		},
	}

	appService := service{
		apps: appRepository,
	}

	err := appService.CheckAppId("myapp", createAppId())

	assert.Equal(t, ErrInvalidAppId, err)
}

func Test_CheckAppId_WhenAppDoesNotExist_ReturnsErr(t *testing.T) {
	appRepository := &core.FakeAppRepository{
		GetFn: func(name string) (*core.App, error) {
			return nil, core.ErrNotFound
		},
	}

	appService := service{
		apps: appRepository,
	}

	err := appService.CheckAppId("myapp", createAppId())

	assert.Equal(t, ErrAppNotFound, err)
}

func Test_CheckAppId_WhenRepositoryError_ReturnsErr(t *testing.T) {
	expectedErr := errors.New("error")
	appRepository := &core.FakeAppRepository{
		GetFn: func(name string) (*core.App, error) {
			return nil, expectedErr
		},
	}

	appService := service{
		apps: appRepository,
	}

	err := appService.CheckAppId("myapp", createAppId())

	assert.Equal(t, "Error getting app: error", err.Error())
}

func Test_createAppId(t *testing.T) {
	result1 := createAppId()
	result2 := createAppId()

	assert.Equal(t, appIdSizeInBytes, len(result1))
	assert.Regexp(t, appIdSizeInBytes*2, len(result1.String()))
	assert.Regexp(t, "[a-f0-9]", result1.String())
	assert.NotEqual(t, result1, result2, "the appId should be unique every time")
}
