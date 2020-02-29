package deployment

import (
	"time"

	"github.com/riser-platform/riser-server/pkg/deploymentreservation"
	"github.com/riser-platform/riser-server/pkg/state"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
)

// Note: See snapshot_test for state based testing of deployment artifacts

func Test_Delete(t *testing.T) {
	name := core.NewNamespacedName("mydep", "apps")
	deploymentRepository := &core.FakeDeploymentRepository{
		DeleteFn: func(nameArg *core.NamespacedName, stageName string) error {
			assert.Equal(t, name, nameArg)
			assert.Equal(t, "mystage", stageName)
			return nil
		},
	}

	committer := state.NewDryRunCommitter()

	service := service{deployments: deploymentRepository}

	err := service.Delete(name, "mystage", committer)

	assert.NoError(t, err)
	assert.Equal(t, 1, deploymentRepository.DeleteCallCount)
	assert.Len(t, committer.Commits, 1)
	assert.Equal(t, `Deleting deployment "mydep.apps"`, committer.Commits[0].Message)
	assert.Len(t, committer.Commits[0].Files, 2)
	assert.Equal(t, "stages/mystage/kube-resources/riser-managed/apps/deployments/mydep", committer.Commits[0].Files[0].Name)
	assert.True(t, committer.Commits[0].Files[0].Delete)
	assert.Equal(t, "stages/mystage/configs/apps/mydep.yaml", committer.Commits[0].Files[1].Name)
	assert.True(t, committer.Commits[0].Files[1].Delete)
}

func Test_Delete_SoftDeleteFails(t *testing.T) {
	deploymentRepository := &core.FakeDeploymentRepository{
		DeleteFn: func(*core.NamespacedName, string) error {
			return errors.New("test")
		},
	}

	committer := state.NewDryRunCommitter()

	service := service{deployments: deploymentRepository}

	err := service.Delete(core.NewNamespacedName("mydep", "myns"), "mystage", committer)

	assert.Equal(t, "error deleting deployment: test", err.Error())
}

func Test_Delete_DeploymentNotFound(t *testing.T) {
	deploymentRepository := &core.FakeDeploymentRepository{
		DeleteFn: func(*core.NamespacedName, string) error {
			return core.ErrNotFound
		},
	}

	service := service{deployments: deploymentRepository}

	err := service.Delete(core.NewNamespacedName("mydep", "myns"), "mystage", nil)

	assert.Equal(t, `There is no deployment by the name "mydep.myns" in stage "mystage"`, err.Error())
	assert.IsType(t, &core.ValidationError{}, err)
}

func Test_prepareForDeployment_whenNewDeploymentCreates(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:      "myapp-mydep",
		Namespace: "myns",
		Stage:     "mystage",
		App: &model.AppConfig{
			Id:   uuid.New(),
			Name: "myapp",
		},
	}

	reservation := &core.DeploymentReservation{
		Id: uuid.New(),
	}

	reservationService := &deploymentreservation.FakeService{
		EnsureReservationFn: func(appIdArg uuid.UUID, nameArg *core.NamespacedName) (*core.DeploymentReservation, error) {
			assert.Equal(t, deployment.App.Id, appIdArg)
			assert.Equal(t, core.NewNamespacedName(deployment.Name, deployment.Namespace), nameArg)
			return reservation, nil
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetByReservationFn: func(reservationId uuid.UUID, stageNameArg string) (*core.Deployment, error) {
			assert.Equal(t, reservation.Id, reservationId)
			assert.Equal(t, "mystage", stageNameArg)
			return nil, core.ErrNotFound
		},
		CreateFn: func(deploymentArg *core.DeploymentRecord) error {
			assert.NotEqual(t, uuid.Nil, deploymentArg.Id)
			assert.Equal(t, reservation.Id, deploymentArg.ReservationId)
			assert.Equal(t, "mystage", deploymentArg.StageName)
			assert.Equal(t, int64(1), deploymentArg.RiserRevision)
			return nil
		},
	}

	service := service{deployments: deploymentRepository, reservationService: reservationService}
	result, err := service.prepareForDeployment(deployment, false)

	assert.NoError(t, err)
	assert.Equal(t, "myns", deployment.Namespace)
	assert.Equal(t, int64(1), result)
	assert.Equal(t, 1, deploymentRepository.GetByReservationCallCount)
	assert.Equal(t, 1, deploymentRepository.CreateCallCount)
}

