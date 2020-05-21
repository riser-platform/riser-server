package v1

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/pkg/util"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/api/v1/model"
)

func Test_mapDeploymentToStatusModel(t *testing.T) {
	deployment := &core.Deployment{
		DeploymentReservation: core.DeploymentReservation{
			AppId:     uuid.New(),
			Name:      "mydeployment",
			Namespace: "myns",
		},
		DeploymentRecord: core.DeploymentRecord{
			EnvironmentName: "myenv",
			RiserRevision:   4,
			Doc: core.DeploymentDoc{
				Status: &core.DeploymentStatus{
					ObservedRiserRevision:     3,
					LatestCreatedRevisionName: "rev2",
					LatestReadyRevisionName:   "rev1",
					Revisions: []core.DeploymentRevisionStatus{
						{
							Name:                 "rev1",
							AvailableReplicas:    1,
							RevisionStatus:       "myrevisionstatus",
							RevisionStatusReason: "myrevisionstatusreason",
							DockerImage:          "mydockerimage",
							RiserRevision:        3,
						},
						{
							Name:                 "rev2",
							AvailableReplicas:    1,
							RevisionStatus:       "myrevisionstatus2",
							RevisionStatusReason: "myrevisionstatusreason2",
							DockerImage:          "mydockerimage2",
							RiserRevision:        4,
						},
					},
					Traffic: []core.DeploymentTrafficStatus{
						{
							Percent:      util.PtrInt64(90),
							RevisionName: "rev1",
							Tag:          "r1",
						},
						{
							Percent:      util.PtrInt64(10),
							RevisionName: "rev2",
							Tag:          "r2",
						},
					},
				},
			},
		},
	}

	result := mapDeploymentToStatusModel(deployment)

	assert.Equal(t, deployment.AppId, result.AppId)
	assert.Equal(t, "mydeployment", result.DeploymentName)
	assert.Equal(t, "myns", result.Namespace)
	assert.Equal(t, "myenv", result.EnvironmentName)
	assert.Equal(t, int64(3), result.ObservedRiserRevision)
	assert.Equal(t, int64(4), result.RiserRevision)
	assert.Equal(t, "rev2", result.LatestCreatedRevisionName)
	assert.Equal(t, "rev1", result.LatestReadyRevisionName)

	// Traffic
	assert.Len(t, result.Traffic, 2)
	assert.Equal(t, "rev1", result.Traffic[0].RevisionName)
	assert.Equal(t, int64(90), *result.Traffic[0].Percent)
	assert.Equal(t, "r1", result.Traffic[0].Tag)
	assert.Equal(t, "rev2", result.Traffic[1].RevisionName)
	assert.Equal(t, int64(10), *result.Traffic[1].Percent)
	assert.Equal(t, "r2", result.Traffic[1].Tag)

	// Revisions
	assert.Len(t, result.Revisions, 2)
	assert.Equal(t, "rev1", result.Revisions[0].Name)
	assert.Equal(t, int32(1), result.Revisions[0].AvailableReplicas)
	assert.Equal(t, "myrevisionstatus", result.Revisions[0].RevisionStatus)
	assert.Equal(t, "myrevisionstatusreason", result.Revisions[0].RevisionStatusReason)
	assert.Equal(t, "mydockerimage", result.Revisions[0].DockerImage)
	assert.Equal(t, int64(3), result.Revisions[0].RiserRevision)
	assert.Equal(t, "rev2", result.Revisions[1].Name)
	assert.Equal(t, int32(1), result.Revisions[1].AvailableReplicas)
	assert.Equal(t, "myrevisionstatus2", result.Revisions[1].RevisionStatus)
	assert.Equal(t, "myrevisionstatusreason2", result.Revisions[1].RevisionStatusReason)
	assert.Equal(t, "mydockerimage2", result.Revisions[1].DockerImage)
	assert.Equal(t, int64(4), result.Revisions[1].RiserRevision)
}

