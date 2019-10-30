package core

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("the object could not be found")

// ValidationError provides an error consumable by a client. This is safe to return to the API as the errorHandler is aware of this error and will handle it correctly.
type ValidationError struct {
	Message string
	// ValidationErrors is a map of errors. The key should be set to the field name.
	ValidationErrors validation.Errors
	// Internal is an optional error that should not be exposed to the client.
	Internal error
}

// NewValidationError creates a ValidationError conditional on the type of error passed in. It is expected that validationErrors
// is typically of type  ozzo-validation.Errors object, containing a key/value pair of string/errors. If validationErrors is another error,
// Internal will be set.
func NewValidationError(message string, validationErrors error) error {
	if ve, ok := validationErrors.(validation.Errors); ok {
		return &ValidationError{Message: message, ValidationErrors: ve}
	}

	return &ValidationError{Message: message, Internal: validationErrors}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.ValidationErrors.Error())
}
