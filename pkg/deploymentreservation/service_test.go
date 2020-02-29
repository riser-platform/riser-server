package deploymentreservation

import (
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_EnsureReservation_ExistingReservation(t *testing.T) {
	appId := uuid.New()
	name := core.NewNamespacedName("mydep", "myns")
	expectedReservation := &core.DeploymentReservation{
		Id:        uuid.New(),
		AppId:     appId,
		Name:      name.Name,
		Namespace: name.Namespace,
	}

	reservations := &core.FakeDeploymentReservationRepository{
		GetByNameFn: func(nameArg *core.NamespacedName) (*core.DeploymentReservation, error) {
			assert.Equal(t, name, nameArg)
			return expectedReservation, nil
		},
	}

	svc := service{reservations}

	result, err := svc.EnsureReservation(appId, name)

	assert.NoError(t, err)
	assert.Equal(t, 0, reservations.CreateCallCount)
	assert.Equal(t, expectedReservation, result)
}

func Test_EnsureReservation_NewReservation(t *testing.T) {
	appId := uuid.New()
	name := core.NewNamespacedName("mydep", "myns")
	reservations := &core.FakeDeploymentReservationRepository{
		GetByNameFn: func(nameArg *core.NamespacedName) (*core.DeploymentReservation, error) {
			assert.Equal(t, name, nameArg)
			return nil, core.ErrNotFound
		},
		CreateFn: func(reservationArg *core.DeploymentReservation) error {
			assert.NotEqual(t, uuid.Nil, reservationArg.Id)
			assert.Equal(t, appId, reservationArg.AppId)
			assert.Equal(t, name.Name, reservationArg.Name)
			assert.Equal(t, name.Namespace, reservationArg.Namespace)
			return nil
		},
	}

	svc := service{reservations}

	result, err := svc.EnsureReservation(appId, name)

	assert.NoError(t, err)
	assert.Equal(t, 1, reservations.CreateCallCount)
	assert.NotEqual(t, uuid.Nil, result.Id)
	assert.Equal(t, appId, result.AppId)
	assert.Equal(t, name.Name, result.Name)
	assert.Equal(t, name.Namespace, result.Namespace)
}

func Test_EnsureReservation_ExistingReservationIsOwnedByAnotherApp(t *testing.T) {
	appId := uuid.New()
	name := core.NewNamespacedName("mydep", "myns")
	expectedReservation := &core.DeploymentReservation{
		Id:        uuid.New(),
		AppId:     uuid.New(),
		Name:      name.Name,
		Namespace: name.Namespace,
	}

	reservations := &core.FakeDeploymentReservationRepository{
		GetByNameFn: func(nameArg *core.NamespacedName) (*core.DeploymentReservation, error) {
			return expectedReservation, nil
		},
	}

	svc := service{reservations}

	result, err := svc.EnsureReservation(appId, name)

	assert.Nil(t, result)
	assert.Equal(t, ErrNameAlreadyReserved, err)
}

func Test_EnsureReservation_GetErr(t *testing.T) {
	reservations := &core.FakeDeploymentReservationRepository{
		GetByNameFn: func(nameArg *core.NamespacedName) (*core.DeploymentReservation, error) {
			return nil, errors.New("test")
		},
	}

	svc := service{reservations}

	result, err := svc.EnsureReservation(uuid.New(), core.NewNamespacedName("myapp", "myns"))

	assert.Nil(t, result)
	assert.Equal(t, "error retrieving reservation: test", err.Error())
}

func Test_EnsureReservation_CreateErr(t *testing.T) {
	reservations := &core.FakeDeploymentReservationRepository{
		GetByNameFn: func(nameArg *core.NamespacedName) (*core.DeploymentReservation, error) {
			return nil, core.ErrNotFound
		},

		CreateFn: func(reservation *core.DeploymentReservation) error {
			return errors.New("test")
		},
	}

	svc := service{reservations}

	result, err := svc.EnsureReservation(uuid.New(), core.NewNamespacedName("myapp", "myns"))

	assert.Nil(t, result)
	assert.Equal(t, "error creating reservation: test", err.Error())
}
