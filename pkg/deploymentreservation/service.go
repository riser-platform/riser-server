package deploymentreservation

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
)

var ErrNameAlreadyReserved = core.NewValidationErrorMessage("the deployment name is already reserved in this namespace by another application")

type Service interface {
	EnsureReservation(appId uuid.UUID, name *core.NamespacedName) (*core.DeploymentReservation, error)
}

type service struct {
	reservations core.DeploymentReservationRepository
}

func NewService(reservations core.DeploymentReservationRepository) Service {
	return &service{reservations}
}

func (s *service) EnsureReservation(appId uuid.UUID, name *core.NamespacedName) (*core.DeploymentReservation, error) {
	reservation, err := s.reservations.GetByName(name)
	if err == core.ErrNotFound {
		reservation = &core.DeploymentReservation{
			Id:        uuid.New(),
			AppId:     appId,
			Name:      name.Name,
			Namespace: name.Namespace,
		}
		err = s.reservations.Create(reservation)
		if err != nil {
			return nil, errors.Wrap(err, "error creating reservation")
		}
	} else if err != nil {
		return nil, errors.Wrap(err, "error retrieving reservation")
	}

	if reservation.AppId != appId {
		return nil, ErrNameAlreadyReserved
	}

	return reservation, nil
}
