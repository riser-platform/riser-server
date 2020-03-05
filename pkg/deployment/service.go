package deployment

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/pkg/deploymentreservation"
	"github.com/riser-platform/riser-server/pkg/namespace"

	validation "github.com/go-ozzo/ozzo-validation/v3"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/state/resources"

	"github.com/riser-platform/riser-server/pkg/secret"

	"github.com/pkg/errors"

	"github.com/riser-platform/riser-server/pkg/state"
)

type Service interface {
	Update(deployment *core.DeploymentConfig, committer state.Committer, dryRun bool) error
	Delete(name *core.NamespacedName, stageName string, committer state.Committer) error
}

type service struct {
	namespaceService   namespace.Service
	secrets            secret.Service
	stages             core.StageRepository
	deployments        core.DeploymentRepository
	reservationService deploymentreservation.Service
}

func NewService(
	apps core.AppRepository,
	namespaceService namespace.Service,
	secrets secret.Service,
	stages core.StageRepository,
	deployments core.DeploymentRepository,
	reservationService deploymentreservation.Service) Service {
	return &service{namespaceService, secrets, stages, deployments, reservationService}
}

func (s *service) Delete(name *core.NamespacedName, stageName string, committer state.Committer) error {
	// This is safe to do before we perform the commit since it's idempotent
	err := s.deployments.Delete(name, stageName)
	if err != nil {
		if err == core.ErrNotFound {
			return core.NewValidationErrorMessage(fmt.Sprintf("There is no deployment by the name %q in stage %q", name, stageName))
		}
		return errors.Wrap(err, "error deleting deployment")
	}

	files := state.RenderDeleteDeployment(name.Name, name.Namespace, stageName)
	return committer.Commit(fmt.Sprintf("Deleting deployment %q", name), files)
}

func (s *service) Update(deploymentConfig *core.DeploymentConfig, committer state.Committer, dryRun bool) error {
	err := s.namespaceService.EnsureNamespaceInStage(deploymentConfig.Namespace, deploymentConfig.Stage, committer)
	if err != nil {
		return err
	}

	riserRevision, err := s.prepareForDeployment(deploymentConfig, dryRun)
	if err != nil {
		return err
	}
	stage, err := s.stages.Get(deploymentConfig.Stage)
	if err != nil {
		return err
	}

	secrets, err := s.secrets.FindByStage(deploymentConfig.App.Id, deploymentConfig.Stage)
	if err != nil {
		return err
	}
	ctx := &core.DeploymentContext{
		Deployment:    deploymentConfig,
		Stage:         &stage.Doc.Config,
		RiserRevision: riserRevision,
		Secrets:       secrets,
	}
	err = deploy(ctx, committer)
	if err != nil {
		// TODO: Log rollback error but don't return since we want the original deployment error to flow to caller
		_, _ = s.deployments.RollbackRevision(
			core.NewNamespacedName(deploymentConfig.Name, deploymentConfig.Namespace), deploymentConfig.Stage, riserRevision)
		return err
	}

	return nil
}

func (s *service) prepareForDeployment(deploymentConfig *core.DeploymentConfig, dryRun bool) (riserRevision int64, err error) {
	if err := validateDeploymentConfig(deploymentConfig); err != nil {
		return 0, err
	}

	reservation, err := s.reservationService.EnsureReservation(
		deploymentConfig.App.Id,
		core.NewNamespacedName(deploymentConfig.Name, deploymentConfig.Namespace))
	if err != nil {
		return 0, errors.Wrap(err, "Error ensuring deployment reservation")
	}

	existingDeployment, err := s.deployments.GetByReservation(reservation.Id, deploymentConfig.Stage)
	if err != nil && err != core.ErrNotFound {
		return 0, errors.Wrap(err, fmt.Sprintf("Error retrieving deployment %q in stage %q", deploymentConfig.Name, deploymentConfig.Stage))
	}
	if err == core.ErrNotFound {
		riserRevision = 1
		deploymentConfig.Traffic = computeTraffic(riserRevision, deploymentConfig, nil)
		err = s.deployments.Create(&core.DeploymentRecord{
			Id:            uuid.New(),
			ReservationId: reservation.Id,
			StageName:     deploymentConfig.Stage,
			RiserRevision: riserRevision,
			Doc: core.DeploymentDoc{
				Traffic: deploymentConfig.Traffic,
			},
		})
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("Error creating deployment %q in stage %q", deploymentConfig.Name, deploymentConfig.Stage))
		}
	} else if existingDeployment.AppId != deploymentConfig.App.Id {
		return 0, &core.ValidationError{Message: fmt.Sprintf("A deployment with the name %q is owned by app %q", deploymentConfig.Name, existingDeployment.AppId)}
	} else {
		if !dryRun {
			riserRevision, err = s.deployments.IncrementRevision(
				core.NewNamespacedName(deploymentConfig.Name, deploymentConfig.Namespace), deploymentConfig.Stage)
			if err != nil {
				return 0, errors.Wrap(err, "Error incrementing deployment revision")
			}
		}

		// When a deployment was previously deleted, we don't want to compute traffic with the old traffic rules
		if existingDeployment.DeletedAt == nil {
			deploymentConfig.Traffic = computeTraffic(riserRevision, deploymentConfig, &existingDeployment.DeploymentRecord)
		} else {
			deploymentConfig.Traffic = computeTraffic(riserRevision, deploymentConfig, nil)
		}

		if !dryRun {
			err = s.deployments.UpdateTraffic(
				core.NewNamespacedName(deploymentConfig.Name, deploymentConfig.Namespace),
				deploymentConfig.Stage,
				riserRevision,
				deploymentConfig.Traffic)
			if err != nil {
				return 0, errors.Wrap(err, "Error updating traffic")
			}
		}
	}

	return riserRevision, nil
}

func computeTraffic(riserRevision int64, deploymentConfig *core.DeploymentConfig, existingDeployment *core.DeploymentRecord) core.TrafficConfig {
	newRule := core.TrafficConfigRule{
		RiserRevision: riserRevision,
		RevisionName:  fmt.Sprintf("%s-%d", deploymentConfig.Name, riserRevision),
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