func Test_prepareForDeployment_whenExistingDeployment(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:      "myapp-mydep",
		Namespace: "myns",
		Stage:     "mystage",
		App: &model.AppConfig{
			Id:   uuid.New(),
			Name: "myapp",
		},
	}

	deploymentId := uuid.New()
	reservation := core.DeploymentReservation{
		Id:        uuid.New(),
		AppId:     deployment.App.Id,
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
	}

	reservationService := &deploymentreservation.FakeService{
		EnsureReservationFn: func(appIdArg uuid.UUID, nameArg *core.NamespacedName) (*core.DeploymentReservation, error) {
			return &reservation, nil
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetByReservationFn: func(reservationId uuid.UUID, stageNameArg string) (*core.Deployment, error) {
			return &core.Deployment{
				DeploymentReservation: reservation,
				DeploymentRecord: core.DeploymentRecord{
					Id:            deploymentId,
					ReservationId: reservation.Id,
					StageName:     "mystage"}}, nil
		},
		IncrementRevisionFn: func(name string, stageName string) (int64, error) {
			assert.Equal(t, "myapp-mydep", name)
			assert.Equal(t, "mystage", stageName)
			return 3, nil
		},
		UpdateTrafficFn: func(name string, stageName string, riserRevision int64, traffic core.TrafficConfig) error {
			assert.Equal(t, "myapp-mydep", name)
			assert.Equal(t, "mystage", stageName)
			assert.Len(t, traffic, 1)
			assert.Equal(t, int64(3), traffic[0].RiserRevision)
			assert.Equal(t, "myapp-mydep-3", traffic[0].RevisionName)
			assert.Equal(t, 100, traffic[0].Percent)
			return nil
		},
	}

	service := service{deployments: deploymentRepository, reservationService: reservationService}
	result, err := service.prepareForDeployment(deployment, false)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), result)
	assert.Equal(t, 1, deploymentRepository.GetByReservationCallCount)
	assert.Equal(t, 1, deploymentRepository.IncrementRevisionCallCount)
	assert.Equal(t, 1, deploymentRepository.UpdateTrafficCallCount)
	assert.Equal(t, 0, deploymentRepository.CreateCallCount)
}

// If a manual rollout is requested for a previously deleted deployment, don't try to update traffic rules with
// the old deployment as they will not be valid. ManualRollout is effectively ignored in this case.
func Test_prepareForDeployment_manualRollout_previouslyDeletedDeployment(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "myapp-mydep",
		Stage: "mystage",
		App: &model.AppConfig{
			Id:   uuid.New(),
			Name: "myapp",
		},
		ManualRollout: true,
	}

	deploymentId := uuid.New()
	reservation := core.DeploymentReservation{
		Id:        uuid.New(),
		AppId:     deployment.App.Id,
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
	}

	reservationService := &deploymentreservation.FakeService{
		EnsureReservationFn: func(appIdArg uuid.UUID, nameArg *core.NamespacedName) (*core.DeploymentReservation, error) {
			assert.Equal(t, deployment.App.Id, appIdArg)
			assert.Equal(t, core.NewNamespacedName(deployment.Name, deployment.Namespace), nameArg)
			return &reservation, nil
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetByReservationFn: func(reservationId uuid.UUID, stageNameArg string) (*core.Deployment, error) {
			deletedAt := time.Now()
			return &core.Deployment{
				DeploymentReservation: reservation,
				DeploymentRecord: core.DeploymentRecord{
					Id:            deploymentId,
					ReservationId: reservation.Id,
					StageName:     "mystage",
					DeletedAt:     &deletedAt,
					Doc: core.DeploymentDoc{
						// This rule should be ignored since the deployment was previously deleted
						Traffic: core.TrafficConfig{
							core.TrafficConfigRule{
								RevisionName:  "myapp-mydep-2",
								Percent:       100,
								RiserRevision: 2,
							},
						},
					},
				},
			}, nil
		},
		IncrementRevisionFn: func(name string, stageName string) (int64, error) {
			assert.Equal(t, "myapp-mydep", name)
			assert.Equal(t, "mystage", stageName)
			return 3, nil
		},
		UpdateTrafficFn: func(name string, stageName string, riserRevision int64, traffic core.TrafficConfig) error {
			assert.Equal(t, "myapp-mydep", name)
			assert.Equal(t, "mystage", stageName)
			// Even though a manual rollout is requested, a previously deleted deployment is treated as if there are no previous traffic rules
			// Therefore we route all traffic to the new revision.
			assert.Len(t, traffic, 1)
			assert.Equal(t, int64(3), traffic[0].RiserRevision)
			assert.Equal(t, "myapp-mydep-3", traffic[0].RevisionName)
			assert.Equal(t, 100, traffic[0].Percent)
			return nil
		},
	}

	service := service{deployments: deploymentRepository, reservationService: reservationService}
	result, err := service.prepareForDeployment(deployment, false)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), result)
	assert.Equal(t, 1, deploymentRepository.GetByReservationCallCount)
	assert.Equal(t, 1, deploymentRepository.IncrementRevisionCallCount)
	assert.Equal(t, 1, deploymentRepository.UpdateTrafficCallCount)
	assert.Equal(t, 0, deploymentRepository.CreateCallCount)
}

