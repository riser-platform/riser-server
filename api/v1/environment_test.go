package v1

import (
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_mapEnvironmentMetaFromDomain(t *testing.T) {
	domain := core.Environment{
		Name: "myenv",
	}

	result := mapEnvironmentMetaFromDomain(domain)

	assert.Equal(t, "myenv", result.Name)
}

func Test_mapEnvironmentMetaArrayFromDomain(t *testing.T) {
	domainArray := []core.Environment{
		{Name: "myenv1"},
		{Name: "myenv2"},
	}

	result := mapEnvironmentMetaArrayFromDomain(domainArray)

	assert.Len(t, result, 2)
	assert.Equal(t, "myenv1", result[0].Name)
	assert.Equal(t, "myenv2", result[1].Name)
}

func Test_mapEnvironmentConfigToDomain(t *testing.T) {
	config := &model.EnvironmentConfig{
		SealedSecretCert:  []byte{0x1},
		PublicGatewayHost: "myhost",
	}

	result := mapEnvironmentConfigToDomain(config)

	assert.Equal(t, []byte{0x1}, result.SealedSecretCert)
	assert.Equal(t, "myhost", result.PublicGatewayHost)
}

func Test_mapEnvironmentConfigFromDomain(t *testing.T) {
	domain := &core.EnvironmentConfig{
		SealedSecretCert:  []byte{0x1},
		PublicGatewayHost: "myhost",
	}

	result := mapEnvironmentConfigFromDomain(domain)

	assert.Equal(t, []byte{0x1}, result.SealedSecretCert)
	assert.Equal(t, "myhost", result.PublicGatewayHost)
}

func Test_validateEnvironmentName_Error(t *testing.T) {
	result := validateEnvironmentName("")
	assert.NotNil(t, result)
	assert.Equal(t, "invalid environment name: cannot be blank", result.Error())
}

func Test_validateEnvironmentName_NoError(t *testing.T) {
	result := validateEnvironmentName("valid")
	assert.Nil(t, result)
}
