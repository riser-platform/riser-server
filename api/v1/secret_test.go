package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/riser-platform/riser-server/pkg/environment"
	"github.com/riser-platform/riser-server/pkg/secret"
	"github.com/riser-platform/riser-server/pkg/state"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PutSecret(t *testing.T) {
	unsealed := model.UnsealedSecret{
		SecretMeta: model.SecretMeta{
			AppName:     "myapp",
			Namespace:   "myns",
			Environment: "dev",
			Name:        "mysecret",
		},
		PlainText: "myplain",
	}

	req := httptest.NewRequest(http.MethodPut, "/secrets/", safeMarshal(unsealed))
	req.Header.Add("CONTENT-TYPE", "application/json")

	ctx, rec := newContextWithRecorder(req)

	secretService := &secret.FakeService{
		SealAndSaveFn: func(plaintextSecret string, secretMeta *core.SecretMeta, committer state.Committer) error {
			assert.Equal(t, "myplain", plaintextSecret)
			assert.Equal(t, secretMeta, mapSecretMetaFromModel(&unsealed.SecretMeta))
			return nil
		},
	}

	environmentService := &environment.FakeService{
		ValidateDeployableFn: func(envName string) error {
			assert.Equal(t, "dev", envName)
			return nil
		},
	}

	err := PutSecret(ctx, nil, secretService, environmentService)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Equal(t, 1, secretService.SealAndSaveCallCount)
}

func Test_PutSecret_WhenRevisionConflict(t *testing.T) {
	unsealed := model.UnsealedSecret{
		SecretMeta: model.SecretMeta{
			AppName:     "myapp",
			Namespace:   "myns",
			Environment: "dev",
			Name:        "mysecret",
		},
		PlainText: "myplain",
	}

	req := httptest.NewRequest(http.MethodPut, "/secrets/", safeMarshal(unsealed))
	req.Header.Add("CONTENT-TYPE", "application/json")

	ctx, _ := newContextWithRecorder(req)

	secretService := &secret.FakeService{
		SealAndSaveFn: func(plaintextSecret string, secretMeta *core.SecretMeta, committer state.Committer) error {
			return core.ErrConflictNewerVersion
		},
	}

	environmentService := &environment.FakeService{
		ValidateDeployableFn: func(envName string) error {
			return nil
		},
	}

	err := PutSecret(ctx, nil, secretService, environmentService)
	require.IsType(t, &echo.HTTPError{}, err)
	httpErr := err.(*echo.HTTPError)
	assert.Equal(t, "A newer revision of the secret was saved while attempting to save this secret. This is usually caused by a race condition due to another user saving the secret at the same time.", httpErr.Message)
	assert.Equal(t, http.StatusConflict, httpErr.Code)
}

func Test_mapSecretMetaStatusFromDomain(t *testing.T) {
	domain := core.SecretMeta{
		App:             core.NewNamespacedName("myapp", "myns"),
		EnvironmentName: "myenv",
		Name:            "mysecret",
		Revision:        1,
	}

	result := mapSecretMetaStatusFromDomain(domain)

	assert.EqualValues(t, "myapp", result.AppName)
	assert.EqualValues(t, "myns", result.Namespace)
	assert.Equal(t, "myenv", result.Environment)
	assert.Equal(t, "mysecret", result.Name)
	assert.EqualValues(t, 1, result.Revision)
}

func Test_mapSecretMetaStatusArrayFromDomain(t *testing.T) {
	domainArray := []core.SecretMeta{
		{Name: "secret1", App: &core.NamespacedName{}},
		{Name: "secret2", App: &core.NamespacedName{}},
	}

	result := mapSecretMetaStatusArrayFromDomain(domainArray)

	assert.Len(t, result, 2)
	assert.Equal(t, "secret1", result[0].Name)
	assert.Equal(t, "secret2", result[1].Name)
}

func Test_mapSecretMetaFromModel(t *testing.T) {
	model := &model.SecretMeta{
		AppName:     "myapp",
		Namespace:   "myns",
		Name:        "mysecret",
		Environment: "myenv",
	}

	result := mapSecretMetaFromModel(model)

	assert.Equal(t, "myapp", result.App.Name)
	assert.Equal(t, "myns", result.App.Namespace)
	assert.Equal(t, "mysecret", result.Name)
	assert.Equal(t, "myenv", result.EnvironmentName)
}
