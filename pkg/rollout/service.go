package rollout

import (
	"fmt"

	"github.com/riser-platform/riser-server/pkg/core"
)

type Service interface {
	UpdateTraffic(deploymentName, stageName string, rollout core.TrafficConfig) error
}

type service struct {
	deployments core.DeploymentRepository
}

func NewService(deployments core.DeploymentRepository) Service {
	return &service{deployments}
}

func (s *service) UpdateTraffic(deploymentName, stageName string, traffic core.TrafficConfig) error {
	return fmt.Errorf("Not implemented: Receieved %s-%s\n%#v", deploymentName, stageName, traffic)
}
