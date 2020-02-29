package deploymentstatus

import (
	"testing"

	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/pkg/stage"

	"github.com/pkg/errors"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/stretchr/testify/assert"
)

func Test_GetByApp(t *testing.T) {
	status := core.Deployment{
		DeploymentReservation: core.DeploymentReservation{
			Name:  "myDeployment",
			AppId: uuid.New(),
		},
		DeploymentRecord: core.DeploymentRecord{
			StageName: "mystage",
			Doc: core.DeploymentDoc{
				Status: &core.DeploymentStatus{
					Revisions: []core.DeploymentRevisionStatus{
						core.DeploymentRevisionStatus{
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

	stageService := &stage.FakeService{
		GetStatusFn: func(stageName string) (*core.StageStatus, error) {
			assert.Equal(t, "mystage", stageName)
			return &core.StageStatus{
				StageName: "mystage",
				Healthy:   true,
			}, nil
		},
	}

	service := service{deployments: deploymentRepository, stageService: stageService}

	result, err := service.GetByApp(status.AppId)

	assert.NoError(t, err)
	assert.Len(t, result.StageStatuses, 1)
	assert.Equal(t, "mystage", result.StageStatuses[0].StageName)
	assert.Len(t, result.Deployments, 1)
	assert.Equal(t, status, result.Deployments[0])
	assert.Equal(t, 1, stageService.GetStatusCallCount)
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

func Test_GetByApp_StageStatusError_ReturnsError(t *testing.T) {
	deployment := core.Deployment{
		DeploymentRecord: core.DeploymentRecord{
			StageName: "mystage",
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		FindByAppFn: func(appId uuid.UUID) ([]core.Deployment, error) {
			return []core.Deployment{deployment}, nil
		},
	}

	stageService := &stage.FakeService{
		GetStatusFn: func(string) (*core.StageStatus, error) {
			return nil, errors.New("test")
		},
	}

	service := service{deployments: deploymentRepository, stageService: stageService}

	result, err := service.GetByApp(uuid.New())

	assert.Nil(t, result)
	assert.Equal(t, "Error retrieving stage status for stage \"mystage\": test", err.Error())

}

func Test_UpdateStatus(t *testing.T) {
	status := &core.DeploymentStatus{}

	deploymentRepository := &core.FakeDeploymentRepository{
		UpdateStatusFn: func(deploymentNameArg string, stageNameArg string, statusArg *core.DeploymentStatus) error {
			assert.Same(t, status, statusArg)
			assert.Equal(t, "mydeployment", deploymentNameArg)
			assert.Equal(t, "mystage", stageNameArg)
			return nil
		},
	}

	service := service{deployments: deploymentRepository}
	err := service.UpdateStatus("mydeployment", "mystage", status)

	assert.NoError(t, err)
	assert.Equal(t, 1, deploymentRepository.UpdateStatusCallCount)
}

func Test_UpdateStatus_WhenUpdateStatusError_ReturnsError(t *testing.T) {
	status := &core.DeploymentStatus{}

	deploymentRepository := &core.FakeDeploymentRepository{
		UpdateStatusFn: func(string, string, *core.DeploymentStatus) error {
			return errors.New("test")
		},
	}

	service := service{deployments: deploymentRepository}

	err := service.UpdateStatus("mydeployment", "mystage", status)

	assert.Equal(t, "Error saving status for deployment \"mydeployment\" in stage \"mystage\": test", err.Error())
}
