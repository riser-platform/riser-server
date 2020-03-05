package core

import "github.com/google/uuid"

type DeploymentRepository interface {
	Create(newDeployment *DeploymentRecord) error
	Delete(name *NamespacedName, stageName string) error
	GetByReservation(reservationId uuid.UUID, stageName string) (*Deployment, error)
	GetByName(name *NamespacedName, stageName string) (*Deployment, error)
	FindByApp(appId uuid.UUID) ([]Deployment, error)
	UpdateStatus(name *NamespacedName, stageName string, status *DeploymentStatus) error
	// TODO: either name/namespace/stage or reservationId
	UpdateTraffic(name, stageName string, riserRevision int64, traffic TrafficConfig) error
	// TODO: either name/namespace/stage or reservationId
	IncrementRevision(name, stageName string) (int64, error)
	// TODO: either name/namespace/stage or reservationId
	RollbackRevision(name, stageName string, failedRevision int64) (int64, error)
}

type FakeDeploymentRepository struct {
	CreateFn                   func(newDeployment *DeploymentRecord) error
	CreateCallCount            int
	DeleteFn                   func(name *NamespacedName, stageName string) error
	DeleteCallCount            int
	GetByNameFn                func(name *NamespacedName, stageName string) (*Deployment, error)
	GetByReservationFn         func(reservationId uuid.UUID, stageName string) (*Deployment, error)
	GetByReservationCallCount  int
	FindByAppFn                func(uuid.UUID) ([]Deployment, error)
	IncrementRevisionFn        func(name, stageName string) (int64, error)
	IncrementRevisionCallCount int
	RollbackRevisionFn         func(name, stageName string, failedRevision int64) (int64, error)
	UpdateStatusFn             func(name *NamespacedName, stageName string, status *DeploymentStatus) error
	UpdateStatusCallCount      int
	UpdateTrafficFn            func(name, stageName string, riserRevision int64, traffic TrafficConfig) error
	UpdateTrafficCallCount     int
}

func (f *FakeDeploymentRepository) Create(newDeployment *DeploymentRecord) error {
	f.CreateCallCount++
	return f.CreateFn(newDeployment)
}

func (f *FakeDeploymentRepository) Delete(name *NamespacedName, stageName string) error {
	f.DeleteCallCount++
	return f.DeleteFn(name, stageName)
}

func (f *FakeDeploymentRepository) GetByName(name *NamespacedName, stageName string) (*Deployment, error) {
	return f.GetByNameFn(name, stageName)
}

func (f *FakeDeploymentRepository) GetByReservation(reservationId uuid.UUID, stageName string) (*Deployment, error) {
	f.GetByReservationCallCount++
	return f.GetByReservationFn(reservationId, stageName)
}

func (fake *FakeDeploymentRepository) FindByApp(appId uuid.UUID) ([]Deployment, error) {
	return fake.FindByAppFn(appId)
}

func (fake *FakeDeploymentRepository) IncrementRevision(deploymentName, stageName string) (int64, error) {
	fake.IncrementRevisionCallCount++
	return fake.IncrementRevisionFn(deploymentName, stageName)
}

func (fake *FakeDeploymentRepository) RollbackRevision(name, stageName string, failedRevision int64) (int64, error) {
	return fake.RollbackRevisionFn(name, stageName, failedRevision)
}

func (fake *FakeDeploymentRepository) UpdateStatus(name *NamespacedName, stageName string, status *DeploymentStatus) error {
	fake.UpdateStatusCallCount++
	return fake.UpdateStatusFn(name, stageName, status)
}

func (fake *FakeDeploymentRepository) UpdateTraffic(name, stageName string, riserRevision int64, traffic TrafficConfig) error {
	fake.UpdateTrafficCallCount++
	return fake.UpdateTrafficFn(name, stageName, riserRevision, traffic)
}
