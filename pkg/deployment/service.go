package deployment

import (
	"fmt"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/state/resources"

	"github.com/riser-platform/riser-server/pkg/secret"

	"github.com/pkg/errors"

	"github.com/riser-platform/riser-server/pkg/state"
)

type Service interface {
	Update(deployment *core.NewDeployment, committer state.Committer) error
}

type service struct {
	secrets secret.Service
	stages  core.StageRepository
}

func NewService(secrets secret.Service, stages core.StageRepository) Service {
	return &service{secrets, stages}
}

func (s *service) Update(deployment *core.NewDeployment, committer state.Committer) error {
	secretNames, err := s.secrets.FindNamesByStage(deployment.App.Name, deployment.Stage)
	if err != nil {
		return err
	}

	stage, err := s.stages.Get(deployment.Stage)
	if err != nil {
		return err
	}

	// TODO: Log a warning if the public gateway host is not configured for this stage
	return s.update(deployment, committer, secretNames, stage.Doc.Config.PublicGatewayHost)
}

func (s *service) update(newDeployment *core.NewDeployment, committer state.Committer, secretNames []string, publicGatewayHost string) error {
	newDeployment = ApplyDefaults(newDeployment)

	deployment, err := ApplyOverrides(newDeployment)
	if err != nil {
		return err
	}

	return deploy(deployment, committer, secretNames, publicGatewayHost)
}

func deploy(deployment *core.Deployment, committer state.Committer, secretNames []string, publicGatewayHost string) error {
	deploymentResource, err := resources.CreateDeployment(deployment, secretNames)
	if err != nil {
		return err
	}

	serviceResource, err := resources.CreateService(deployment)
	if err != nil {
		return err
	}

	virtualServiceResource, err := resources.CreateVirtualService(deployment, publicGatewayHost)
	if err != nil {
		return err
	}

	resourceFiles, err := state.RenderDeployment(deployment, deploymentResource, serviceResource, virtualServiceResource)
	if err != nil {
		return errors.Wrap(err, "Error rendering deployment resources")
	}

	return committer.Commit(fmt.Sprintf("Updating resources for %q in stage %q", deployment.App.Name, deployment.Stage), resourceFiles)
}
