package deploymentstatus

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/stage"
)

// TODO: Consider better homes for these
type Service interface {
	GetByApp(appId uuid.UUID) (*core.AppStatus, error)
}

type service struct {
	deployments  core.DeploymentRepository
	stageService stage.Service
}

func NewService(deployments core.DeploymentRepository, stageService stage.Service) Service {
	return &service{deployments, stageService}
}

func (s *service) GetByApp(appId uuid.UUID) (*core.AppStatus, error) {
	deployments, err := s.deployments.FindByApp(appId)
	if err != nil {
		return nil, errors.Wrap(err, "Error retrieving deployment status")
	}

	appStatus := &core.AppStatus{
		Deployments:   deployments,
		StageStatuses: []core.StageStatus{},
	}

	stageMap := map[string]core.StageStatus{}

	for _, deploymentStatus := range deployments {
		if _, ok := stageMap[deploymentStatus.StageName]; !ok {
			stageStatus, err := s.stageService.GetStatus(deploymentStatus.StageName)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("Error retrieving stage status for stage %q", deploymentStatus.StageName))
			}
			stageMap[deploymentStatus.StageName] = *stageStatus
			appStatus.StageStatuses = append(appStatus.StageStatuses, *stageStatus)
		}
	}

	return appStatus, nil
}
