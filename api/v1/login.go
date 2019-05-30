package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/riser-platform/riser-server/pkg/login"
)

func loginWithApiKey(c echo.Context, loginService login.Service, apikey string) (bool, error) {
	username, err := loginService.LoginWithApiKey(apikey)
	if err != nil {
		if err == login.ErrInvalidLogin {
			return false, nil
		}

		// Echo does not seem to log the returned error so log it here
		c.Logger().Errorf("Error logging in with API key: %s", err)
		return false, err
	}
	c.Set("username", username)
	return true, nil
}
