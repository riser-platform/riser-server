package app

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
)

var ErrAlreadyExists = errors.New("an app already exists with the provided name")
var ErrInvalidAppName = errors.New("the app name does not match the name associated with the provided app ID")
var ErrInvalidAppNamespace = errors.New("the app namespace does not match the name associated with the provided app ID")
var ErrAppNotFound = errors.New("app not found")

type Service interface {
	CreateApp(name *core.NamespacedName) (*core.App, error)
	// CheckAppName ensures that the app name and namespace belongs to the app ID. This prevents an accidental or otherwise name change in the app config.
	CheckAppName(id uuid.UUID, name *core.NamespacedName) error
}

type service struct {
	apps core.AppRepository
}

func NewService(apps core.AppRepository) Service {
	return &service{apps}
}

func (s *service) CreateApp(name *core.NamespacedName) (*core.App, error) {
	_, err := s.apps.GetByName(name)

	if err == nil {
		return nil, ErrAlreadyExists
	} else if err != core.ErrNotFound {
		return nil, errors.Wrap(err, "unable to validate app")
	}
	appId := uuid.New()
	app := &core.App{
		Id:        appId,
		Name:      name.Name,
		Namespace: name.Namespace,
	}
	err = s.apps.Create(app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (s *service) CheckAppName(id uuid.UUID, name *core.NamespacedName) error {
	app, err := s.apps.Get(id)
	if err != nil {
		return handleGetAppErr(err)
	}

	if name.Name != app.Name {
		return ErrInvalidAppName
	}

	if name.Namespace != app.Namespace {
		return ErrInvalidAppNamespace
	}

	return nil
}

func handleGetAppErr(err error) error {
	if err == core.ErrNotFound {
		return ErrAppNotFound
	}
	return errors.Wrap(err, "Error getting app")
}
