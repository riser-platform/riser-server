package api

import (
	"errors"

	"github.com/labstack/echo/v4"
)

// ErrInvalidBindType occurs when the wrong type is sent to a custom data binder. This should never happen.
var ErrInvalidBindType = errors.New("The incoming type must match the struct type")

type DefaultApplier interface {
	ApplyDefaults() error
}

type DataBinder struct{}

func (b *DataBinder) Bind(i interface{}, c echo.Context) (err error) {
	db := new(echo.DefaultBinder)
	if err = db.Bind(i, c); err != nil {
		return err
	}

	if modelDataBinder, ok := i.(DefaultApplier); ok {
		return modelDataBinder.ApplyDefaults()
	}

	return nil
}
