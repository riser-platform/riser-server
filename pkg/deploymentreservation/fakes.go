package deploymentreservation

import (
	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/pkg/core"
)

type FakeService struct {
	EnsureReservationFn func(appId uuid.UUID, name *core.NamespacedName) (*core.DeploymentReservation, error)
}

func (f *FakeService) EnsureReservation(appId uuid.UUID, name *core.NamespacedName) (*core.DeploymentReservation, error) {
	return f.EnsureReservationFn(appId, name)
}
