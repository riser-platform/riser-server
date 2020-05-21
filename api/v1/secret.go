package v1

import (
	"net/http"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/environment"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/git"
	"github.com/riser-platform/riser-server/pkg/secret"
	"github.com/riser-platform/riser-server/pkg/state"
)

func PutSecret(c echo.Context, stateRepo git.Repo, secretService secret.Service, environmentService environment.Service) error {
	unsealedSecret := &model.UnsealedSecret{}
	err := c.Bind(unsealedSecret)
	if err != nil {
		return errors.Wrap(err, "Error binding secret")
	}

	err = environmentService.ValidateDeployable(unsealedSecret.Environment)
	if err != nil {
		return err
	}

	err = secretService.SealAndSave(
		unsealedSecret.PlainText,
		mapSecretMetaFromModel(&unsealedSecret.SecretMeta),
		state.NewGitCommitter(stateRepo))
	if err == core.ErrConflictNewerVersion {
		return echo.NewHTTPError(http.StatusConflict, "A newer revision of the secret was saved while attempting to save this secret. This is usually caused by a race condition due to another user saving the secret at the same time.")
	}

	return err
}

func GetSecrets(c echo.Context, secrets core.SecretMetaRepository, environmentService environment.Service) error {
	envName := c.Param("envName")
	namespace := c.Param("namespace")
	appName := c.Param("appName")

	err := environmentService.ValidateDeployable(envName)
	if err != nil {
		return err
	}

	secretMetas, err := secrets.ListByAppInEnvironment(core.NewNamespacedName(appName, namespace), envName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, mapSecretMetaStatusArrayFromDomain(secretMetas))
}

func mapSecretMetaStatusFromDomain(domain core.SecretMeta) model.SecretMetaStatus {
	return model.SecretMetaStatus{
		SecretMeta: model.SecretMeta{
			AppName:     model.AppName(domain.App.Name),
			Namespace:   model.NamespaceName(domain.App.Namespace),
			Environment: domain.EnvironmentName,
			Name:        domain.Name,
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
		App:             core.NewNamespacedName(string(in.AppName), string(in.Namespace)),
		Name:            in.Name,
		EnvironmentName: in.Environment,
	}
}
