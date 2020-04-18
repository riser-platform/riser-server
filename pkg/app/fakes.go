package app

import (
	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/pkg/core"
)

type FakeService struct {
	CheckAppNameFn func(id uuid.UUID, name *core.NamespacedName) error
}

func (f *FakeService) CheckAppName(id uuid.UUID, name *core.NamespacedName) error {
	return f.CheckAppNameFn(id, name)
}

func (f *FakeService) CreateApp(name *core.NamespacedName) (*core.App, error) {
	panic("NI!")
}