func Test_prepareForDeployment_whenIncrementRevisionFails(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "myapp-mydep",
		Stage: "mystage",
		App: &model.AppConfig{
			Id:   uuid.New(),
			Name: "myapp",
		},
	}

	deploymentId := uuid.New()
	reservation := core.DeploymentReservation{
		Id:        uuid.New(),
		AppId:     deployment.App.Id,
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
	}

	reservationService := &deploymentreservation.FakeService{
		EnsureReservationFn: func(appIdArg uuid.UUID, nameArg *core.NamespacedName) (*core.DeploymentReservation, error) {
			return &reservation, nil
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetByReservationFn: func(reservationId uuid.UUID, stageNameArg string) (*core.Deployment, error) {
			return &core.Deployment{
				DeploymentReservation: reservation,
				DeploymentRecord: core.DeploymentRecord{
					Id:            deploymentId,
					ReservationId: reservation.Id,
					StageName:     "mystage"}}, nil
		},
		IncrementRevisionFn: func(name string, stageName string) (int64, error) {
			return 0, errors.New("test")
		},
	}

	service := service{deployments: deploymentRepository, reservationService: reservationService}
	result, err := service.prepareForDeployment(deployment, false)

	assert.Zero(t, result)
	assert.Equal(t, "Error incrementing deployment revision: test", err.Error())
}

func Test_prepareForDeployment_doesNotUpdateWhenDryRun(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "myapp-mydep",
		Stage: "mystage",
		App: &model.AppConfig{
			Id:   uuid.New(),
			Name: "myapp",
		},
	}

	deploymentId := uuid.New()
	reservation := core.DeploymentReservation{
		Id:        uuid.New(),
		AppId:     deployment.App.Id,
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
	}

	reservationService := &deploymentreservation.FakeService{
		EnsureReservationFn: func(appIdArg uuid.UUID, nameArg *core.NamespacedName) (*core.DeploymentReservation, error) {
			return &reservation, nil
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetByReservationFn: func(reservationId uuid.UUID, stageNameArg string) (*core.Deployment, error) {
			return &core.Deployment{
				DeploymentReservation: reservation,
				DeploymentRecord: core.DeploymentRecord{
					Id:            deploymentId,
					ReservationId: reservation.Id,
					StageName:     "mystage"}}, nil
		},
	}

	service := service{deployments: deploymentRepository, reservationService: reservationService}
	result, err := service.prepareForDeployment(deployment, true)

	assert.NoError(t, err)
	// The RiserRevision is always "0" for a dry-run
	assert.Equal(t, int64(0), result)
	assert.Equal(t, 1, deploymentRepository.GetByReservationCallCount)
	assert.Equal(t, 0, deploymentRepository.IncrementRevisionCallCount)
	assert.Equal(t, 0, deploymentRepository.UpdateTrafficCallCount)
	assert.Equal(t, 0, deploymentRepository.CreateCallCount)
	// Traffic should still be computed in a dry-run, just not persisted
	assert.Len(t, deployment.Traffic, 1)
	assert.Equal(t, int64(0), deployment.Traffic[0].RiserRevision)
	assert.Equal(t, "myapp-mydep-0", deployment.Traffic[0].RevisionName)
	assert.Equal(t, 100, deployment.Traffic[0].Percent)
}

