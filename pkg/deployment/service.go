package deployment

import (
	"database/sql"
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
	Update(deployment *core.DeploymentConfig, committer state.Committer) error
}

type service struct {
	secrets     secret.Service
	stages      core.StageRepository
	deployments core.DeploymentRepository
}

func NewService(secrets secret.Service, stages core.StageRepository, deployments core.DeploymentRepository) Service {
	return &service{secrets, stages, deployments}
}

func (s *service) Update(deploymentConfig *core.DeploymentConfig, committer state.Committer) error {
	riserGeneration, err := s.prepareForDeployment(deploymentConfig)
	if err != nil {
		return err
	}
	stage, err := s.stages.Get(deploymentConfig.Stage)
	if err != nil {
		return err
	}

	secretNames, err := s.secrets.FindNamesByStage(deploymentConfig.App.Name, deploymentConfig.Stage)
	if err != nil {
		return err
	}
	ctx := &core.DeploymentContext{
		Deployment:      deploymentConfig,
		Stage:           &stage.Doc.Config,
		RiserGeneration: riserGeneration,
		SecretNames:     secretNames,
	}
	err = deploy(ctx, committer)
	if err != nil {
		// TODO: Log rollback error but don't return since we want the deployment error to flow to caller
		_, _ = s.deployments.RollbackGeneration(deploymentConfig.Name, deploymentConfig.Stage, riserGeneration)
		return err
	}

	return nil
}

func (s *service) prepareForDeployment(deploymentConfig *core.DeploymentConfig) (riserGeneration int64, err error) {
	applyDefaults(deploymentConfig)
	existingDeployment, err := s.deployments.Get(deploymentConfig.Name, deploymentConfig.Stage)
	if err != nil && err != sql.ErrNoRows {
		return 0, errors.Wrap(err, fmt.Sprintf("Error retrieving deployment %q in stage %q", deploymentConfig.Name, deploymentConfig.Stage))
	}
	if err == sql.ErrNoRows {
		// TODO: Ensure that the deployment name does not exist in another stage by another app (edge case)
		err = s.deployments.Create(&core.Deployment{
			Name:            deploymentConfig.Name,
			StageName:       deploymentConfig.Stage,
			AppName:         deploymentConfig.App.Name,
			RiserGeneration: 1,
		})
		riserGeneration = 1
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("Error creating deployment %q in stage %q", deploymentConfig.Name, deploymentConfig.Stage))
		}
	} else if existingDeployment.AppName != deploymentConfig.App.Name {
		return 0, &core.ValidationError{Message: fmt.Sprintf("A deployment with the name %q is owned by app %q", deploymentConfig.Name, existingDeployment.AppName)}
	} else {
		riserGeneration, err = s.deployments.IncrementGeneration(deploymentConfig.Name, deploymentConfig.Stage)
		if err != nil {
			return 0, errors.Wrap(err, "Error incrementing deployment generation")
		}
	}

	return riserGeneration, nil
}

func deploy(ctx *core.DeploymentContext, committer state.Committer) error {
	// This is a one-off validation until we rationalize our validation strategy.
	// TODO: Once rules are factored out of api/v1/model use RulesNamingIdentifier (creates a circular dep)
	err := validation.Validate(ctx.Deployment.Name,
		validation.Required,
		validation.RuneLength(3, 63),
		validation.Match(regexp.MustCompile("^[a-z][a-z0-9-]+$")).Error("must be lowercase, alphanumeric, and start with a letter"))
	if err != nil {
		// It's important that we print the full deployment name here as the end user can use short hand and just provide the suffix, which can cause
		// confusion (e.g. the suffix may be short enough but not <appName>-<deploymentSuffix>)
		return core.NewValidationError(fmt.Sprintf("invalid deployment name %q", ctx.Deployment.Name), err)
	}

	// TODO: Flag on ctx.Stage.KNativeEnabled
	// deploymentResource, err := resources.CreateDeployment(ctx)
	// if err != nil {
	// 	return err
	// }

	// serviceResource, err := resources.CreateService(ctx)
	// if err != nil {
	// 	return err
	// }

	// virtualServiceResource, err := resources.CreateVirtualService(ctx)
	// if err != nil {
	// 	return err
	// }

	knativeServiceResource := resources.CreateKNativeService(ctx)

	// resourceFiles, err := state.RenderDeployment(ctx.Deployment, deploymentResource, serviceResource, virtualServiceResource)
	resourceFiles, err := state.RenderDeployment(ctx.Deployment, knativeServiceResource)
	if err != nil {
		return err
	}

	return committer.Commit(fmt.Sprintf("Updating resources for %q in stage %q", ctx.Deployment.App.Name, ctx.Deployment.Stage), resourceFiles)
}
