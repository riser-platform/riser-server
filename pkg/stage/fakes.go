package stage

import (
	"github.com/riser-platform/riser-server/pkg/core"
)

type FakeService struct {
	PingFn             func(string) error
	PingCallCount      int
	GetStatusFn        func(stageName string) (*core.StageStatus, error)
	GetStatusCallCount int
}

func (fake *FakeService) Ping(stageName string) error {
	fake.PingCallCount++
	return fake.PingFn(stageName)
}

func (fake *FakeService) GetStatus(stageName string) (*core.StageStatus, error) {
	fake.GetStatusCallCount++
	return fake.GetStatusFn(stageName)
}

func (fake *FakeService) SetConfig(stageName string, stage *core.StageConfig) error {
	panic("NI")
}