func Test_prepareForDeployment_whenUpdateTrafficFails(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "myapp-mydep",
		Stage: "mystage",
		App: &model.AppConfig{
			Id:   uuid.New(),
			Name: "myapp",
		},
	}

	deploymentId := uuid.New()
	reservation := core.DeploymentReservation{
		Id:        uuid.New(),
		AppId:     deployment.App.Id,
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
	}

	reservationService := &deploymentreservation.FakeService{
		EnsureReservationFn: func(appIdArg uuid.UUID, nameArg *core.NamespacedName) (*core.DeploymentReservation, error) {
			return &reservation, nil
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetByReservationFn: func(reservationId uuid.UUID, stageNameArg string) (*core.Deployment, error) {
			return &core.Deployment{
				DeploymentReservation: reservation,
				DeploymentRecord: core.DeploymentRecord{
					Id:            deploymentId,
					ReservationId: reservation.Id,
					StageName:     "mystage"}}, nil
		},
		IncrementRevisionFn: func(name string, stageName string) (int64, error) {
			return 1, nil
		},
		UpdateTrafficFn: func(name string, stageName string, riserRevision int64, traffic core.TrafficConfig) error {
			return errors.New("broke")
		},
	}

	service := service{deployments: deploymentRepository, reservationService: reservationService}
	result, err := service.prepareForDeployment(deployment, false)

	assert.Zero(t, result)
	assert.Equal(t, "Error updating traffic: broke", err.Error())
}

func Test_prepareForDeployment_whenEnsureReservationErr(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "myapp-mydep",
		Stage: "mystage",
		App: &model.AppConfig{
			Id:   uuid.New(),
			Name: "myapp",
		},
	}

	reservationService := &deploymentreservation.FakeService{
		EnsureReservationFn: func(appIdArg uuid.UUID, nameArg *core.NamespacedName) (*core.DeploymentReservation, error) {
			return nil, errors.New("test")
		},
	}

	service := service{reservationService: reservationService}
	result, err := service.prepareForDeployment(deployment, false)

	assert.Zero(t, result)
	assert.Equal(t, `Error ensuring deployment reservation: test`, err.Error())
}

func Test_prepareForDeployment_whenGetFails(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "myapp-mydep",
		Stage: "mystage",
		App: &model.AppConfig{
			Name: "myapp",
		},
	}

	reservation := core.DeploymentReservation{
		Id:        uuid.New(),
		AppId:     deployment.App.Id,
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
	}

	reservationService := &deploymentreservation.FakeService{
		EnsureReservationFn: func(appIdArg uuid.UUID, nameArg *core.NamespacedName) (*core.DeploymentReservation, error) {
			return &reservation, nil
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetByReservationFn: func(uuid.UUID, string) (*core.Deployment, error) {
			return nil, errors.New("test")
		},
	}

	service := service{deployments: deploymentRepository, reservationService: reservationService}
	result, err := service.prepareForDeployment(deployment, false)

	assert.Zero(t, result)
	assert.Equal(t, `Error retrieving deployment "myapp-mydep" in stage "mystage": test`, err.Error())
}

func Test_prepareForDeployment_whenCreateFails(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "myapp-mydep",
		Stage: "mystage",
		App: &model.AppConfig{
			Name: "myapp",
		},
	}

	reservation := core.DeploymentReservation{
		Id:        uuid.New(),
		AppId:     deployment.App.Id,
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
	}

	reservationService := &deploymentreservation.FakeService{
		EnsureReservationFn: func(appIdArg uuid.UUID, nameArg *core.NamespacedName) (*core.DeploymentReservation, error) {
			return &reservation, nil
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetByReservationFn: func(reservationId uuid.UUID, stageNameArg string) (*core.Deployment, error) {
			return nil, core.ErrNotFound
		},
		CreateFn: func(newDeploymentArg *core.DeploymentRecord) error {
			return errors.New("test")
		},
	}

	service := service{deployments: deploymentRepository, reservationService: reservationService}
	result, err := service.prepareForDeployment(deployment, false)

	assert.Zero(t, result)
	assert.Equal(t, `Error creating deployment "myapp-mydep" in stage "mystage": test`, err.Error())
}

