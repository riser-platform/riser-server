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

		status.Problems = make([]model.DeploymentStatusProblem, len(domain.Doc.Status.Problems))
		for idx, problem := range domain.Doc.Status.Problems {
			status.Problems[idx] = model.DeploymentStatusProblem{Count: problem.Count, Message: problem.Message}
		}

		status.Revisions = make([]model.DeploymentRevisionStatus, len(domain.Doc.Status.Revisions))
		for idx, revision := range domain.Doc.Status.Revisions {
			status.Revisions[idx] = model.DeploymentRevisionStatus{
				Name:                revision.Name,
				AvailableReplicas:   revision.AvailableReplicas,
				DockerImage:         revision.DockerImage,
				RiserGeneration:     revision.RiserGeneration,
				RolloutStatus:       revision.RolloutStatus,
				RolloutStatusReason: revision.RolloutStatusReason,
			}
		}

		status.Traffic = make([]model.DeploymentTrafficStatus, len(domain.Doc.Status.Traffic))
		for idx, traffic := range domain.Doc.Status.Traffic {
			status.Traffic[idx] = model.DeploymentTrafficStatus{
				Latest:       traffic.Latest,
				Percent:      traffic.Percent,
				RevisionName: traffic.RevisionName,
			}
		}
	}
	return status
}

func mapDeploymentStatusFromModel(in *model.DeploymentStatusMutable) *core.DeploymentStatus {
	out := &core.DeploymentStatus{
		ObservedRiserGeneration: in.ObservedRiserGeneration,
		LatestReadyRevisionName: in.LatestReadyRevisionName,
		LastUpdated:             time.Now().UTC(),
	}

	out.Problems = make([]core.DeploymentStatusProblem, len(in.Problems))
	for idx, problem := range in.Problems {
		out.Problems[idx] = core.DeploymentStatusProblem{Count: problem.Count, Message: problem.Message}
	}

	out.Revisions = make([]core.DeploymentRevisionStatus, len(in.Revisions))
	for idx, revision := range in.Revisions {
		out.Revisions[idx] = core.DeploymentRevisionStatus{
			Name:                revision.Name,
			AvailableReplicas:   revision.AvailableReplicas,
			DockerImage:         revision.DockerImage,
			RiserGeneration:     revision.RiserGeneration,
			RolloutStatus:       revision.RolloutStatus,
			RolloutStatusReason: revision.RolloutStatusReason,
		}
	}

	out.Traffic = make([]core.DeploymentTrafficStatus, len(in.Traffic))
	for idx, traffic := range in.Traffic {
		out.Traffic[idx] = core.DeploymentTrafficStatus{
			Latest:       traffic.Latest,
			Percent:      traffic.Percent,
			RevisionName: traffic.RevisionName,
		}
	}

	return out
}
