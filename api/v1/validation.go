package v1

import (
	"net/http"

	"github.com/riser-platform/riser-server/pkg/app"

	"github.com/labstack/echo/v4"
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
)

func PostValidateAppConfig(c echo.Context, appService app.Service) error {
	appConfig := &model.AppConfigWithOverrides{}
	err := c.Bind(appConfig)
	if err == nil {
		err = appConfig.Validate()
	}

	if err != nil {
		return err
	}

	err = appService.CheckAppName(appConfig.Id, core.NewNamespacedName(string(appConfig.Name), string(appConfig.Namespace)))
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
