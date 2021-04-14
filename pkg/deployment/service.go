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

	"github.com/pkg/errors"

	"github.com/riser-platform/riser-server/pkg/state"
)

type Service interface {
	Update(deployment *core.DeploymentConfig, committer state.Committer, dryRun bool) (riserRevision int64, err error)
	Delete(name *core.NamespacedName, envName string, committer state.Committer) error
}

type service struct {
	namespaceService   namespace.Service
	secrets            core.SecretMetaRepository
	environments       core.EnvironmentRepository
	deployments        core.DeploymentRepository
	reservationService deploymentreservation.Service
}

func NewService(
	apps core.AppRepository,
	namespaceService namespace.Service,
	secrets core.SecretMetaRepository,
	environments core.EnvironmentRepository,
	deployments core.DeploymentRepository,
	reservationService deploymentreservation.Service) Service {
	return &service{namespaceService, secrets, environments, deployments, reservationService}
}

func (s *service) Delete(name *core.NamespacedName, envName string, committer state.Committer) error {
	// Deleting the deployment is safe to do before we perform the commit since it's a soft delete and therefore idempotent
	err := s.deployments.Delete(name, envName)
	if err != nil {
		if err == core.ErrNotFound {
			return core.NewValidationErrorMessage(fmt.Sprintf("There is no deployment by the name %q in environment %q", name, envName))
		}
		return errors.Wrap(err, "error deleting deployment")
	}

	files := state.RenderDeleteDeployment(name.Name, name.Namespace)
	return committer.Commit(fmt.Sprintf("Deleting deployment %q", name), files)
}

func (s *service) Update(deploymentConfig *core.DeploymentConfig, committer state.Committer, dryRun bool) (riserRevision int64, err error) {
	riserRevision, err = s.prepareForDeployment(deploymentConfig, dryRun)
	if err != nil {
		return 0, err
	}
	environment, err := s.environments.Get(deploymentConfig.EnvironmentName)
	if err != nil {
		return 0, err
	}

	secrets, err := s.secrets.ListByAppInEnvironment(core.NewNamespacedName(deploymentConfig.Name, deploymentConfig.Namespace), deploymentConfig.EnvironmentName)
	if err != nil {
		return 0, err
	}
	ctx := &core.DeploymentContext{
		DeploymentConfig:  deploymentConfig,
		EnvironmentConfig: &environment.Doc.Config,
		RiserRevision:     riserRevision,
		Secrets:           secrets,
	}
	err = deploy(ctx, committer)
	if err != nil {
		// TODO: Log rollback error but don't return since we want the original deployment error to flow to caller
		_, _ = s.deployments.RollbackRevision(
			core.NewNamespacedName(deploymentConfig.Name, deploymentConfig.Namespace), deploymentConfig.EnvironmentName, riserRevision)
		return 0, err
	}

	return riserRevision, nil
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

	existingDeployment, err := s.deployments.GetByReservation(reservation.Id, deploymentConfig.EnvironmentName)
	if err != nil && err != core.ErrNotFound {
		return 0, errors.Wrap(err, fmt.Sprintf("Error retrieving deployment %q in environment %q", deploymentConfig.Name, deploymentConfig.EnvironmentName))
	}
	if err == core.ErrNotFound {
		riserRevision = 1
		deploymentConfig.Traffic = computeTraffic(riserRevision, deploymentConfig, nil)
		err = s.deployments.Create(&core.DeploymentRecord{
			Id:              uuid.New(),
			ReservationId:   reservation.Id,
			EnvironmentName: deploymentConfig.EnvironmentName,
			RiserRevision:   riserRevision,
			Doc: core.DeploymentDoc{
				Traffic: deploymentConfig.Traffic,
			},
		})
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("Error creating deployment %q in environment %q", deploymentConfig.Name, deploymentConfig.EnvironmentName))
		}
	} else if existingDeployment.AppId != deploymentConfig.App.Id {
		return 0, &core.ValidationError{Message: fmt.Sprintf("A deployment with the name %q is owned by app %q", deploymentConfig.Name, existingDeployment.AppId)}
	} else {
		if !dryRun {
			riserRevision, err = s.deployments.IncrementRevision(
				core.NewNamespacedName(deploymentConfig.Name, deploymentConfig.Namespace), deploymentConfig.EnvironmentName)
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
				deploymentConfig.EnvironmentName,
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
	resourceFiles, err := state.RenderDeployment(ctx.DeploymentConfig, createDeployResources(ctx)...)
	if err != nil {
		return err
	}

	// Create the namespace resource whether we need to or not to ensure that it exists and that it's up-to-date
	clusterResourceFiles, err := state.RenderGeneric(ctx.DeploymentConfig.EnvironmentName,
		resources.CreateNamespace(ctx.DeploymentConfig.Namespace, ctx.DeploymentConfig.EnvironmentName))
	if err != nil {
		return nil
	}

	resourceFiles = append(resourceFiles, clusterResourceFiles...)

	return committer.Commit(fmt.Sprintf("Updating resources for \"%s.%s\" in environment %q", ctx.DeploymentConfig.Name, ctx.DeploymentConfig.Namespace, ctx.DeploymentConfig.EnvironmentName), resourceFiles)
}

func createDeployResources(ctx *core.DeploymentContext) []state.KubeResource {
	return []state.KubeResource{
		resources.CreateHealthcheckDenyPolicy(ctx),
		resources.CreateKNativeConfiguration(ctx),
		resources.CreateKNativeRoute(ctx),
	}
}
