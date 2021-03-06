package v1

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/riser-platform/riser-server/pkg/deployment"
	"github.com/riser-platform/riser-server/pkg/environment"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/app"

	model "github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/git"
	"github.com/riser-platform/riser-server/pkg/state"

	"github.com/labstack/echo/v4"
)

// TODO: Refactor and add unit test coverage
func PostDeployment(c echo.Context, repoCache *environment.RepoCache, appService app.Service, deploymentService deployment.Service, environmentService environment.Service) error {
	deploymentRequest := &model.SaveDeploymentRequest{}
	err := c.Bind(deploymentRequest)
	if err != nil {
		return err
	}

	isDryRun := c.QueryParam("dryRun") == "true"

	err = environmentService.ValidateDeployable(deploymentRequest.Environment)
	if err != nil {
		return err
	}

	newDeployment, err := mapDeploymentRequestToDomain(deploymentRequest)
	if err != nil {
		return err
	}

	err = appService.CheckID(deploymentRequest.App.AppConfig.Id, core.NewNamespacedName(string(deploymentRequest.App.Name), string(deploymentRequest.App.Namespace)))
	if err != nil {
		return err
	}

	var committer state.Committer

	if isDryRun {
		committer = state.NewDryRunCommitter()
	} else {
		gitRepo, err := repoCache.GetRepo(newDeployment.EnvironmentName)
		if err != nil {
			return err
		}
		committer = state.NewGitCommitter(gitRepo)
	}

	riserRevision, err := deploymentService.Update(newDeployment, committer, isDryRun)
	if err != nil {
		if err == git.ErrNoChanges {
			return c.JSON(http.StatusOK, model.SaveDeploymentResponse{Message: "No changes to deploy"})
		}
		return err
	}

	if isDryRun {
		dryRunCommitter := committer.(*state.DryRunCommitter)
		return c.JSON(http.StatusAccepted, model.SaveDeploymentResponse{

			Message:       "Dry run: changes not applied",
			DryRunCommits: mapDryRunCommitsFromDomain(dryRunCommitter.Commits),
		})
	}

	return c.JSON(http.StatusAccepted, model.SaveDeploymentResponse{RiserRevision: riserRevision, Message: "Deployment requested"})
}

func DeleteDeployment(c echo.Context, repoCache *environment.RepoCache, deploymentService deployment.Service) error {
	envName := c.Param("envName")
	gitRepo, err := repoCache.GetRepo(envName)
	if err != nil {
		return err
	}

	err = deploymentService.Delete(
		core.NewNamespacedName(c.Param("deploymentName"), c.Param("namespace")),
		envName,
		state.NewGitCommitter(gitRepo))

	if err != nil {
		if err == git.ErrNoChanges {
			return c.JSON(http.StatusNotFound, model.APIResponse{Message: "Deployment not found"})
		}
		return err
	}

	return c.JSON(http.StatusAccepted, model.APIResponse{Message: "Deployment deletion requested"})
}

func PutDeploymentStatus(c echo.Context, deployments core.DeploymentRepository) error {
	deploymentStatus := &model.DeploymentStatusMutable{}
	err := c.Bind(deploymentStatus)
	if err != nil {
		return errors.Wrap(err, "Error binding status")
	}

	deploymentName := c.Param("deploymentName")
	namespace := c.Param("namespace")
	envName := c.Param("envName")

	err = deployments.UpdateStatus(core.NewNamespacedName(deploymentName, namespace), envName, mapDeploymentStatusFromModel(deploymentStatus))
	if err == core.ErrConflictNewerVersion {
		return echo.NewHTTPError(http.StatusConflict, "A newer revision of the deployment has been observed or the deployment does not exist in this environment")
	}

	return err
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

func mapDeploymentRequestToDomain(deploymentRequest *model.SaveDeploymentRequest) (*core.DeploymentConfig, error) {
	app, err := deploymentRequest.App.ApplyOverrides(deploymentRequest.Environment)
	if err != nil {
		return nil, err
	}
	return &core.DeploymentConfig{
		Name:            deploymentRequest.Name,
		Namespace:       string(app.Namespace),
		EnvironmentName: deploymentRequest.Environment,
		Docker: core.DeploymentDocker{
			Tag: deploymentRequest.Docker.Tag,
		},
		App:           app,
		ManualRollout: deploymentRequest.ManualRollout,
	}, nil
}
