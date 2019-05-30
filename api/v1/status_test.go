package v1

import (
	"testing"
	"time"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/api/v1/model"
)

func Test_mapDeploymentStatusToModel(t *testing.T) {
	deploymentStatus := &core.DeploymentStatus{
		AppName:        "myapp",
		DeploymentName: "mydeployment",
		StageName:      "mystage",
		Doc: &core.DeploymentStatusDoc{
			RolloutStatus:       "myrolloutstatus",
			RolloutRevision:     1337,
			RolloutStatusReason: "myrolloutstatusreason",
			DockerImage:         "mydockerimage",
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
	}

	result := mapDeploymentStatusToModel(deploymentStatus)

	assert.Equal(t, "myapp", result.AppName)
	assert.Equal(t, "mydeployment", result.DeploymentName)
	assert.Equal(t, "mystage", result.StageName)
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

func Test_mapDeploymentStatusFromModel(t *testing.T) {
	deploymentStatus := &model.DeploymentStatus{
		AppName:             "myapp",
		DeploymentName:      "mydeployment",
		StageName:           "mystage",
		RolloutStatus:       "myrolloutstatus",
		RolloutRevision:     1337,
		RolloutStatusReason: "myrolloutstatusreason",
		DockerImage:         "mydockerimage",
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

	assert.Equal(t, "myapp", result.AppName)
	assert.Equal(t, "mydeployment", result.DeploymentName)
	assert.Equal(t, "mystage", result.StageName)
	assert.InDelta(t, now, result.Doc.LastUpdated.Unix(), 3)
	assert.Equal(t, "myrolloutstatus", result.Doc.RolloutStatus)
	assert.EqualValues(t, 1337, result.Doc.RolloutRevision)
	assert.Equal(t, "myrolloutstatusreason", result.Doc.RolloutStatusReason)
	assert.Equal(t, "mydockerimage", result.Doc.DockerImage)
	assert.Len(t, result.Doc.Problems, 2)
	assert.Equal(t, "myproblem1", result.Doc.Problems[0].Message)
	assert.Equal(t, 1, result.Doc.Problems[0].Count)
	assert.Equal(t, "myproblem2", result.Doc.Problems[1].Message)
	assert.Equal(t, 2, result.Doc.Problems[1].Count)
}
