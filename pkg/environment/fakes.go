package environment

import (
	"github.com/riser-platform/riser-server/pkg/core"
)

type FakeService struct {
	PingFn               func(string) error
	PingCallCount        int
	GetStatusFn          func(envName string) (*core.EnvironmentStatus, error)
	GetStatusCallCount   int
	ValidateDeployableFn func(envName string) error
}

func (fake *FakeService) Ping(envName string) error {
	fake.PingCallCount++
	return fake.PingFn(envName)
}

func (fake *FakeService) GetStatus(envName string) (*core.EnvironmentStatus, error) {
	fake.GetStatusCallCount++
	return fake.GetStatusFn(envName)
}

func (fake *FakeService) GetConfig(string) (*core.EnvironmentConfig, error) {
	panic("NI")
}

func (fake *FakeService) SetConfig(string, *core.EnvironmentConfig) error {
	panic("NI")
}

func (fake *FakeService) ValidateDeployable(envName string) error {
	return fake.ValidateDeployableFn(envName)
}
