package deployment

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/state/resources"

	"github.com/riser-platform/riser-server/pkg/secret"

	"github.com/pkg/errors"

	"github.com/riser-platform/riser-server/pkg/state"
)

type Service interface {
	Update(deployment *core.Deployment, committer state.Committer) error
}

type service struct {
	secrets secret.Service
	stages  core.StageRepository
}

func NewService(secrets secret.Service, stages core.StageRepository) Service {
	return &service{secrets, stages}
}

func (s *service) Update(deployment *core.Deployment, committer state.Committer) error {
	secretNames, err := s.secrets.FindNamesByStage(deployment.App.Name, deployment.Stage)
	if err != nil {
		return err
	}

	stage, err := s.stages.Get(deployment.Stage)
	if err != nil {
		return err
	}

	// TODO: Log a warning if the public gateway host is not configured for this stage
	return deploy(deployment, stage.Doc.Config, committer, secretNames)
}

func deploy(deployment *core.Deployment, stageConfig core.StageConfig, committer state.Committer, secretNames []string) error {
	deployment = ApplyDefaults(deployment)
	// This is a one-off validation until we rationalize our validation strategy.
	// TODO: Once rules are factored out of api/v1/model use RulesNamingIdentifier (creates a circular dep)
	err := validation.Validate(deployment.Name,
		validation.Required,
		validation.RuneLength(3, 63),
		validation.Match(regexp.MustCompile("^[a-z][a-z0-9-]+$")).Error("must be lowercase, alphanumeric, and start with a letter"))
	if err != nil {
		// It's important that we print the full deployment name here as the end user can use short hand and just provide the suffix, which can cause
		// confusion (e.g. the suffix may be short enough but not <appName>-<deploymentSuffix>)
		return core.NewValidationError(fmt.Sprintf("invalid deployment name %q", deployment.Name), err)
	}

	var resources []state.KubeResource
	resources, err = createKNativeDeploymentResources(deployment, stageConfig, secretNames)
	// TODO: Don't commit!
	// if stageConfig.KNativeEnabled {
	// 	resources, err = createKNativeDeploymentResources(deployment, stageConfig, secretNames)
	// } else {
	// 	resources, err = createDeploymentResources(deployment, stageConfig, secretNames)
	// }
	if err != nil {
		return err
	}

	resourceFiles, err := state.RenderDeployment(deployment, resources...)
	if err != nil {
		return errors.Wrap(err, "Error rendering deployment resources")
	}

	return committer.Commit(fmt.Sprintf("Updating resources for %q in stage %q", deployment.App.Name, deployment.Stage), resourceFiles)
}

func createKNativeDeploymentResources(deployment *core.Deployment, stageConfig core.StageConfig, secretNames []string) ([]state.KubeResource, error) {
	return []state.KubeResource{resources.CreateKNativeService(deployment, secretNames)}, nil
}

func createDeploymentResources(deployment *core.Deployment, stageConfig core.StageConfig, secretNames []string) ([]state.KubeResource, error) {
	deploymentResource, err := resources.CreateDeployment(deployment, secretNames)
	if err != nil {
		return nil, err
	}

	serviceResource, err := resources.CreateService(deployment)
	if err != nil {
		return nil, err
	}

	virtualServiceResource, err := resources.CreateVirtualService(deployment, stageConfig.PublicGatewayHost)
	if err != nil {
		return nil, err
	}

	return []state.KubeResource{
		deploymentResource,
		serviceResource,
		virtualServiceResource,
	}, nil
}