func Test_computeTraffic_NewDeployment(t *testing.T) {
	cfg := &core.DeploymentConfig{
		Name: "myapp",
	}

	result := computeTraffic(1, cfg, nil)

	assert.Len(t, result, 1)
	assert.EqualValues(t, result[0].RiserRevision, 1)
	assert.Equal(t, result[0].RevisionName, "myapp-1")
	assert.EqualValues(t, result[0].Percent, 100)
}

// Manual rollout for a first time deployment is effectively not allowed
func Test_computeTraffic_NewDeployment_ManualRollout(t *testing.T) {
	cfg := &core.DeploymentConfig{
		Name:          "myapp",
		ManualRollout: true,
	}

	result := computeTraffic(1, cfg, nil)

	assert.Len(t, result, 1)
	assert.EqualValues(t, result[0].RiserRevision, 1)
	assert.Equal(t, result[0].RevisionName, "myapp-1")
	assert.EqualValues(t, result[0].Percent, 100)
}

func Test_computeTraffic_ExistingDeployment_ManualRollout(t *testing.T) {
	cfg := &core.DeploymentConfig{
		Name:          "myapp",
		ManualRollout: true,
	}

	existingDeployment := &core.DeploymentRecord{
		Doc: core.DeploymentDoc{
			Traffic: core.TrafficConfig{
				core.TrafficConfigRule{
					RiserRevision: 1,
					RevisionName:  "myapp-1",
					Percent:       100,
				},
			},
		},
	}

	result := computeTraffic(2, cfg, existingDeployment)

	assert.Len(t, result, 2)
	assert.EqualValues(t, result[0].RiserRevision, 2)
	assert.Equal(t, result[0].RevisionName, "myapp-2")
	assert.EqualValues(t, result[0].Percent, 0)
	assert.EqualValues(t, result[1].RiserRevision, 1)
	assert.Equal(t, result[1].RevisionName, "myapp-1")
	assert.EqualValues(t, result[1].Percent, 100)
}

func Test_computeTraffic_ExistingDeployment_ManualRollout_RemovesExistingZeroPercentRules(t *testing.T) {
	cfg := &core.DeploymentConfig{
		Name:          "myapp",
		ManualRollout: true,
	}

	existingDeployment := &core.DeploymentRecord{
		Doc: core.DeploymentDoc{
			Traffic: core.TrafficConfig{
				core.TrafficConfigRule{
					RiserRevision: 1,
					RevisionName:  "myapp-1",
					Percent:       100,
				},
				core.TrafficConfigRule{
					RiserRevision: 2,
					RevisionName:  "myapp-2",
					Percent:       0,
				},
			},
		},
	}

	result := computeTraffic(3, cfg, existingDeployment)

	assert.Len(t, result, 2)
	assert.EqualValues(t, result[0].RiserRevision, 3)
	assert.Equal(t, result[0].RevisionName, "myapp-3")
	assert.EqualValues(t, result[0].Percent, 0)
	assert.EqualValues(t, result[1].RiserRevision, 1)
	assert.Equal(t, result[1].RevisionName, "myapp-1")
	assert.EqualValues(t, result[1].Percent, 100)
}

func Test_validateDeploymentConfig_ValidatesName(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"app", nil},
		{"app-good", nil},
		{"app-", core.NewValidationError(`invalid deployment name "app-"`, errors.New(`must be either "app" or start with "app-"`))},
		{"mydep", core.NewValidationError(`invalid deployment name "mydep"`, errors.New(`must be either "app" or start with "app-"`))},
		{"app-b@d", core.NewValidationError(`invalid deployment name "app-b@d"`, errors.New("must be lowercase, alphanumeric, and start with a letter"))},
	}

	for _, tt := range tests {
		deployment := &core.DeploymentConfig{
			Name: tt.name,
			App: &model.AppConfig{
				Name: "app",
			},
		}

		result := validateDeploymentConfig(deployment)

		if tt.err == nil {
			assert.Nil(t, result, tt.name)
		} else {
			assert.IsType(t, &core.ValidationError{}, result, tt.name)
			ve := result.(*core.ValidationError)
			ttve := tt.err.(*core.ValidationError)
			assert.Equal(t, ttve.Error(), ve.Error(), tt.name)
		}
	}
}
