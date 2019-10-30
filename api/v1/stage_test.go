package v1

import (
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_mapStageMetaFromDomain(t *testing.T) {
	domain := core.Stage{
		Name: "mystage",
	}

	result := mapStageMetaFromDomain(domain)

	assert.Equal(t, "mystage", result.Name)
}

func Test_mapStageMetaArrayFromDomain(t *testing.T) {
	domainArray := []core.Stage{
		core.Stage{Name: "mystage1"},
		core.Stage{Name: "mystage2"},
	}

	result := mapStageMetaArrayFromDomain(domainArray)

	assert.Len(t, result, 2)
	assert.Equal(t, "mystage1", result[0].Name)
	assert.Equal(t, "mystage2", result[1].Name)
}

func Test_mapStageConfigToDomain(t *testing.T) {
	config := &model.StageConfig{
		SealedSecretCert:  []byte{0x1},
		PublicGatewayHost: "myhost",
	}

	result := mapStageConfigToDomain(config)

	assert.Equal(t, []byte{0x1}, result.SealedSecretCert)
	assert.Equal(t, "myhost", result.PublicGatewayHost)
}

func Test_validateStageName_Error(t *testing.T) {
	result := validateStageName("")
	assert.NotNil(t, result)
	assert.IsType(t, &core.ValidationError{}, result)
}

func Test_validateStageName_NoError(t *testing.T) {
	result := validateStageName("valid")
	assert.Nil(t, result)
}
