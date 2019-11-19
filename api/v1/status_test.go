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
				LatestReadyRevisionName: "rev2",
				Revisions: []core.DeploymentRevisionStatus{
					core.DeploymentRevisionStatus{
						Name:                "rev1",
						RolloutStatus:       "myrolloutstatus",
						RolloutStatusReason: "myrolloutstatusreason",
						DockerImage:         "mydockerimage",
						RiserGeneration:     3,
					},
					core.DeploymentRevisionStatus{
						Name:                "rev2",
						RolloutStatus:       "myrolloutstatus2",
						RolloutStatusReason: "myrolloutstatusreason2",
						DockerImage:         "mydockerimage2",
						RiserGeneration:     4,
					},
				},
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
	assert.Equal(t, int64(3), result.ObservedRiserGeneration)
	assert.Equal(t, int64(4), result.RiserGeneration)
	assert.Equal(t, "rev2", result.LatestReadyRevisionName)

	// Revisions
	assert.Len(t, result.Revisions, 2)
	assert.Equal(t, "rev1", result.Revisions[0].Name)
	assert.Equal(t, "myrolloutstatus", result.Revisions[0].RolloutStatus)
	assert.Equal(t, "myrolloutstatusreason", result.Revisions[0].RolloutStatusReason)
	assert.Equal(t, "mydockerimage", result.Revisions[0].DockerImage)
	assert.Equal(t, int64(3), result.Revisions[0].RiserGeneration)
	assert.Equal(t, "rev2", result.Revisions[1].Name)
	assert.Equal(t, "myrolloutstatus2", result.Revisions[1].RolloutStatus)
	assert.Equal(t, "myrolloutstatusreason2", result.Revisions[1].RolloutStatusReason)
	assert.Equal(t, "mydockerimage2", result.Revisions[1].DockerImage)
	assert.Equal(t, int64(4), result.Revisions[1].RiserGeneration)

	// Problems
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
}

func Test_mapDeploymentStatusFromModel(t *testing.T) {
	deploymentStatus := &model.DeploymentStatusMutable{
		ObservedRiserGeneration: 3,
		LatestReadyRevisionName: "rev2",
		Revisions: []model.DeploymentRevisionStatus{
			model.DeploymentRevisionStatus{
				Name:                "rev1",
				RiserGeneration:     2,
				RolloutStatus:       "myrolloutstatus",
				RolloutStatusReason: "myrolloutstatusreason",
				DockerImage:         "mydockerimage",
			},
			model.DeploymentRevisionStatus{
				Name:                "rev2",
				RiserGeneration:     3,
				RolloutStatus:       "myrolloutstatus2",
				RolloutStatusReason: "myrolloutstatusreason2",
				DockerImage:         "mydockerimage2",
			},
		},
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
	assert.Equal(t, "rev2", result.LatestReadyRevisionName)

	// Revisions
	assert.Len(t, result.Revisions, 2)
	assert.Equal(t, int64(2), result.Revisions[0].RiserGeneration)
	assert.Equal(t, "myrolloutstatus", result.Revisions[0].RolloutStatus)
	assert.Equal(t, "myrolloutstatusreason", result.Revisions[0].RolloutStatusReason)
	assert.Equal(t, "mydockerimage", result.Revisions[0].DockerImage)
	assert.Equal(t, int64(3), result.Revisions[1].RiserGeneration)
	assert.Equal(t, "myrolloutstatus2", result.Revisions[1].RolloutStatus)
	assert.Equal(t, "myrolloutstatusreason2", result.Revisions[1].RolloutStatusReason)
	assert.Equal(t, "mydockerimage2", result.Revisions[1].DockerImage)

	// Problems
	assert.Len(t, result.Problems, 2)
	assert.Equal(t, "myproblem1", result.Problems[0].Message)
	assert.Equal(t, 1, result.Problems[0].Count)
	assert.Equal(t, "myproblem2", result.Problems[1].Message)
	assert.Equal(t, 2, result.Problems[1].Count)
}
