package v1

import (
	"net/http"

	"github.com/riser-platform/riser-server/pkg/namespace"

	"github.com/riser-platform/riser-server/pkg/deployment"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/app"

	model "github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/git"
	"github.com/riser-platform/riser-server/pkg/state"

	"github.com/labstack/echo/v4"
)

func PostDeployment(c echo.Context, stateRepo git.GitRepoProvider, appService app.Service, deploymentService deployment.Service) error {
	deploymentRequest := &model.DeploymentRequest{}
	err := c.Bind(deploymentRequest)
	if err != nil {
		return err
	}

	isDryRun := c.QueryParam("dryRun") == "true"

	err = deploymentRequest.App.Validate()
	if err != nil {
		return handleValidationError(c, err, "Invalid app config")
	}

	appId, err := core.DecodeAppId(deploymentRequest.App.AppConfig.Id)
	if err != nil {
		return NewAPIError(http.StatusBadRequest, "App Id must be a hex string")
	}
	err = appService.CheckAppId(deploymentRequest.App.Name, appId)
	if err == app.ErrInvalidAppId {
		return NewAPIError(http.StatusBadRequest, "Invalid App Id")
	}
	if err != nil {
		return err
	}

	var committer state.Committer

	if isDryRun {
		committer = state.NewDryRunComitter()
	} else {
		committer = state.NewGitComitter(stateRepo)
	}

	// TODO: This is a hack that exists for ease of use since we only support the "apps" namespace.
	// Once we support multiple namespace this should be in its own route
	namespaceService := namespace.NewService()
	err = namespaceService.Save(&core.Namespace{Name: "apps", Stage: deploymentRequest.Stage}, committer)
	if err != nil && err != git.ErrNoChanges {
		return err
	}

	err = deploymentService.Update(mapDeploymentRequestToDomain(deploymentRequest), committer)
	if err != nil {
		if err == git.ErrNoChanges {
			return c.JSON(http.StatusOK, model.DeploymentResponse{Message: "No changes to deploy"})
		}
		return err
	}

	if isDryRun {
		dryRunCommitter := committer.(*state.DryRunComitter)

		return c.JSON(http.StatusAccepted, model.DeploymentResponse{
			Message:       "Dry run: changes not applied",
			DryRunCommits: mapDryRunCommitsFromDomain(dryRunCommitter.Commits),
		})
	}

	return c.JSON(http.StatusAccepted, model.APIResponse{Message: "Deployment requested"})
}

func mapDryRunCommitsFromDomain(commits []state.DryRunCommit) []model.DryRunCommit {
	out := []model.DryRunCommit{}
	for _, commit := range commits {
		modelCommit := model.DryRunCommit{}
		modelCommit.Message = commit.Message
		modelCommit.Files = []model.DryRunFile{}
		for _, file := range commit.Files {
			modelCommit.Files = append(modelCommit.Files, model.DryRunFile{Name: file.Name, Contents: string(file.Contents)})
		}
		out = append(out, modelCommit)
	}

	return out
}

func mapDeploymentRequestToDomain(deploymentRequest *model.DeploymentRequest) *core.NewDeployment {
	return &core.NewDeployment{
		DeploymentMeta: core.DeploymentMeta{
			Name:      deploymentRequest.Name,
			Namespace: deploymentRequest.Namespace,
			Stage:     deploymentRequest.Stage,
			Docker: core.DeploymentDocker{
				Tag: deploymentRequest.Docker.Tag,
			},
		},
		App: deploymentRequest.App,
	}
}
