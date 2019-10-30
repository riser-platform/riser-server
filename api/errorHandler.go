package api

import (
	"fmt"
	"net/http"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/labstack/echo/v4"
)

// Note: This model is shared between API versions. Therefore any change here is breaking for all API versions.
type ValidationErrorResponse struct {
	Message          string            `json:"message"`
	ValidationErrors map[string]string `json:"validationErrors"`
}

func ErrorHandler(err error, c echo.Context) {
	var (
		code          = http.StatusInternalServerError
		jsonResponse  interface{}
		internalError = err
	)

	if httpError, ok := err.(*echo.HTTPError); ok {
		internalError = httpError.Internal
		code = httpError.Code
		jsonResponse = echo.Map{"message": httpError.Message}
	}

	if validationError, ok := err.(*core.ValidationError); ok {
		internalError = validationError.Internal
		code = http.StatusBadRequest
		jsonResponse = &ValidationErrorResponse{
			Message:          validationError.Message,
			ValidationErrors: formatValidationErrors(validationError.ValidationErrors),
		}
	}

	// Checking Response().Committed is required to prevent duplicate log entries
	// I could not figure out a way to repro this in a unit test so tests will still pass if removed
	if !c.Response().Committed {
		if internalError != nil {
			logErrorWithStack(internalError, c)
		}

		if jsonResponse == nil {
			jsonResponse = echo.Map{"message": http.StatusText(code)}
		}

		err = c.JSON(code, jsonResponse)
		if err != nil {
			c.Logger().Error(err)
		}
	}
}

func formatValidationErrors(errors map[string]error) map[string]string {
	result := map[string]string{}
	for k, v := range errors {
		result[k] = v.Error()
	}

	return result
}

func logErrorWithStack(err error, c echo.Context) {
	c.Logger().Error(fmt.Sprintf("%+v", err))
}
