package rollout

import (
	"errors"
	"testing"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_UpdateTraffic_ReturnsGetDeploymentError(t *testing.T) {
	deployments := &core.FakeDeploymentRepository{
		GetFn: func(name string, stageName string) (*core.Deployment, error) {
			return nil, errors.New("test")
		},
	}

	svc := service{deployments: deployments}

	result := svc.UpdateTraffic("myapp", "dev", core.TrafficConfig{})

	assert.Equal(t, "error getting deployment: test", result.Error())
}

func Test_UpdateTraffic_WhenDeploymentDoesNotExist(t *testing.T) {
	deployments := &core.FakeDeploymentRepository{
		GetFn: func(name string, stageName string) (*core.Deployment, error) {
			return nil, core.ErrNotFound
		},
	}

	svc := service{deployments: deployments}

	result := svc.UpdateTraffic("myapp", "dev", core.TrafficConfig{})

	assert.IsType(t, &core.ValidationError{}, result)
	vErr := result.(*core.ValidationError)
	assert.Equal(t, `a deployment with the name "myapp" does not exist in stage "dev"`, vErr.Error())
}

func Test_UpdateTraffic_ValidatesRevisionStatus(t *testing.T) {
	deployments := &core.FakeDeploymentRepository{
		GetFn: func(name string, stageName string) (*core.Deployment, error) {
			assert.Equal(t, "myapp", name)
			assert.Equal(t, "dev", stageName)
			return &core.Deployment{
				Doc: core.DeploymentDoc{
					Status: &core.DeploymentStatus{
						Revisions: []core.DeploymentRevisionStatus{
							core.DeploymentRevisionStatus{
								RiserGeneration: 1,
							},
						},
					},
				},
			}, nil
		},
	}

	traffic := core.TrafficConfig{
		core.TrafficConfigRule{
			RiserGeneration: 1,
			Percent:         50,
		},
		core.TrafficConfigRule{
			RiserGeneration: 2,
			Percent:         50,
		},
	}

	svc := service{deployments: deployments}

	result := svc.UpdateTraffic("myapp", "dev", traffic)

	assert.Equal(t, `revision "2" either does not exist or has not reported its status yet`, result.Error())
}

// Race condition for a brand new deployment in a new stage
func Test_UpdateTraffic_ValidatesRevisionStatus_NoStatus(t *testing.T) {
	deployments := &core.FakeDeploymentRepository{
		GetFn: func(name string, stageName string) (*core.Deployment, error) {
			assert.Equal(t, "myapp", name)
			assert.Equal(t, "dev", stageName)
			return &core.Deployment{}, nil
		},
	}

	traffic := core.TrafficConfig{
		core.TrafficConfigRule{
			RiserGeneration: 1,
			Percent:         100,
		},
	}

	svc := service{deployments: deployments}

	result := svc.UpdateTraffic("myapp", "dev", traffic)

	assert.Equal(t, `revision "1" either does not exist or has not reported its status yet`, result.Error())
}
