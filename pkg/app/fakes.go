package app

import (
	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/pkg/core"
)

type FakeService struct {
	CheckIDFn func(id uuid.UUID, name *core.NamespacedName) error
}

func (f *FakeService) CheckID(id uuid.UUID, name *core.NamespacedName) error {
	return f.CheckIDFn(id, name)
}

func (f *FakeService) Create(name *core.NamespacedName) (*core.App, error) {
	panic("NI!")
}

func (f *FakeService) GetByName(name *core.NamespacedName) (*core.App, error) {
	panic("NI!")
}
