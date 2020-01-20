package api

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v3"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/labstack/echo/v4"
)

// Note: This model is shared between API versions. Therefore any change here is breaking for all API versions.
type ValidationErrorResponse struct {
	Message          string            `json:"message"`
	ValidationErrors map[string]string `json:"validationErrors,omitempty"`
}

func ErrorHandler(err error, c echo.Context) {
	var (
		code          = http.StatusInternalServerError
		internalError = err
		jsonResponse  interface{}
	)

	if httpError, ok := err.(*echo.HTTPError); ok {
		internalError = httpError.Internal
		code = httpError.Code
		jsonResponse = echo.Map{"message": httpError.Message}
	}

	if validationError, ok := err.(*core.ValidationError); ok {
		// An ozzo-validation Internal error means that something went wrong (e.g. a misconfigured validation rule).
		if ozzoInternal, ok := validationError.ValidationError.(validation.InternalError); ok {
			internalError = ozzoInternal
		} else {
			internalError = nil
			code = http.StatusBadRequest
			// ozzoValidation returns field specific errors, resulting in more user friendly error messages
			if ozzoValidation, ok := validationError.ValidationError.(validation.Errors); ok {
				jsonResponse = &ValidationErrorResponse{
					Message:          validationError.Message,
					ValidationErrors: formatValidationErrors(ozzoValidation),
				}
			} else {
				message := validationError.Message
				if validationError.ValidationError != nil {
					message = fmt.Sprintf("%s: %s", validationError.Message, validationError.ValidationError.Error())
				}
				jsonResponse = &ValidationErrorResponse{
					Message: message,
				}
			}
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

func formatValidationErrors(errors validation.Errors) map[string]string {
	result := map[string]string{}
	if errors == nil {
		return result
	}
	for k, v := range errors {
		result[k] = v.Error()
	}

	return result
}

func logErrorWithStack(err error, c echo.Context) {
	c.Logger().Error(fmt.Sprintf("%+v", err))
}
