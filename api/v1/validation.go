package v1

import (
	"net/http"

	"github.com/riser-platform/riser-server/pkg/environment"

	"github.com/riser-platform/riser-server/pkg/app"

	"github.com/labstack/echo/v4"
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
)

func PostValidateAppConfig(c echo.Context, appService app.Service, environmentService environment.Service) error {
	appConfig := &model.AppConfigWithOverrides{}
	err := c.Bind(appConfig)
	// if err == nil {
	// 	// TODO: Do we need this or do we get this for free?
	// 	err = appConfig.Validate()
	// }

	if err != nil {
		return err
	}

	err = validateAppConfig(appConfig, appService, environmentService)
	if err == nil {
		return c.NoContent(http.StatusNoContent)
	}
	return err
}

// validateAppConfig performs additional validation beyond type validation of the appConfig model (i.e. appConfig.Validate())
// such as checking the database if the name is valid
// This should move into its own package once we move AppConfig outside of v1/model
func validateAppConfig(appConfig *model.AppConfigWithOverrides, appService app.Service, environmentService environment.Service) error {
	err := appService.CheckID(appConfig.Id, core.NewNamespacedName(string(appConfig.Name), string(appConfig.Namespace)))
	if err != nil {
		return err
	}
	for env := range appConfig.Overrides {
		err = environmentService.ValidateDeployable(env)
		if err != nil {
			return core.NewValidationError("Invalid environmentOverride", err)
		}
	}
	return nil
}
