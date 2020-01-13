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
	if err := validateDeploymentConfig(deploymentConfig); err != nil {
		return 0, err
	}
	existingDeployment, err := s.deployments.Get(deploymentConfig.Name, deploymentConfig.Stage)
	if err != nil && err != sql.ErrNoRows {
		return 0, errors.Wrap(err, fmt.Sprintf("Error retrieving deployment %q in stage %q", deploymentConfig.Name, deploymentConfig.Stage))
	}
	if err == sql.ErrNoRows {
		riserGeneration = 1
		deploymentConfig.Traffic = computeTraffic(riserGeneration, deploymentConfig, nil)
		// TODO: Ensure that the deployment name does not exist in another stage by another app (edge case)
		err = s.deployments.Create(&core.Deployment{
			Name:            deploymentConfig.Name,
			StageName:       deploymentConfig.Stage,
			AppName:         deploymentConfig.App.Name,
			RiserGeneration: riserGeneration,
			Doc: core.DeploymentDoc{
				Traffic: deploymentConfig.Traffic,
			},
		})
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
		deploymentConfig.Traffic = computeTraffic(riserGeneration, deploymentConfig, existingDeployment)

		err = s.deployments.UpdateTraffic(deploymentConfig.Name, deploymentConfig.Stage, riserGeneration, deploymentConfig.Traffic)
		if err != nil {
			return 0, errors.Wrap(err, "Error updating traffic")
		}
	}

	return riserGeneration, nil
}

func computeTraffic(riserGeneration int64, deploymentConfig *core.DeploymentConfig, existingDeployment *core.Deployment) core.TrafficConfig {
	newRule := core.TrafficConfigRule{
		RiserGeneration: riserGeneration,
		RevisionName:    fmt.Sprintf("%s-%d", deploymentConfig.Name, riserGeneration),
	}

	if deploymentConfig.ManualRollout && existingDeployment != nil {
		newRule.Percent = 0
		trafficConfig := core.TrafficConfig{newRule}
		for _, rule := range existingDeployment.Doc.Traffic {
			if rule.Percent > 0 {
				trafficConfig = append(trafficConfig, rule)
			}
		}
		return trafficConfig
	}

	newRule.Percent = 100
	return core.TrafficConfig{newRule}
}

// This is a one-off validation until we rationalize our validation strategy (API layer or service layer).
func validateDeploymentConfig(deployment *core.DeploymentConfig) error {
	// TODO: Once rules are factored out of api/v1/model use RulesNamingIdentifier (creates a circular dep)
	err := validation.Validate(deployment.Name,
		validation.Required,
		validation.RuneLength(3, 63),
		validation.
			Match(regexp.MustCompile(fmt.Sprintf("^%s(-.+)?$", deployment.App.Name))).
			Error(fmt.Sprintf("must be either %q or start with \"%s-\"", deployment.App.Name, deployment.App.Name)),
		validation.Match(regexp.MustCompile("^[a-z][a-z0-9-]+$")).Error("must be lowercase, alphanumeric, and start with a letter"))
	if err != nil {
		return core.NewValidationError(fmt.Sprintf("invalid deployment name %q", deployment.Name), err)
	}
	return nil
}

func deploy(ctx *core.DeploymentContext, committer state.Committer) error {
	resourceFiles, err := state.RenderDeployment(ctx.Deployment,
		resources.CreateKNativeConfiguration(ctx),
		resources.CreateKNativeRoute(ctx))
	if err != nil {
		return err
	}

	return committer.Commit(fmt.Sprintf("Updating resources for %q in stage %q", ctx.Deployment.Name, ctx.Deployment.Stage), resourceFiles)
}
