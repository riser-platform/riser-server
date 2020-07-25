package v1

import (
	"net/http"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/app"

	"github.com/labstack/echo/v4"
	"github.com/riser-platform/riser-server/api/v1/model"
)

func PostApp(c echo.Context, appService app.Service) error {
	newAppRequest := &model.NewApp{}
	err := c.Bind(newAppRequest)
	if err != nil {
		return err
	}

	createdApp, err := appService.Create(core.NewNamespacedName(string(newAppRequest.Name), string(newAppRequest.Namespace)))
	if err != nil {
		return err

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

func GetApp(c echo.Context, apps core.AppRepository) error {
	domainApp, err := apps.GetByName(core.NewNamespacedName(c.Param("appName"), c.Param("namespace")))

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, mapAppFromDomain(*domainApp))
}

func mapAppFromDomain(domain core.App) model.App {
	return model.App{
		Id:        domain.Id,
		Name:      model.AppName(domain.Name),
		Namespace: model.NamespaceName(domain.Namespace),
	}
}

func mapAppArrayFromDomain(domainArray []core.App) []model.App {
	apps := []model.App{}
	for _, domain := range domainArray {
		apps = append(apps, mapAppFromDomain(domain))
	}

	return apps
}
