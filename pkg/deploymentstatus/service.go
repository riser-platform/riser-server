package deploymentstatus

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/stage"
)

type Service interface {
	Save(status *core.DeploymentStatus) error
	GetSummary(appName string) (*core.DeploymentStatusSummary, error)
}

type service struct {
	statuses     core.DeploymentStatusRepository
	stageService stage.Service
}

func NewService(statuses core.DeploymentStatusRepository, stageService stage.Service) Service {
	return &service{statuses, stageService}
}

func (s *service) GetSummary(appName string) (*core.DeploymentStatusSummary, error) {
	deploymentStatuses, err := s.statuses.FindByApp(appName)
	if err != nil {
		return nil, errors.Wrap(err, "Error retrieving deployment status")
	}

	summary := &core.DeploymentStatusSummary{
		DeploymentStatuses: deploymentStatuses,
		StageStatuses:      []core.StageStatus{},
	}

	stageMap := map[string]core.StageStatus{}

	for _, deploymentStatus := range deploymentStatuses {
		if _, ok := stageMap[deploymentStatus.StageName]; !ok {
			stageStatus, err := s.stageService.GetStatus(deploymentStatus.StageName)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("Error retrieving stage status for stage %q", deploymentStatus.StageName))
			}
			stageMap[deploymentStatus.StageName] = *stageStatus
			summary.StageStatuses = append(summary.StageStatuses, *stageStatus)
		}
	}

	return summary, nil
}

func (s *service) Save(status *core.DeploymentStatus) error {
	err := s.stageService.Ping(status.StageName)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error saving stage %q", status.StageName))
	}

	err = s.statuses.Save(status)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error saving deployment status for app %q", status.AppName))
	}

	return nil
}
