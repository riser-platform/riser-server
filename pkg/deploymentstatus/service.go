package deploymentstatus

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/environment"
)

// TODO: Consider better homes for these
type Service interface {
	GetByApp(appId uuid.UUID) (*core.AppStatus, error)
}

type service struct {
	deployments core.DeploymentRepository
	envService  environment.Service
}

func NewService(deployments core.DeploymentRepository, envService environment.Service) Service {
	return &service{deployments, envService}
}

func (s *service) GetByApp(appId uuid.UUID) (*core.AppStatus, error) {
	deployments, err := s.deployments.FindByApp(appId)
	if err != nil {
		return nil, errors.Wrap(err, "Error retrieving deployment status")
	}

	appStatus := &core.AppStatus{
		Deployments:       deployments,
		EnvironmentStatus: []core.EnvironmentStatus{},
	}

	environmentMap := map[string]core.EnvironmentStatus{}

	for _, deploymentStatus := range deployments {
		if _, ok := environmentMap[deploymentStatus.EnvironmentName]; !ok {
			environmentStatus, err := s.envService.GetStatus(deploymentStatus.EnvironmentName)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("Error retrieving status for environment %q", deploymentStatus.EnvironmentName))
			}
			environmentMap[deploymentStatus.EnvironmentName] = *environmentStatus
			appStatus.EnvironmentStatus = append(appStatus.EnvironmentStatus, *environmentStatus)
		}
	}

	return appStatus, nil
}
