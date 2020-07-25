package v1

import (
	"net/http"
	"time"

	"github.com/riser-platform/riser-server/pkg/app"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/deploymentstatus"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/labstack/echo/v4"
)

func GetAppStatus(c echo.Context, appService app.Service, statusService deploymentstatus.Service) error {
	domainApp, err := appService.GetByName(core.NewNamespacedName(c.Param("appName"), c.Param("namespace")))
	if err != nil {
		return err
	}

	appStatus, err := statusService.GetByApp(domainApp.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, mapAppStatusFromDomain(appStatus))
}

func mapAppStatusFromDomain(domain *core.AppStatus) *model.AppStatus {
	out := &model.AppStatus{
		Environments: mapEnvironmentStatusesFromDomain(domain.EnvironmentStatus),
		Deployments:  []model.DeploymentStatus{},
	}

	for _, deployment := range domain.Deployments {
		out.Deployments = append(out.Deployments, *mapDeploymentToStatusModel(&deployment))
	}

	return out
}

func mapEnvironmentStatusesFromDomain(domain []core.EnvironmentStatus) []model.EnvironmentStatus {
	out := []model.EnvironmentStatus{}
	for _, envStatus := range domain {
		out = append(out, model.EnvironmentStatus{
			EnvironmentName: envStatus.EnvironmentName,
			Healthy:         envStatus.Healthy,
			Reason:          envStatus.Reason,
		})
	}

	return out
}

func mapDeploymentToStatusModel(domain *core.Deployment) *model.DeploymentStatus {
	status := &model.DeploymentStatus{
		AppId:           domain.AppId,
		DeploymentName:  domain.Name,
		Namespace:       domain.Namespace,
		EnvironmentName: domain.EnvironmentName,
		RiserRevision:   domain.RiserRevision,
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
