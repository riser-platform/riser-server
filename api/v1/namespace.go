package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/git"
	"github.com/riser-platform/riser-server/pkg/namespace"
	"github.com/riser-platform/riser-server/pkg/state"
)

func PostNamespace(c echo.Context, namespaceService namespace.Service, repo git.Repo) error {
	ns := &model.Namespace{}
	err := c.Bind(ns)
	if err != nil {
		return nil
	}

	return namespaceService.Create(ns.Name, state.NewGitCommitter(repo))
}
