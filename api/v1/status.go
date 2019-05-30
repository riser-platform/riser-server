package v1

import (
	"net/http"
	"time"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/deploymentstatus"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func PostStatus(c echo.Context, deploymentStatusService deploymentstatus.Service) error {
	deploymentStatus := &model.DeploymentStatus{}
	err := c.Bind(deploymentStatus)
	if err != nil {
		return errors.Wrap(err, "Error binding status")
	}

	return deploymentStatusService.Save(mapDeploymentStatusFromModel(deploymentStatus))
}

func GetStatus(c echo.Context, deploymentStatusService deploymentstatus.Service) error {
	// TODO: Support by stage and deploymentName too
	appName := c.QueryParam("app")

	statusSummary, err := deploymentStatusService.GetSummary(appName)
	if err != nil {
		return err
	}

	statusModel := model.Status{
		Stages:      []model.StageStatus{},
		Deployments: []model.DeploymentStatus{},
	}

	// TODO: Move and test model conversion.
	for _, stageStatus := range statusSummary.StageStatuses {
		statusModel.Stages = append(statusModel.Stages, model.StageStatus{
			StageName: stageStatus.StageName,
			Healthy:   stageStatus.Healthy,
			Reason:    stageStatus.Reason,
		})
	}

	for _, deploymentStatus := range statusSummary.DeploymentStatuses {
		statusModel.Deployments = append(statusModel.Deployments, *mapDeploymentStatusToModel(&deploymentStatus))
	}

	return c.JSON(http.StatusOK, statusModel)
}

func mapDeploymentStatusToModel(domain *core.DeploymentStatus) *model.DeploymentStatus {
	status := &model.DeploymentStatus{
		AppName:             domain.AppName,
		DeploymentName:      domain.DeploymentName,
		StageName:           domain.StageName,
		RolloutStatus:       domain.Doc.RolloutStatus,
		RolloutStatusReason: domain.Doc.RolloutStatusReason,
		RolloutRevision:     domain.Doc.RolloutRevision,
		DockerImage:         domain.Doc.DockerImage,
	}

	status.Problems = []model.DeploymentStatusProblem{}
	for _, problem := range domain.Doc.Problems {
		status.Problems = append(status.Problems, model.DeploymentStatusProblem{Count: problem.Count, Message: problem.Message})
	}
	return status
}

func mapDeploymentStatusFromModel(in *model.DeploymentStatus) *core.DeploymentStatus {
	out := &core.DeploymentStatus{
		AppName:        in.AppName,
		DeploymentName: in.DeploymentName,
		StageName:      in.StageName,
		Doc: &core.DeploymentStatusDoc{
			RolloutStatus:       in.RolloutStatus,
			RolloutStatusReason: in.RolloutStatusReason,
			RolloutRevision:     in.RolloutRevision,
			DockerImage:         in.DockerImage,
			LastUpdated:         time.Now().UTC(),
		},
	}

	out.Doc.Problems = []core.DeploymentStatusProblem{}
	for _, problem := range in.Problems {
		out.Doc.Problems = append(out.Doc.Problems, core.DeploymentStatusProblem{Count: problem.Count, Message: problem.Message})
	}

	return out
}
