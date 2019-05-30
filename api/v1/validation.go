package v1

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"

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
		return handleValidationError(c, err, "Invalid app config")
	}

	return c.NoContent(http.StatusNoContent)
}

// TODO: Consider pros/cons of using echo validation middleware
func handleValidationError(c echo.Context, inErr error, message string) error {
	if _, ok := inErr.(validation.Errors); ok {
		response := model.ValidationResponse{
			APIResponse: model.APIResponse{
				Message: message,
			},
			ValidationErrors: inErr,
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	return &APIError{
		HTTPCode: http.StatusBadRequest,
		Message:  fmt.Sprintf("%s: %s", message, inErr),
	}
}
