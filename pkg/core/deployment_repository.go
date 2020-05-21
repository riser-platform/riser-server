package core

import "github.com/google/uuid"

type DeploymentRepository interface {
	Create(newDeployment *DeploymentRecord) error
	Delete(name *NamespacedName, envName string) error
	GetByReservation(reservationId uuid.UUID, envName string) (*Deployment, error)
	GetByName(name *NamespacedName, envName string) (*Deployment, error)
	FindByApp(appId uuid.UUID) ([]Deployment, error)
	UpdateStatus(name *NamespacedName, envName string, status *DeploymentStatus) error
	UpdateTraffic(name *NamespacedName, envName string, riserRevision int64, traffic TrafficConfig) error
	IncrementRevision(name *NamespacedName, envName string) (int64, error)
	RollbackRevision(name *NamespacedName, envName string, failedRevision int64) (int64, error)
}

type FakeDeploymentRepository struct {
	CreateFn                   func(newDeployment *DeploymentRecord) error
	CreateCallCount            int
	DeleteFn                   func(name *NamespacedName, envName string) error
	DeleteCallCount            int
	GetByNameFn                func(name *NamespacedName, envName string) (*Deployment, error)
	GetByReservationFn         func(reservationId uuid.UUID, envName string) (*Deployment, error)
	GetByReservationCallCount  int
	FindByAppFn                func(uuid.UUID) ([]Deployment, error)
	IncrementRevisionFn        func(name *NamespacedName, envName string) (int64, error)
	IncrementRevisionCallCount int
	RollbackRevisionFn         func(name *NamespacedName, envName string, failedRevision int64) (int64, error)
	UpdateStatusFn             func(name *NamespacedName, envName string, status *DeploymentStatus) error
	UpdateStatusCallCount      int
	UpdateTrafficFn            func(name *NamespacedName, envName string, riserRevision int64, traffic TrafficConfig) error
	UpdateTrafficCallCount     int
}

func (f *FakeDeploymentRepository) Create(newDeployment *DeploymentRecord) error {
	f.CreateCallCount++
	return f.CreateFn(newDeployment)
}

func (f *FakeDeploymentRepository) Delete(name *NamespacedName, envName string) error {
	f.DeleteCallCount++
	return f.DeleteFn(name, envName)
}

func (f *FakeDeploymentRepository) GetByName(name *NamespacedName, envName string) (*Deployment, error) {
	return f.GetByNameFn(name, envName)
}

func (f *FakeDeploymentRepository) GetByReservation(reservationId uuid.UUID, envName string) (*Deployment, error) {
	f.GetByReservationCallCount++
	return f.GetByReservationFn(reservationId, envName)
}

func (fake *FakeDeploymentRepository) FindByApp(appId uuid.UUID) ([]Deployment, error) {
	return fake.FindByAppFn(appId)
}

func (fake *FakeDeploymentRepository) IncrementRevision(name *NamespacedName, envName string) (int64, error) {
	fake.IncrementRevisionCallCount++
	return fake.IncrementRevisionFn(name, envName)
}

func (fake *FakeDeploymentRepository) RollbackRevision(name *NamespacedName, envName string, failedRevision int64) (int64, error) {
	return fake.RollbackRevisionFn(name, envName, failedRevision)
}

func (fake *FakeDeploymentRepository) UpdateStatus(name *NamespacedName, envName string, status *DeploymentStatus) error {
	fake.UpdateStatusCallCount++
	return fake.UpdateStatusFn(name, envName, status)
}

func (fake *FakeDeploymentRepository) UpdateTraffic(name *NamespacedName, envName string, riserRevision int64, traffic TrafficConfig) error {
	fake.UpdateTrafficCallCount++
	return fake.UpdateTrafficFn(name, envName, riserRevision, traffic)
}
