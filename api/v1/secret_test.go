package v1

import (
	"testing"
	"time"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_mapSecretMetaStatusFromDomain(t *testing.T) {
	domain := core.SecretMeta{
		AppName:    "myapp",
		StageName:  "mystage",
		SecretName: "mysecret",
		Doc: core.SecretMetaDoc{
			LastUpdated: time.Now(),
		},
	}

	result := mapSecretMetaStatusFromDomain(domain)

	assert.Equal(t, "myapp", result.App)
	assert.Equal(t, "mystage", result.Stage)
	assert.Equal(t, "mysecret", result.Name)
	assert.Equal(t, domain.Doc.LastUpdated, result.LastUpdated)
}

func Test_mapSecretMetaStatusArrayFromDomain(t *testing.T) {
	domainArray := []core.SecretMeta{
		core.SecretMeta{SecretName: "secret1"},
		core.SecretMeta{SecretName: "secret2"},
	}

	result := mapSecretMetaStatusArrayFromDomain(domainArray)

	assert.Len(t, result, 2)
	assert.Equal(t, "secret1", result[0].Name)
	assert.Equal(t, "secret2", result[1].Name)
}

func Test_mapSecretMetaFromModel(t *testing.T) {
	model := &model.SecretMeta{
		App:   "myapp",
		Name:  "mysecret",
		Stage: "mystage",
	}

	result := mapSecretMetaFromModel(model)

	assert.Equal(t, "myapp", result.AppName)
	assert.Equal(t, "mysecret", result.SecretName)
	assert.Equal(t, "mystage", result.StageName)
}
