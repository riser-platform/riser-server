package api

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/riser-platform/riser-server/pkg/core"
)

// ErrInvalidBindType occurs when the wrong type is sent to a custom data binder. This should never happen.
var ErrInvalidBindType = errors.New("The incoming type must match the struct type")

type DefaultApplier interface {
	ApplyDefaults() error
}

type Validator interface {
	Validate() error
}

type DataBinder struct{}

func (b *DataBinder) Bind(i interface{}, c echo.Context) (err error) {
	db := new(echo.DefaultBinder)
	if err = db.Bind(i, c); err != nil {
		return err
	}

	if modelWithDefaults, ok := i.(DefaultApplier); ok {
		err = modelWithDefaults.ApplyDefaults()
		if err != nil {
			return err
		}
	}

	// Use our own validation instead of echo validation.
	if modelWithValidator, ok := i.(Validator); ok {
		err := modelWithValidator.Validate()
		if err != nil {
			return core.NewValidationError("validation error", err)
		}
	}

	return nil
}
