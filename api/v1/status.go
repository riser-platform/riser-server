package v1

import (
	"net/http"
	"time"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/deploymentstatus"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/labstack/echo/v4"
)

func GetStatus(c echo.Context, statusService deploymentstatus.Service) error {
	appName := c.Param("appName")

	appStatus, err := statusService.GetByApp(appName)
	if err != nil {
		return err
	}

	statusModel := model.AppStatus{
		Stages:      []model.StageStatus{},
		Deployments: []model.DeploymentStatus{},
	}

	// TODO: Move and test model conversion.
	for _, stageStatus := range appStatus.StageStatuses {
		statusModel.Stages = append(statusModel.Stages, model.StageStatus{
			StageName: stageStatus.StageName,
			Healthy:   stageStatus.Healthy,
			Reason:    stageStatus.Reason,
		})
	}

	for _, deployment := range appStatus.Deployments {
		statusModel.Deployments = append(statusModel.Deployments, *mapDeploymentToStatusModel(&deployment))
	}

	return c.JSON(http.StatusOK, statusModel)
}

func mapDeploymentToStatusModel(domain *core.Deployment) *model.DeploymentStatus {
	status := &model.DeploymentStatus{
		DeploymentName:  domain.Name,
		StageName:       domain.StageName,
		RiserGeneration: domain.RiserGeneration,
	}
	if domain.Doc.Status == nil {
		status.DeploymentStatusMutable = model.DeploymentStatusMutable{}
	} else {
		status.DeploymentStatusMutable = model.DeploymentStatusMutable{
			ObservedRiserGeneration: domain.Doc.Status.ObservedRiserGeneration,
			LatestReadyRevisionName: domain.Doc.Status.LatestReadyRevisionName,
		}

		status.Problems = []model.DeploymentStatusProblem{}
		for _, problem := range domain.Doc.Status.Problems {
			status.Problems = append(status.Problems, model.DeploymentStatusProblem{Count: problem.Count, Message: problem.Message})
		}

		status.Revisions = []model.DeploymentRevision{}
		for _, revision := range domain.Doc.Status.Revisions {
			status.Revisions = append(status.Revisions, model.DeploymentRevision{
				Name:                revision.Name,
				DockerImage:         revision.DockerImage,
				RiserGeneration:     revision.RiserGeneration,
				RolloutStatus:       revision.RolloutStatus,
				RolloutStatusReason: revision.RolloutStatusReason,
			})
		}
	}
	return status
}

func mapDeploymentStatusFromModel(in *model.DeploymentStatusMutable) *core.DeploymentStatus {
	out := &core.DeploymentStatus{
		ObservedRiserGeneration: in.ObservedRiserGeneration,
		LatestReadyRevisionName: "rev2",
		LastUpdated:             time.Now().UTC(),
	}

	out.Problems = []core.DeploymentStatusProblem{}
	for _, problem := range in.Problems {
		out.Problems = append(out.Problems, core.DeploymentStatusProblem{Count: problem.Count, Message: problem.Message})
	}

	out.Revisions = []core.DeploymentRevision{}
	for _, revision := range in.Revisions {
		out.Revisions = append(out.Revisions, core.DeploymentRevision{
			Name:                revision.Name,
			DockerImage:         revision.DockerImage,
			RiserGeneration:     revision.RiserGeneration,
			RolloutStatus:       revision.RolloutStatus,
			RolloutStatusReason: revision.RolloutStatusReason,
		})
	}

	return out
}
