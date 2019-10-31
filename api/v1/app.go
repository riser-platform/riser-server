package v1

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/app"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/api/v1/model"
)

func PostApp(c echo.Context, appService app.Service) error {
	newAppRequest := &model.NewApp{}
	err := c.Bind(newAppRequest)
	if err != nil {
		return errors.Wrap(err, "Error binding App")
	}

	err = validation.Validate(&newAppRequest.Name, model.RulesAppName()...)
	if err != nil {
		return core.NewValidationError("invalid app name", err)
	}

	createdApp, err := appService.CreateApp(newAppRequest.Name)
	if err != nil {
		if err == app.ErrAlreadyExists {
			return echo.NewHTTPError(http.StatusConflict, fmt.Sprintf("The app \"%s\" already exists", newAppRequest.Name))
		} else {
			return errors.Wrap(err, "Error creating App")
		}
	}
	return c.JSON(http.StatusCreated, mapAppFromDomain(*createdApp))
}

func ListApps(c echo.Context, appRepo core.AppRepository) error {
	apps, err := appRepo.ListApps()
	if err != nil {
		return err
	}
	return c.JSON(200, mapAppArrayFromDomain(apps))
}

func mapAppFromDomain(domain core.App) model.App {
	return model.App{
		Name: domain.Name,
		Id:   domain.Hashid.String(),
	}
}

func mapAppArrayFromDomain(domainArray []core.App) []model.App {
	apps := []model.App{}
	for _, domain := range domainArray {
		apps = append(apps, mapAppFromDomain(domain))
	}

	return apps
}
