package v1

import (
	"testing"

	"github.com/labstack/echo/v4"

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

var stageNameTests = []struct {
	stageName string
	expected  error
}{
	{"", echo.NewHTTPError(400, "cannot be blank")},
	{"a", echo.NewHTTPError(400, "the length must be between 3 and 63")},
	{"0abcd", echo.NewHTTPError(400, "must be alphanumeric and start with a letter")},
	{"A123456789012345678901234567890123456789012345678901234567891234", echo.NewHTTPError(400, "the length must be between 3 and 63")},
	{"A!@#", echo.NewHTTPError(400, "must be alphanumeric and start with a letter")},
	{"valid", nil},
}

func Test_validateStageName(t *testing.T) {
	for _, tt := range stageNameTests {
		result := validateStageName(tt.stageName)
		assert.Equal(t, tt.expected, result, tt.stageName)
	}
}
