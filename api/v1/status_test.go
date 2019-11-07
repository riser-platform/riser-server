package v1

import (
	"testing"
	"time"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/api/v1/model"
)

func Test_mapDeploymentToStatusModel(t *testing.T) {
	deployment := &core.Deployment{
		Name:            "mydeployment",
		StageName:       "mystage",
		RiserGeneration: 4,
		Doc: core.DeploymentDoc{
			Status: &core.DeploymentStatus{
				ObservedRiserGeneration: 3,
				RolloutStatus:           "myrolloutstatus",
				RolloutRevision:         1337,
				RolloutStatusReason:     "myrolloutstatusreason",
				DockerImage:             "mydockerimage",
				Problems: []core.DeploymentStatusProblem{
					core.DeploymentStatusProblem{
						Message: "myproblem1",
						Count:   1,
					},
					core.DeploymentStatusProblem{
						Message: "myproblem2",
						Count:   2,
					},
				},
			},
		},
	}

	result := mapDeploymentToStatusModel(deployment)

	assert.Equal(t, "mydeployment", result.DeploymentName)
	assert.Equal(t, "mystage", result.StageName)
	assert.Equal(t, "myrolloutstatus", result.RolloutStatus)
	assert.Equal(t, int64(3), result.ObservedRiserGeneration)
	assert.Equal(t, int64(4), result.RiserGeneration)
	assert.Equal(t, int64(1337), result.RolloutRevision)
	assert.Equal(t, "myrolloutstatusreason", result.RolloutStatusReason)
	assert.Equal(t, "mydockerimage", result.DockerImage)
	assert.Len(t, result.Problems, 2)
	assert.Equal(t, "myproblem1", result.Problems[0].Message)
	assert.Equal(t, 1, result.Problems[0].Count)
	assert.Equal(t, "myproblem2", result.Problems[1].Message)
	assert.Equal(t, 2, result.Problems[1].Count)
}

func Test_mapDeploymentToStatusModel_NilStatus(t *testing.T) {
	deployment := &core.Deployment{
		Name:      "mydeployment",
		StageName: "mystage",
		Doc:       core.DeploymentDoc{},
	}

	result := mapDeploymentToStatusModel(deployment)

	assert.Equal(t, "mydeployment", result.DeploymentName)
	assert.Equal(t, "mystage", result.StageName)
	assert.Equal(t, "Unknown", result.RolloutStatus)
}

func Test_mapDeploymentStatusFromModel(t *testing.T) {
	deploymentStatus := &model.DeploymentStatusMutable{
		ObservedRiserGeneration: 3,
		RolloutStatus:           "myrolloutstatus",
		RolloutRevision:         1337,
		RolloutStatusReason:     "myrolloutstatusreason",
		DockerImage:             "mydockerimage",
		Problems: []model.DeploymentStatusProblem{
			model.DeploymentStatusProblem{
				Message: "myproblem1",
				Count:   1,
			},
			model.DeploymentStatusProblem{
				Message: "myproblem2",
				Count:   2,
			},
		},
	}

	now := time.Now().Unix()

	result := mapDeploymentStatusFromModel(deploymentStatus)

	assert.InDelta(t, now, result.LastUpdated.Unix(), 3)
	assert.Equal(t, int64(3), result.ObservedRiserGeneration)
	assert.Equal(t, "myrolloutstatus", result.RolloutStatus)
	assert.EqualValues(t, 1337, result.RolloutRevision)
	assert.Equal(t, "myrolloutstatusreason", result.RolloutStatusReason)
	assert.Equal(t, "mydockerimage", result.DockerImage)
	assert.Len(t, result.Problems, 2)
	assert.Equal(t, "myproblem1", result.Problems[0].Message)
	assert.Equal(t, 1, result.Problems[0].Count)
	assert.Equal(t, "myproblem2", result.Problems[1].Message)
	assert.Equal(t, 2, result.Problems[1].Count)
}
