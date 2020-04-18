package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/app"
	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/stretchr/testify/assert"
)

var validAppConfig = &model.AppConfig{
	Name:      "myapp",
	Namespace: "myns",
	Id:        uuid.New(),
	Image:     "myimage",
	Expose: &model.AppConfigExpose{
		ContainerPort: 80,
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

	err := PostValidateAppConfig(ctx, appService)

	assert.NoError(t, err)
}

func Test_PostValidateAppConfig_InvalidAppName(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/validate/appconfig", safeMarshal(validAppConfig))
	req.Header.Add("CONTENT-TYPE", "application/json")
	ctx, _ := newContextWithRecorder(req)

	appService := &app.FakeService{
		CheckAppNameFn: func(id uuid.UUID, name *core.NamespacedName) error {
			assert.Equal(t, validAppConfig.Id, id)
			assert.EqualValues(t, validAppConfig.Name, name.Name)
			assert.EqualValues(t, validAppConfig.Namespace, name.Namespace)
			return app.ErrInvalidAppName
		},
	}

	err := PostValidateAppConfig(ctx, appService)

	assert.Equal(t, app.ErrInvalidAppName, err)
}
