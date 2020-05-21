package rollout

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/state"
	"github.com/riser-platform/riser-server/pkg/state/resources"
)

type Service interface {
	UpdateTraffic(name *core.NamespacedName, envName string, rollout core.TrafficConfig, committer state.Committer) error
}

type service struct {
	apps        core.AppRepository
	deployments core.DeploymentRepository
}

func NewService(apps core.AppRepository, deployments core.DeploymentRepository) Service {
	return &service{apps, deployments}
}

func (s *service) UpdateTraffic(name *core.NamespacedName, envName string, traffic core.TrafficConfig, committer state.Committer) error {
	deployment, err := s.deployments.GetByName(name, envName)
	if err != nil {
		if err == core.ErrNotFound {
			return &core.ValidationError{Message: fmt.Sprintf("a deployment with the name %q does not exist in environment %q", name, envName)}
		}
		return errors.Wrap(err, "error getting deployment")
	}

	app, err := s.apps.Get(deployment.AppId)
	if err != nil {
		return errors.Wrap(err, "error getting app")
	}

	err = validateTrafficRules(traffic, deployment)
	if err != nil {
		return err
	}

	// TODO: Refactor underlying code to not require the entire deployment context. Currently this is hydrated only with fields that we know are needed
	ctx := &core.DeploymentContext{
		DeploymentConfig: &core.DeploymentConfig{
			Name:            name.Name,
			Namespace:       name.Namespace,
			EnvironmentName: envName,
			Traffic:         traffic,
			App: &model.AppConfig{
				Name: model.AppName(app.Name),
			},
		},
		RiserRevision: deployment.RiserRevision,
	}

	resourceFiles, err := state.RenderRoute(name.Name, name.Namespace, envName, resources.CreateKNativeRoute(ctx))
	if err != nil {
		return err
	}

	return committer.Commit(fmt.Sprintf("Updating resources for %q in environment %q", name, ctx.DeploymentConfig.EnvironmentName), resourceFiles)
}

func validateTrafficRules(traffic core.TrafficConfig, deployment *core.Deployment) error {
	revisions := map[int64]bool{}
	if deployment.Doc.Status != nil {
		for _, rev := range deployment.Doc.Status.Revisions {
			revisions[rev.RiserRevision] = true
		}
	}
	for _, rule := range traffic {
		if _, ok := revisions[rule.RiserRevision]; !ok {
			return &core.ValidationError{Message: fmt.Sprintf(`revision "%d" either does not exist or has not reported its status yet`, rule.RiserRevision)}
		}
	}

	return nil
}
