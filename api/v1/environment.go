package v1

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v3"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/labstack/echo/v4"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/environment"
)

// TODO: Once RBAC is implemented this should be limited to the controller.
func PostEnvironmentPing(c echo.Context, environmentService environment.Service) error {
	envName := c.Param("envName")
	err := validateEnvironmentName(envName)
	if err != nil {
		return err
	}
	return environmentService.Ping(envName)
}

func PutEnvironmentConfig(c echo.Context, environmentService environment.Service) error {
	environmentConfig := &model.EnvironmentConfig{}
	err := c.Bind(environmentConfig)
	if err != nil {
		return err
	}

	envName := c.Param("envName")

	err = validateEnvironmentName(envName)
	if err != nil {
		return err
	}

	err = environmentService.SetConfig(envName, mapEnvironmentConfigToDomain(environmentConfig))
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}

func ListEnvironments(c echo.Context, environmentRepository core.EnvironmentRepository) error {
	environments, err := environmentRepository.List()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, mapEnvironmentMetaArrayFromDomain(environments))
}

func validateEnvironmentName(envName string) error {
	rules := model.RulesNamingIdentifier()
	rules = append(rules, validation.Required)
	err := validation.Validate(&envName, rules...)
	if err != nil {
		return core.NewValidationError("invalid environment name", err)
	}
	return nil
}

func mapEnvironmentMetaFromDomain(domain core.Environment) model.EnvironmentMeta {
	return model.EnvironmentMeta{
		Name: domain.Name,
	}
}

func mapEnvironmentMetaArrayFromDomain(domainArray []core.Environment) []model.EnvironmentMeta {
	environments := []model.EnvironmentMeta{}
	for _, domain := range domainArray {
		environments = append(environments, mapEnvironmentMetaFromDomain(domain))
	}
	return environments
}

func mapEnvironmentConfigToDomain(in *model.EnvironmentConfig) *core.EnvironmentConfig {
	return &core.EnvironmentConfig{
		SealedSecretCert:  in.SealedSecretCert,
		PublicGatewayHost: in.PublicGatewayHost,
	}
}
