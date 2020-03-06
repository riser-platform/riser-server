package core

type DeploymentReservationRepository interface {
	Create(reservation *DeploymentReservation) error
	GetByName(name *NamespacedName) (*DeploymentReservation, error)
}

type FakeDeploymentReservationRepository struct {
	CreateFn        func(reservation *DeploymentReservation) error
	CreateCallCount int
	GetByNameFn     func(name *NamespacedName) (*DeploymentReservation, error)
}

func (f *FakeDeploymentReservationRepository) Create(reservation *DeploymentReservation) error {
	f.CreateCallCount++
	return f.CreateFn(reservation)
}

func (f *FakeDeploymentReservationRepository) GetByName(name *NamespacedName) (*DeploymentReservation, error) {
	return f.GetByNameFn(name)
}
