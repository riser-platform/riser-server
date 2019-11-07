package app

import (
	"crypto/sha1"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
)

var ErrAlreadyExists = errors.New("app already exists")
var ErrInvalidAppId = errors.New("invalid app ID")
var ErrAppNotFound = errors.New("app not found")

const appIdSizeInBytes = 4

type Service interface {
	CreateApp(name string) (*core.App, error)
	CheckAppId(name string, appId core.AppId) error
}

type service struct {
	apps core.AppRepository
}

func NewService(apps core.AppRepository) Service {
	return &service{apps}
}

func (s *service) CreateApp(name string) (*core.App, error) {
	_, err := s.apps.Get(name)

	if err == nil {
		return nil, ErrAlreadyExists
	} else if err != core.ErrNotFound {
		return nil, errors.Wrap(err, "unable to validate app")
	}
	appId := createAppId()
	app := &core.App{
		Hashid: appId,
		Name:   name,
	}
	err = s.apps.Create(app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (s *service) CheckAppId(name string, appId core.AppId) error {
	app, err := s.apps.Get(name)
	if err != nil {
		if err == core.ErrNotFound {
			return ErrAppNotFound
		}
		return errors.Wrap(err, "Error getting app")
	}

	if appId.String() != app.Hashid.String() {
		return ErrInvalidAppId
	}

	return nil
}

func createAppId() core.AppId {
	hashBytes := sha1.Sum([]byte(uuid.New().String()))
	return hashBytes[:appIdSizeInBytes]
}
