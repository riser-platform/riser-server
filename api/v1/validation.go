package v1

import (
	"net/http"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/labstack/echo/v4"
	"github.com/riser-platform/riser-server/api/v1/model"
)

func PostValidateAppConfig(c echo.Context) error {
	appConfig := &model.AppConfigWithOverrides{}
	err := c.Bind(appConfig)
	if err == nil {
		err = appConfig.Validate()
	}

	if err != nil {
		return core.NewValidationError("Invalid app config", err)
	}

	return c.NoContent(http.StatusNoContent)
}
