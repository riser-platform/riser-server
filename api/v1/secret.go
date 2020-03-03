package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/pkg/stage"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/git"
	"github.com/riser-platform/riser-server/pkg/secret"
	"github.com/riser-platform/riser-server/pkg/state"
)

func PutSecret(c echo.Context, stateRepo git.Repo, secretService secret.Service, stageService stage.Service) error {
	unsealedSecret := &model.UnsealedSecret{}
	err := c.Bind(unsealedSecret)
	if err != nil {
		return errors.Wrap(err, "Error binding secret")
	}

	err = stageService.ValidateDeployable(unsealedSecret.Stage)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = secretService.SealAndSave(
		unsealedSecret.PlainText,
		mapSecretMetaFromModel(&unsealedSecret.SecretMeta),
		// TODO(ns)
		core.DefaultNamespace,
		state.NewGitCommitter(stateRepo))
	if err == core.ErrConflictNewerVersion {
		return echo.NewHTTPError(http.StatusConflict, "A newer revision of the secret was saved while attempting to save this secret. This is usually caused by a race condition due to another user saving the secret at the same time.")
	}

	return err
}

func GetSecrets(c echo.Context, secretService secret.Service, stageService stage.Service) error {
	stageName := c.Param("stageName")
	// TODO(ns) use app name nd namespace
	appId := uuid.New()
	err := stageService.ValidateDeployable(stageName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	secretMetas, err := secretService.FindByStage(appId, stageName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, mapSecretMetaStatusArrayFromDomain(secretMetas))
}

func mapSecretMetaStatusFromDomain(domain core.SecretMeta) model.SecretMetaStatus {
	return model.SecretMetaStatus{
		SecretMeta: model.SecretMeta{
			AppId: domain.AppId,
			Stage: domain.StageName,
			Name:  domain.Name,
		},
		Revision: domain.Revision,
	}
}

func mapSecretMetaStatusArrayFromDomain(domainArray []core.SecretMeta) []model.SecretMetaStatus {
	statuses := []model.SecretMetaStatus{}
	for _, domain := range domainArray {
		statuses = append(statuses, mapSecretMetaStatusFromDomain(domain))
	}

	return statuses
}

func mapSecretMetaFromModel(in *model.SecretMeta) *core.SecretMeta {
	return &core.SecretMeta{
		AppId:     in.AppId,
		Name:      in.Name,
		StageName: in.Stage,
	}
}
