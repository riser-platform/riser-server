package core

type DeploymentStatusRepository interface {
	FindByApp(appName string) ([]DeploymentStatus, error)
	Save(*DeploymentStatus) error
}

type FakeDeploymentStatusRepository struct {
	FindByAppFn   func(string) ([]DeploymentStatus, error)
	SaveFn        func(*DeploymentStatus) error
	SaveCallCount int
}

func (fake *FakeDeploymentStatusRepository) FindByApp(appName string) ([]DeploymentStatus, error) {
	return fake.FindByAppFn(appName)
}

func (fake *FakeDeploymentStatusRepository) Save(status *DeploymentStatus) error {
	fake.SaveCallCount++
	return fake.SaveFn(status)
}
