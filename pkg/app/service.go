package app

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
)

var ErrAlreadyExists = errors.New("an app already exists with the provided name")
var ErrInvalidAppName = errors.New("the app name does not match the name associated with the provided app ID")
var ErrAppNotFound = errors.New("app not found")

type Service interface {
	GetByIdOrName(idOrName string) (*core.App, error)
	CreateApp(name string) (*core.App, error)
	// CheckAppName ensures that the app name belongs to the app ID. This prevents an accidental or otherwise name change in the app config.
	CheckAppName(id uuid.UUID, name string) error
}

type service struct {
	apps core.AppRepository
}

func NewService(apps core.AppRepository) Service {
	return &service{apps}
}

func (s *service) GetByIdOrName(idOrName string) (app *core.App, err error) {
	appId, _ := uuid.Parse(idOrName)
	if appId != uuid.Nil {
		app, err = s.apps.Get(appId)
	} else {
		app, err = s.apps.GetByName(idOrName)
	}
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (s *service) CreateApp(name string) (*core.App, error) {
	_, err := s.apps.GetByName(name)

	if err == nil {
		return nil, ErrAlreadyExists
	} else if err != core.ErrNotFound {
		return nil, errors.Wrap(err, "unable to validate app")
	}
	appId := uuid.New()
	app := &core.App{
		Id:   appId,
		Name: name,
	}
	err = s.apps.Create(app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (s *service) CheckAppName(id uuid.UUID, name string) error {
	app, err := s.apps.Get(id)
	if err != nil {
		return handleGetAppErr(err)
	}

	if name != app.Name {
		return ErrInvalidAppName
	}

	return nil
}

func handleGetAppErr(err error) error {
	if err == core.ErrNotFound {
		return ErrAppNotFound
	}
	return errors.Wrap(err, "Error getting app")
}
