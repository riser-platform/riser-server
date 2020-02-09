package core

type DeploymentRepository interface {
	Create(newDeployment *Deployment) error
	Delete(name, stageName string) error
	Get(name, stageName string) (*Deployment, error)
	FindByApp(appName string) ([]Deployment, error)
	UpdateStatus(name, stageName string, status *DeploymentStatus) error
	UpdateTraffic(name, stageName string, riserRevision int64, traffic TrafficConfig) error
	IncrementRevision(name, stageName string) (int64, error)
	RollbackRevision(name, stageName string, failedRevision int64) (int64, error)
}

type FakeDeploymentRepository struct {
	CreateFn                   func(newDeployment *Deployment) error
	CreateCallCount            int
	DeleteFn                   func(name, stageName string) error
	DeleteCallCount            int
	GetFn                      func(name, stageName string) (*Deployment, error)
	GetCallCount               int
	FindByAppFn                func(string) ([]Deployment, error)
	IncrementRevisionFn        func(name, stageName string) (int64, error)
	IncrementRevisionCallCount int
	RollbackRevisionFn         func(name, stageName string, failedRevision int64) (int64, error)
	UpdateStatusFn             func(name, stageName string, status *DeploymentStatus) error
	UpdateStatusCallCount      int
	UpdateTrafficFn            func(name, stageName string, riserRevision int64, traffic TrafficConfig) error
	UpdateTrafficCallCount     int
}

func (f *FakeDeploymentRepository) Create(newDeployment *Deployment) error {
	f.CreateCallCount++
	return f.CreateFn(newDeployment)
}

func (f *FakeDeploymentRepository) Delete(name, stageName string) error {
	f.DeleteCallCount++
	return f.DeleteFn(name, stageName)
}

func (f *FakeDeploymentRepository) Get(name, stageName string) (*Deployment, error) {
	f.GetCallCount++
	return f.GetFn(name, stageName)
}

func (fake *FakeDeploymentRepository) FindByApp(appName string) ([]Deployment, error) {
	return fake.FindByAppFn(appName)
}

func (fake *FakeDeploymentRepository) IncrementRevision(deploymentName, stageName string) (int64, error) {
	fake.IncrementRevisionCallCount++
	return fake.IncrementRevisionFn(deploymentName, stageName)
}

func (fake *FakeDeploymentRepository) RollbackRevision(name, stageName string, failedRevision int64) (int64, error) {
	return fake.RollbackRevisionFn(name, stageName, failedRevision)
}

func (fake *FakeDeploymentRepository) UpdateStatus(deploymentName, stageName string, status *DeploymentStatus) error {
	fake.UpdateStatusCallCount++
	return fake.UpdateStatusFn(deploymentName, stageName, status)
}

func (fake *FakeDeploymentRepository) UpdateTraffic(name, stageName string, riserRevision int64, traffic TrafficConfig) error {
	fake.UpdateTrafficCallCount++
	return fake.UpdateTrafficFn(name, stageName, riserRevision, traffic)
}
