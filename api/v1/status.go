package v1

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/deploymentstatus"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/labstack/echo/v4"
)

func GetStatus(c echo.Context, statusService deploymentstatus.Service) error {
	appId, err := uuid.Parse(c.Param("appId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	appStatus, err := statusService.GetByApp(appId)
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
		DeploymentName: domain.Name,
		StageName:      domain.StageName,
		RiserRevision:  domain.RiserRevision,
	}
	if domain.Doc.Status == nil {
		status.DeploymentStatusMutable = model.DeploymentStatusMutable{}
	} else {
		status.DeploymentStatusMutable = model.DeploymentStatusMutable{
			ObservedRiserRevision:     domain.Doc.Status.ObservedRiserRevision,
			LatestCreatedRevisionName: domain.Doc.Status.LatestCreatedRevisionName,
			LatestReadyRevisionName:   domain.Doc.Status.LatestReadyRevisionName,
		}

		status.Revisions = make([]model.DeploymentRevisionStatus, len(domain.Doc.Status.Revisions))
		for idx, revision := range domain.Doc.Status.Revisions {
			status.Revisions[idx] = model.DeploymentRevisionStatus{
				Name:                 revision.Name,
				AvailableReplicas:    revision.AvailableReplicas,
				DockerImage:          revision.DockerImage,
				RiserRevision:        revision.RiserRevision,
				RevisionStatus:       revision.RevisionStatus,
				RevisionStatusReason: revision.RevisionStatusReason,
			}
		}

		status.Traffic = make([]model.DeploymentTrafficStatus, len(domain.Doc.Status.Traffic))
		for idx, traffic := range domain.Doc.Status.Traffic {
			status.Traffic[idx] = model.DeploymentTrafficStatus{
				Percent:      traffic.Percent,
				RevisionName: traffic.RevisionName,
				Tag:          traffic.Tag,
			}
		}
	}
	return status
}

func mapDeploymentStatusFromModel(in *model.DeploymentStatusMutable) *core.DeploymentStatus {
	out := &core.DeploymentStatus{
		ObservedRiserRevision:     in.ObservedRiserRevision,
		LatestCreatedRevisionName: in.LatestCreatedRevisionName,
		LatestReadyRevisionName:   in.LatestReadyRevisionName,
		LastUpdated:               time.Now().UTC(),
	}

	out.Revisions = make([]core.DeploymentRevisionStatus, len(in.Revisions))
	for idx, revision := range in.Revisions {
		out.Revisions[idx] = core.DeploymentRevisionStatus{
			Name:                 revision.Name,
			AvailableReplicas:    revision.AvailableReplicas,
			DockerImage:          revision.DockerImage,
			RiserRevision:        revision.RiserRevision,
			RevisionStatus:       revision.RevisionStatus,
			RevisionStatusReason: revision.RevisionStatusReason,
		}

	}

	out.Traffic = make([]core.DeploymentTrafficStatus, len(in.Traffic))
	for idx, traffic := range in.Traffic {
		out.Traffic[idx] = core.DeploymentTrafficStatus{
			Percent:      traffic.Percent,
			RevisionName: traffic.RevisionName,
			Tag:          traffic.Tag,
		}
	}

	return out
}
