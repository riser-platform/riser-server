package rollout

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

// See snapshot test for happy path

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
	assert.Equal(t, `a deployment with the name "myapp.myns" does not exist in environment "dev"`, vErr.Error())
}

func Test_UpdateTraffic_ValidatesRevisionStatus(t *testing.T) {
	deployments := &core.FakeDeploymentRepository{
		GetByNameFn: func(name *core.NamespacedName, envName string) (*core.Deployment, error) {
			return &core.Deployment{
				DeploymentRecord: core.DeploymentRecord{
					Doc: core.DeploymentDoc{
						Status: &core.DeploymentStatus{
							Revisions: []core.DeploymentRevisionStatus{
								{
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

// Race condition for a brand new deployment in a new environment
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
