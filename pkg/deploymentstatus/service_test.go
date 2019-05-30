package deploymentstatus

import (
	"testing"

	"github.com/riser-platform/riser-server/pkg/stage"

	"github.com/pkg/errors"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/stretchr/testify/assert"
)

func Test_GetSummary(t *testing.T) {
	status := core.DeploymentStatus{
		AppName:        "myapp",
		DeploymentName: "myDeployment",
		StageName:      "mystage",
		Doc: &core.DeploymentStatusDoc{
			RolloutStatus:       "InProgress",
			RolloutStatusReason: "Deploying",
			RolloutRevision:     1,
			DockerImage:         "foo:v1.0",
			Problems: []core.DeploymentStatusProblem{
				core.DeploymentStatusProblem{Count: 1, Message: "testProblem"},
			},
		},
	}

	statusRepository := &core.FakeDeploymentStatusRepository{
		FindByAppFn: func(appName string) ([]core.DeploymentStatus, error) {
			assert.Equal(t, "myapp", appName)
			return []core.DeploymentStatus{status}, nil
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

	service := service{statuses: statusRepository, stageService: stageService}

	result, err := service.GetSummary("myapp")

	assert.NoError(t, err)
	assert.Len(t, result.StageStatuses, 1)
	assert.Equal(t, "mystage", result.StageStatuses[0].StageName)
	assert.Len(t, result.DeploymentStatuses, 1)
	assert.Equal(t, status, result.DeploymentStatuses[0])
	assert.Equal(t, 1, stageService.GetStatusCallCount)
}

func Test_GetSummary_StatusRepoErr_ReturnsErr(t *testing.T) {
	statusRepository := &core.FakeDeploymentStatusRepository{
		FindByAppFn: func(string) ([]core.DeploymentStatus, error) {
			return nil, errors.New("test")
		},
	}

	service := service{statuses: statusRepository}

	result, err := service.GetSummary("myapp")

	assert.Nil(t, result)
	assert.Equal(t, err.Error(), "Error retrieving deployment status: test")
}

func Test_GetSummary_StageStatusError_ReturnsError(t *testing.T) {
	status := core.DeploymentStatus{
		AppName:   "myapp",
		StageName: "mystage",
	}

	statusRepository := &core.FakeDeploymentStatusRepository{
		FindByAppFn: func(appName string) ([]core.DeploymentStatus, error) {
			assert.Equal(t, "myapp", appName)
			return []core.DeploymentStatus{status}, nil
		},
	}

	stageService := &stage.FakeService{
		GetStatusFn: func(string) (*core.StageStatus, error) {
			return nil, errors.New("test")
		},
	}

	service := service{statuses: statusRepository, stageService: stageService}

	result, err := service.GetSummary("myapp")

	assert.Nil(t, result)
	assert.Equal(t, "Error retrieving stage status for stage \"mystage\": test", err.Error())

}

func Test_Save(t *testing.T) {
	status := &core.DeploymentStatus{StageName: "mystage"}

	statusRepository := &core.FakeDeploymentStatusRepository{
		SaveFn: func(statusArg *core.DeploymentStatus) error {
			assert.Equal(t, status, statusArg)
			return nil
		},
	}

	stageService := &stage.FakeService{
		PingFn: func(stageName string) error {
			assert.Equal(t, "mystage", stageName)
			return nil
		},
	}

	service := service{statuses: statusRepository, stageService: stageService}

	err := service.Save(status)

	assert.NoError(t, err)
	assert.Equal(t, 1, statusRepository.SaveCallCount)
	assert.Equal(t, 1, stageService.PingCallCount)
}

func Test_Save_WhenPingStageError_ReturnsError(t *testing.T) {
	status := &core.DeploymentStatus{StageName: "mystage"}

	stageService := &stage.FakeService{
		PingFn: func(string) error {
			return errors.New("test")
		},
	}

	service := service{stageService: stageService}

	err := service.Save(status)

	assert.Equal(t, "Error saving stage \"mystage\": test", err.Error())
}

func Test_Save_WhenSaveError_ReturnsError(t *testing.T) {
	status := &core.DeploymentStatus{
		AppName: "myapp",
	}

	statusRepository := &core.FakeDeploymentStatusRepository{
		SaveFn: func(statusArg *core.DeploymentStatus) error {
			return errors.New("test")
		},
	}

	stageService := &stage.FakeService{
		PingFn: func(string) error {
			return nil
		},
	}

	service := service{statuses: statusRepository, stageService: stageService}

	err := service.Save(status)

	assert.Equal(t, "Error saving deployment status for app \"myapp\": test", err.Error())
}
