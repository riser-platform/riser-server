package app

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/namespace"
)

var (
	ErrAlreadyExists       = core.NewValidationErrorMessage("an app already exists with the provided name")
	ErrInvalidAppName      = core.NewValidationErrorMessage("the app name does not match the app ID: you may not change the app's name after creation")
	ErrInvalidAppNamespace = core.NewValidationErrorMessage("the app namespace does not match app ID: you may not change the app's namespace after creation")
	// TODO: Consider removing and coming up with a better client error for non-fatal not founds (core.ErrorNotFound is considered fatal by default)
	ErrAppNotFound = errors.New("app not found")
)

type Service interface {
	CreateApp(name *core.NamespacedName) (*core.App, error)
	// CheckAppName ensures that the app name and namespace belongs to the app ID. This prevents an accidental or otherwise name change in the app config.
	CheckAppName(id uuid.UUID, name *core.NamespacedName) error
}

type service struct {
	apps             core.AppRepository
	namespaceService namespace.Service
}

func NewService(apps core.AppRepository, namespaceService namespace.Service) Service {
	return &service{apps, namespaceService}
}

func (s *service) CreateApp(name *core.NamespacedName) (*core.App, error) {
	_, err := s.apps.GetByName(name)

	if err == nil {
		return nil, ErrAlreadyExists
	} else if err != core.ErrNotFound {
		return nil, errors.Wrap(err, "unable to validate app")
	}

	err = s.namespaceService.ValidateDeployable(name.Namespace)
	if err != nil {
		return nil, err
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
