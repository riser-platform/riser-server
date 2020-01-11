package rollout

import (
	"fmt"
	"github.com/pkg/errors"

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
	deployment, err := s.deployments.Get(deploymentName, stageName)
	if err != nil {
		if err == core.ErrNotFound {
			return &core.ValidationError{Message: fmt.Sprintf("a deployment with the name %q does not exist in stage %q", deploymentName, stageName)}
		}
		return errors.Wrap(err, "error getting deployment")
	}

	err = validateTrafficRules(traffic, deployment)
	if err != nil {
		return err
	}

	return fmt.Errorf("Not implemented: Received %s-%s\n%#v", deploymentName, stageName, traffic)
}

func validateTrafficRules(traffic core.TrafficConfig, deployment *core.Deployment) error {
	revisions := map[int64]bool{}
	if deployment.Doc.Status != nil {
		for _, rev := range deployment.Doc.Status.Revisions {
			revisions[rev.RiserGeneration] = true
		}
	}
	for _, rule := range traffic {
		if _, ok := revisions[rule.RiserGeneration]; !ok {
			return &core.ValidationError{Message: fmt.Sprintf(`revision "%d" either does not exist or has not reported its status yet`, rule.RiserGeneration)}
		}
	}

	return nil
}
