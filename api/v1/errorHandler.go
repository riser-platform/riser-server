package v1

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func ApiErrorHandler(err error, c echo.Context) {
	// TODO: Revisit why we're using our own APIError and not echo's HTTPError
	if apiError, ok := err.(*APIError); ok {
		// Log the internal error if it exists.
		if apiError.Internal != nil {
			c.Logger().Error(fmt.Sprintf("%+v", apiError.Internal))
		}
		c.Echo().DefaultHTTPErrorHandler(echo.NewHTTPError(apiError.HTTPCode, apiError.Message), c)
	} else {
		if !c.Response().Committed {
			// TODO: Find a better way to add the stack trace to the log after the request failure log
			c.Logger().Error(fmt.Sprintf("%+v", err))
		}
		c.Echo().DefaultHTTPErrorHandler(err, c)
	}
}
