package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/git"
	"github.com/riser-platform/riser-server/pkg/namespace"
	"github.com/riser-platform/riser-server/pkg/state"
)

func PostNamespace(c echo.Context, namespaceService namespace.Service, repo git.Repo) error {
	ns := &model.Namespace{}
	err := c.Bind(ns)
	if err != nil {
		return err
	}

	return namespaceService.Create(string(ns.Name), state.NewGitCommitter(repo))
}

func GetNamespaces(c echo.Context, namespaces core.NamespaceRepository) error {
	domainArray, err := namespaces.List()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, mapNamespaceArrayFromDomain(domainArray))
}

func mapNamespaceArrayFromDomain(domainArray []core.Namespace) []model.Namespace {
	modelArray := []model.Namespace{}
	for _, domain := range domainArray {
		modelArray = append(modelArray, mapNamespaceFromDomain(domain))
	}

	return modelArray
}

func mapNamespaceFromDomain(domain core.Namespace) model.Namespace {
	return model.Namespace{Name: model.NamespaceName(domain.Name)}
}