func Test_mapDeploymentToStatusModel_NilStatus(t *testing.T) {
	deployment := &core.Deployment{
		DeploymentReservation: core.DeploymentReservation{
			AppId:     uuid.New(),
			Name:      "mydeployment",
			Namespace: "myns",
		},
		DeploymentRecord: core.DeploymentRecord{

			EnvironmentName: "myenv",
			Doc:             core.DeploymentDoc{},
		},
	}

	result := mapDeploymentToStatusModel(deployment)

	assert.Equal(t, deployment.AppId, result.AppId)
	assert.Equal(t, "mydeployment", result.DeploymentName)
	assert.Equal(t, "myns", result.Namespace)
	assert.Equal(t, "myenv", result.EnvironmentName)
}

func Test_mapDeploymentStatusFromModel(t *testing.T) {
	deploymentStatus := &model.DeploymentStatusMutable{
		ObservedRiserRevision:     3,
		LatestReadyRevisionName:   "rev1",
		LatestCreatedRevisionName: "rev2",
		Revisions: []model.DeploymentRevisionStatus{
			{
				Name:                 "rev1",
				AvailableReplicas:    1,
				RiserRevision:        2,
				RevisionStatus:       "myrevisionstatus",
				RevisionStatusReason: "myrevisionstatusreason",
				DockerImage:          "mydockerimage",
			},
			{
				Name:                 "rev2",
				AvailableReplicas:    1,
				RiserRevision:        3,
				RevisionStatus:       "myrevisionstatus2",
				RevisionStatusReason: "myrevisionstatusreason2",
				DockerImage:          "mydockerimage2",
			},
		},
		Traffic: []model.DeploymentTrafficStatus{
			{
				Percent:      util.PtrInt64(90),
				RevisionName: "rev1",
				Tag:          "r1",
			},
			{
				Percent:      util.PtrInt64(10),
				RevisionName: "rev2",
				Tag:          "r2",
			},
		},
	}

	now := time.Now().Unix()

	result := mapDeploymentStatusFromModel(deploymentStatus)

	assert.InDelta(t, now, result.LastUpdated.Unix(), 3)
	assert.Equal(t, int64(3), result.ObservedRiserRevision)
	assert.Equal(t, "rev2", result.LatestCreatedRevisionName)
	assert.Equal(t, "rev1", result.LatestReadyRevisionName)

	// Revisions
	assert.Len(t, result.Revisions, 2)
	assert.Equal(t, int64(2), result.Revisions[0].RiserRevision)
	assert.Equal(t, int32(1), result.Revisions[0].AvailableReplicas)
	assert.Equal(t, "myrevisionstatus", result.Revisions[0].RevisionStatus)
	assert.Equal(t, "myrevisionstatusreason", result.Revisions[0].RevisionStatusReason)
	assert.Equal(t, "mydockerimage", result.Revisions[0].DockerImage)
	assert.Equal(t, int64(3), result.Revisions[1].RiserRevision)
	assert.Equal(t, int32(1), result.Revisions[1].AvailableReplicas)
	assert.Equal(t, "myrevisionstatus2", result.Revisions[1].RevisionStatus)
	assert.Equal(t, "myrevisionstatusreason2", result.Revisions[1].RevisionStatusReason)
	assert.Equal(t, "mydockerimage2", result.Revisions[1].DockerImage)

	// Traffic
	assert.Len(t, result.Traffic, 2)
	assert.Equal(t, "rev1", result.Traffic[0].RevisionName)
	assert.Equal(t, int64(90), *result.Traffic[0].Percent)
	assert.Equal(t, "r1", result.Traffic[0].Tag)
	assert.Equal(t, "rev2", result.Traffic[1].RevisionName)
	assert.Equal(t, int64(10), *result.Traffic[1].Percent)
	assert.Equal(t, "r2", result.Traffic[1].Tag)
}
