package rollout

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/pkg/state"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_UpdateTraffic(t *testing.T) {
	appId := uuid.New()
	name := core.NewNamespacedName("myapp", "myns")
	committer := state.NewDryRunCommitter()

	deployments := &core.FakeDeploymentRepository{
		GetByNameFn: func(nameArg *core.NamespacedName, stageName string) (*core.Deployment, error) {
			assert.Equal(t, name, nameArg)
			return &core.Deployment{
				DeploymentReservation: core.DeploymentReservation{
					Name:  "myapp",
					AppId: appId,
				},
				DeploymentRecord: core.DeploymentRecord{
					Doc: core.DeploymentDoc{
						Status: &core.DeploymentStatus{
							Revisions: []core.DeploymentRevisionStatus{
								core.DeploymentRevisionStatus{
									RiserRevision: 1,
								},
							},
						},
					},
				},
			}, nil
		},
	}

	apps := &core.FakeAppRepository{
		GetFn: func(id uuid.UUID) (*core.App, error) {
			assert.Equal(t, appId, id)
			return &core.App{
				Id:   id,
				Name: "myapp",
			}, nil
		},
	}

	traffic := core.TrafficConfig{
		core.TrafficConfigRule{
			RiserRevision: 1,
			Percent:       100,
		},
	}

	svc := service{apps, deployments}

	err := svc.UpdateTraffic(name, "dev", traffic, committer)

	assert.NoError(t, err)
	assert.Len(t, committer.Commits, 1)
	assert.Equal(t, `Updating resources for "myapp.myns" in stage "dev"`, committer.Commits[0].Message)
	assert.Len(t, committer.Commits[0].Files, 1)
	// TODO: Factor out snapshot code from pkg/deployment and add snapshot test for route file contents
}

func Test_UpdateTraffic_ReturnsGetDeploymentError(t *testing.T) {
	deployments := &core.FakeDeploymentRepository{
		GetByNameFn: func(*core.NamespacedName, string) (*core.Deployment, error) {
			return nil, errors.New("test")
		},
	}

	svc := service{deployments: deployments}

	result := svc.UpdateTraffic(core.NewNamespacedName("myapp", "myns"), "dev", core.TrafficConfig{}, nil)

	assert.Equal(t, "error getting deployment: test", result.Error())
}

func Test_UpdateTraffic_WhenDeploymentDoesNotExist(t *testing.T) {
	deployments := &core.FakeDeploymentRepository{
		GetByNameFn: func(*core.NamespacedName, string) (*core.Deployment, error) {
			return nil, core.ErrNotFound
		},
	}

	svc := service{deployments: deployments}

	result := svc.UpdateTraffic(core.NewNamespacedName("myapp", "myns"), "dev", core.TrafficConfig{}, nil)

	assert.IsType(t, &core.ValidationError{}, result)
	vErr := result.(*core.ValidationError)
	assert.Equal(t, `a deployment with the name "myapp.myns" does not exist in stage "dev"`, vErr.Error())
}

func Test_UpdateTraffic_ValidatesRevisionStatus(t *testing.T) {
	deployments := &core.FakeDeploymentRepository{
		GetByNameFn: func(name *core.NamespacedName, stageName string) (*core.Deployment, error) {
			return &core.Deployment{
				DeploymentRecord: core.DeploymentRecord{
					Doc: core.DeploymentDoc{
						Status: &core.DeploymentStatus{
							Revisions: []core.DeploymentRevisionStatus{
								core.DeploymentRevisionStatus{
									RiserRevision: 1,
								},
							},
						},
					},
				},
			}, nil
		},
	}

	apps := &core.FakeAppRepository{
		GetFn: func(id uuid.UUID) (*core.App, error) {
			return &core.App{
				Id:   id,
				Name: "myapp",
			}, nil
		},
	}

	traffic := core.TrafficConfig{
		core.TrafficConfigRule{
			RiserRevision: 1,
			Percent:       50,
		},
		core.TrafficConfigRule{
			RiserRevision: 2,
			Percent:       50,
		},
	}

	svc := service{apps, deployments}

	result := svc.UpdateTraffic(core.NewNamespacedName("myapp", "myns"), "dev", traffic, nil)

	assert.Equal(t, `revision "2" either does not exist or has not reported its status yet`, result.Error())
}

// Race condition for a brand new deployment in a new stage
func Test_UpdateTraffic_ValidatesRevisionStatus_NoStatus(t *testing.T) {
	deployments := &core.FakeDeploymentRepository{
		GetByNameFn: func(*core.NamespacedName, string) (*core.Deployment, error) {
			return &core.Deployment{}, nil
		},
	}

	apps := &core.FakeAppRepository{
		GetFn: func(id uuid.UUID) (*core.App, error) {
			return &core.App{
				Id:   id,
				Name: "myapp",
			}, nil
		},
	}

	traffic := core.TrafficConfig{
		core.TrafficConfigRule{
			RiserRevision: 1,
			Percent:       100,
		},
	}

	svc := service{apps, deployments}

	result := svc.UpdateTraffic(core.NewNamespacedName("myapp", "myns"), "dev", traffic, nil)

	assert.Equal(t, `revision "1" either does not exist or has not reported its status yet`, result.Error())
}
