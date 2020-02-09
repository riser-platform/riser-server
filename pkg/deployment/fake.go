package deployment

import (
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/state"
)

type FakeService struct {
	DeleteFn        func(deploymentName, namespace, stageName string, committer state.Committer) error
	DeleteCallCount int
}

func (f *FakeService) Update(deployment *core.DeploymentConfig, committer state.Committer, dryRun bool) error {
	panic("NI!")
}

func (f *FakeService) Delete(deploymentName, namespace, stageName string, committer state.Committer) error {
	f.DeleteCallCount++
	return f.DeleteFn(deploymentName, namespace, stageName, committer)
}
