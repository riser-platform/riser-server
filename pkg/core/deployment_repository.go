package core

type DeploymentRepository interface {
	Create(newDeployment *Deployment) error
	FindByApp(appName string) ([]Deployment, error)
	Get(name, stageName string) (*Deployment, error)
	UpdateStatus(name, stageName string, status *DeploymentStatus) error
	IncrementGeneration(name, stageName string) (int64, error)
	RollbackGeneration(name, stageName string, failedGeneration int64) (int64, error)
}

type FakeDeploymentRepository struct {
	CreateFn                     func(newDeployment *Deployment) error
	CreateCallCount              int
	GetFn                        func(name, stageName string) (*Deployment, error)
	GetCallCount                 int
	FindByAppFn                  func(string) ([]Deployment, error)
	IncrementGenerationFn        func(name, stageName string) (int64, error)
	IncrementGenerationCallCount int
	RollbackGenerationFn         func(name, stageName string, failedGeneration int64) (int64, error)
	UpdateStatusFn               func(name, stageName string, status *DeploymentStatus) error
	UpdateStatusCallCount        int
}

func (f *FakeDeploymentRepository) Create(newDeployment *Deployment) error {
	f.CreateCallCount++
	return f.CreateFn(newDeployment)
}

func (f *FakeDeploymentRepository) Get(name, stageName string) (*Deployment, error) {
	f.GetCallCount++
	return f.GetFn(name, stageName)
}

func (fake *FakeDeploymentRepository) FindByApp(appName string) ([]Deployment, error) {
	return fake.FindByAppFn(appName)
}

func (fake *FakeDeploymentRepository) IncrementGeneration(deploymentName, stageName string) (int64, error) {
	fake.IncrementGenerationCallCount++
	return fake.IncrementGenerationFn(deploymentName, stageName)
}

func (fake *FakeDeploymentRepository) RollbackGeneration(name, stageName string, failedGeneration int64) (int64, error) {
	return fake.RollbackGenerationFn(name, stageName, failedGeneration)
}

func (fake *FakeDeploymentRepository) UpdateStatus(deploymentName, stageName string, status *DeploymentStatus) error {
	fake.UpdateStatusCallCount++
	return fake.UpdateStatusFn(deploymentName, stageName, status)
}
