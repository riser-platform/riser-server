package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/riser-platform/riser-server/pkg/stage"

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
			AppId: uuid.New(),
			Stage: "dev",
			Name:  "mysecret",
		},
		PlainText: "myplain",
	}

	req := httptest.NewRequest(http.MethodPut, "/secrets/", safeMarshal(unsealed))
	req.Header.Add("CONTENT-TYPE", "application/json")

	ctx, rec := newContextWithRecorder(req)

	secretService := &secret.FakeService{
		SealAndSaveFn: func(plaintextSecret string, secretMeta *core.SecretMeta, namespace string, committer state.Committer) error {
			assert.Equal(t, "myplain", plaintextSecret)
			assert.Equal(t, secretMeta, mapSecretMetaFromModel(&unsealed.SecretMeta))
			assert.Equal(t, DefaultNamespace, namespace)
			return nil
		},
	}

	stageService := &stage.FakeService{
		ValidateDeployableFn: func(stageName string) error {
			assert.Equal(t, "dev", stageName)
			return nil
		},
	}

	err := PutSecret(ctx, nil, secretService, stageService)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Equal(t, 1, secretService.SealAndSaveCallCount)
}

func Test_PutSecret_WhenRevisionConflict(t *testing.T) {
	unsealed := model.UnsealedSecret{
		SecretMeta: model.SecretMeta{
			AppId: uuid.New(),
			Stage: "dev",
			Name:  "mysecret",
		},
		PlainText: "myplain",
	}

	req := httptest.NewRequest(http.MethodPut, "/secrets/", safeMarshal(unsealed))
	req.Header.Add("CONTENT-TYPE", "application/json")

	ctx, _ := newContextWithRecorder(req)

	secretService := &secret.FakeService{
		SealAndSaveFn: func(plaintextSecret string, secretMeta *core.SecretMeta, namespace string, committer state.Committer) error {
			return core.ErrConflictNewerVersion
		},
	}

	stageService := &stage.FakeService{
		ValidateDeployableFn: func(stageName string) error {
			return nil
		},
	}

	err := PutSecret(ctx, nil, secretService, stageService)
	require.IsType(t, &echo.HTTPError{}, err)
	httpErr := err.(*echo.HTTPError)
	assert.Equal(t, "A newer revision of the secret was saved while attempting to save this secret. This is usually caused by a race condition due to another user saving the secret at the same time.", httpErr.Message)
	assert.Equal(t, http.StatusConflict, httpErr.Code)
}

func Test_mapSecretMetaStatusFromDomain(t *testing.T) {
	domain := core.SecretMeta{
		AppId:     uuid.New(),
		StageName: "mystage",
		Name:      "mysecret",
		Revision:  1,
	}

	result := mapSecretMetaStatusFromDomain(domain)

	assert.Equal(t, domain.AppId, result.AppId)
	assert.Equal(t, "mystage", result.Stage)
	assert.Equal(t, "mysecret", result.Name)
	assert.EqualValues(t, 1, result.Revision)
}

func Test_mapSecretMetaStatusArrayFromDomain(t *testing.T) {
	domainArray := []core.SecretMeta{
		core.SecretMeta{Name: "secret1"},
		core.SecretMeta{Name: "secret2"},
	}

	result := mapSecretMetaStatusArrayFromDomain(domainArray)

	assert.Len(t, result, 2)
	assert.Equal(t, "secret1", result[0].Name)
	assert.Equal(t, "secret2", result[1].Name)
}

func Test_mapSecretMetaFromModel(t *testing.T) {
	model := &model.SecretMeta{
		AppId: uuid.New(),
		Name:  "mysecret",
		Stage: "mystage",
	}

	result := mapSecretMetaFromModel(model)

	assert.Equal(t, model.AppId, result.AppId)
	assert.Equal(t, "mysecret", result.Name)
	assert.Equal(t, "mystage", result.StageName)
}
