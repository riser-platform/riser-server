package deploymentstatus

import (
	"testing"

	"github.com/google/uuid"

	"github.com/pkg/errors"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/environment"

	"github.com/stretchr/testify/assert"
)

func Test_GetByApp(t *testing.T) {
	status := core.Deployment{
		DeploymentReservation: core.DeploymentReservation{
			Name:  "myDeployment",
			AppId: uuid.New(),
		},
		DeploymentRecord: core.DeploymentRecord{
			EnvironmentName: "myenv",
			Doc: core.DeploymentDoc{
				Status: &core.DeploymentStatus{
					Revisions: []core.DeploymentRevisionStatus{
						{
							Name:                 "rev1",
							RiserRevision:        1,
							RevisionStatus:       "InProgress",
							RevisionStatusReason: "Deploying",
							DockerImage:          "foo:v1.0",
						},
					},
				},
			},
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		FindByAppFn: func(appId uuid.UUID) ([]core.Deployment, error) {
			assert.Equal(t, status.AppId, appId)
			return []core.Deployment{status}, nil
		},
	}

	environmentService := &environment.FakeService{
		GetStatusFn: func(envName string) (*core.EnvironmentStatus, error) {
			assert.Equal(t, "myenv", envName)
			return &core.EnvironmentStatus{
				EnvironmentName: "myenv",
				Healthy:         true,
			}, nil
		},
	}

	service := service{deployments: deploymentRepository, envService: environmentService}

	result, err := service.GetByApp(status.AppId)

	assert.NoError(t, err)
	assert.Equal(t, status.AppId, result.AppId)
	assert.Len(t, result.EnvironmentStatus, 1)
	assert.Equal(t, "myenv", result.EnvironmentStatus[0].EnvironmentName)
	assert.Len(t, result.Deployments, 1)
	assert.Equal(t, status, result.Deployments[0])
	assert.Equal(t, 1, environmentService.GetStatusCallCount)
}

func Test_GetByApp_StatusRepoErr_ReturnsErr(t *testing.T) {
	deploymentRepository := &core.FakeDeploymentRepository{
		FindByAppFn: func(uuid.UUID) ([]core.Deployment, error) {
			return nil, errors.New("test")
		},
	}

	service := service{deployments: deploymentRepository}

	result, err := service.GetByApp(uuid.New())

	assert.Nil(t, result)
	assert.Equal(t, err.Error(), "Error retrieving deployment status: test")
}

func Test_GetByApp_EnvironmentStatusError_ReturnsError(t *testing.T) {
	deployment := core.Deployment{
		DeploymentRecord: core.DeploymentRecord{
			EnvironmentName: "myenv",
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		FindByAppFn: func(appId uuid.UUID) ([]core.Deployment, error) {
			return []core.Deployment{deployment}, nil
		},
	}

	environmentService := &environment.FakeService{
		GetStatusFn: func(string) (*core.EnvironmentStatus, error) {
			return nil, errors.New("test")
		},
	}

	service := service{deployments: deploymentRepository, envService: environmentService}

	result, err := service.GetByApp(uuid.New())

	assert.Nil(t, result)
	assert.Equal(t, "Error retrieving status for environment \"myenv\": test", err.Error())
}
