package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/login"
)

func loginWithApiKey(c echo.Context, loginService login.Service, apikey string) (bool, error) {
	username, err := loginService.LoginWithApiKey(apikey)
	if err != nil {
		if err == login.ErrInvalidLogin {
			return false, nil
		}

		return false, errors.Wrap(err, "Error logging in with API key")
	}
	c.Set("username", username)
	return true, nil
}
