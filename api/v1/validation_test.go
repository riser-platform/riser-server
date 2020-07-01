package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/environment"

	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/app"
	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/stretchr/testify/assert"
)

var validAppConfig = &model.AppConfigWithOverrides{
	AppConfig: model.AppConfig{
		Name:      "myapp",
		Namespace: "myns",
		Id:        uuid.New(),
		Image:     "myimage",
		Expose: &model.AppConfigExpose{
			ContainerPort: 80,
		},
	},
}

func Test_PostValidateAppConfig(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/validate/appconfig", safeMarshal(validAppConfig))
	req.Header.Add("CONTENT-TYPE", "application/json")
	ctx, _ := newContextWithRecorder(req)

	appService := &app.FakeService{
		CheckAppNameFn: func(id uuid.UUID, name *core.NamespacedName) error {
			assert.Equal(t, validAppConfig.Id, id)
			assert.EqualValues(t, validAppConfig.Name, name.Name)
			assert.EqualValues(t, validAppConfig.Namespace, name.Namespace)
			return nil
		},
	}

	err := PostValidateAppConfig(ctx, appService, &environment.FakeService{})

	assert.NoError(t, err)
}

func Test_PostValidateAppConfig_InvalidAppName(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/validate/appconfig", safeMarshal(validAppConfig))
	req.Header.Add("CONTENT-TYPE", "application/json")
	ctx, _ := newContextWithRecorder(req)

	appService := &app.FakeService{
		CheckAppNameFn: func(id uuid.UUID, name *core.NamespacedName) error {
			return app.ErrInvalidAppName
		},
	}

	err := PostValidateAppConfig(ctx, appService, &environment.FakeService{})

	assert.Equal(t, app.ErrInvalidAppName, err)
}

// Mostly a dupe of Post test. In the future this should be factored out into its own service and stubbed from above tests
func Test_validateAppConfig(t *testing.T) {
	appService := &app.FakeService{
		CheckAppNameFn: func(id uuid.UUID, name *core.NamespacedName) error {
			assert.Equal(t, validAppConfig.Id, id)
			assert.EqualValues(t, validAppConfig.Name, name.Name)
			assert.EqualValues(t, validAppConfig.Namespace, name.Namespace)
			return nil
		},
	}

	result := validateAppConfig(validAppConfig, appService, &environment.FakeService{})

	assert.NoError(t, result)
}

func Test_validateAppConfig_InvalidAppName(t *testing.T) {
	appService := &app.FakeService{
		CheckAppNameFn: func(id uuid.UUID, name *core.NamespacedName) error {
			return app.ErrInvalidAppName
		},
	}

	result := validateAppConfig(validAppConfig, appService, &environment.FakeService{})

	assert.Equal(t, result, app.ErrInvalidAppName)
}

func Test_validateAppConfig_InvalidEnvOverride(t *testing.T) {
	appConfig := *validAppConfig
	appConfig.Overrides = map[string]model.OverrideableAppConfig{}
	appConfig.Overrides["prod"] = model.OverrideableAppConfig{}
	appConfig.Overrides["foo"] = model.OverrideableAppConfig{}

	appService := &app.FakeService{
		CheckAppNameFn: func(id uuid.UUID, name *core.NamespacedName) error {
			assert.Equal(t, validAppConfig.Id, id)
			assert.EqualValues(t, validAppConfig.Name, name.Name)
			assert.EqualValues(t, validAppConfig.Namespace, name.Namespace)
			return nil
		},
	}

	envService := &environment.FakeService{
		ValidateDeployableFn: func(envName string) error {
			if envName == "foo" {
				return errors.New("Invalid env")
			}

			return nil
		},
	}

	result := validateAppConfig(&appConfig, appService, envService)

	assert.IsType(t, &core.ValidationError{}, result)
	assert.EqualError(t, result, "Invalid environmentOverride: Invalid env")
}
