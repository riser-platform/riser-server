package v1

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/labstack/echo/v4"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/stage"
)

// TODO: Once RBAC is implemented this should be limited to the controller.
func PostStagePing(c echo.Context, stageService stage.Service) error {
	stageName := c.Param("stageName")
	err := validateStageName(stageName)
	if err != nil {
		return err
	}
	return stageService.Ping(stageName)
}

func PutStageConfig(c echo.Context, stageService stage.Service) error {
	stageRequest := &model.StageConfig{}
	err := c.Bind(stageRequest)
	if err != nil {
		return err
	}

	stageName := c.Param("stageName")

	err = validateStageName(stageName)
	if err != nil {
		return err
	}

	err = stageService.SetConfig(stageName, mapStageConfigToDomain(stageRequest))
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusAccepted)
}

func ListStages(c echo.Context, stageRepository core.StageRepository) error {
	stages, err := stageRepository.List()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, mapStageMetaArrayFromDomain(stages))
}

func validateStageName(stageName string) error {
	err := validation.Validate(&stageName, model.RulesNamingIdentifier()...)
	if err != nil {
		return core.NewValidationError("invalid stage name", err)
	}
	return nil
}

func mapStageMetaFromDomain(domain core.Stage) model.StageMeta {
	return model.StageMeta{
		Name: domain.Name,
	}
}

func mapStageMetaArrayFromDomain(domainArray []core.Stage) []model.StageMeta {
	stages := []model.StageMeta{}
	for _, domain := range domainArray {
		stages = append(stages, mapStageMetaFromDomain(domain))
	}
	return stages
}

func mapStageConfigToDomain(in *model.StageConfig) *core.StageConfig {
	return &core.StageConfig{
		SealedSecretCert:  in.SealedSecretCert,
		PublicGatewayHost: in.PublicGatewayHost,
	}
}
