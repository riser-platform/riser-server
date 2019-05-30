package v1

import (
	"net/http"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/git"
	"github.com/riser-platform/riser-server/pkg/secret"
	"github.com/riser-platform/riser-server/pkg/state"
)

func PutSecret(c echo.Context, stateRepo git.GitRepoProvider, secretService secret.Service) error {
	unsealedSecret := &model.UnsealedSecret{}
	err := c.Bind(unsealedSecret)
	if err != nil {
		return errors.Wrap(err, "Error binding secret")
	}

	// We don't know what namespace an app is associated with yet as namespace support is not fully supported.
	// This could get tricky if we support apps being deployed to mulitple namespaces.
	return secretService.SealAndSave(unsealedSecret.PlainText, mapSecretMetaFromModel(&unsealedSecret.SecretMeta), "apps", state.NewGitComitter(stateRepo))
}

func GetSecrets(c echo.Context, secretService secret.Service) error {
	appName := c.Param("appName")
	stageName := c.Param("stageName")

	secretMetas, err := secretService.FindByStage(appName, stageName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, mapSecretMetaStatusArrayFromDomain(secretMetas))
}

func mapSecretMetaStatusFromDomain(domain core.SecretMeta) model.SecretMetaStatus {
	return model.SecretMetaStatus{
		SecretMeta: model.SecretMeta{
			App:   domain.AppName,
			Stage: domain.StageName,
			Name:  domain.SecretName,
		},
		LastUpdated: domain.Doc.LastUpdated,
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
		AppName:    in.App,
		SecretName: in.Name,
		StageName:  in.Stage,
	}
}
