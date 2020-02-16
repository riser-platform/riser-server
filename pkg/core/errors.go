package core

import (
	"fmt"

	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("the object could not be found")
var ErrConflictNewerVersion = errors.New("a newer version of the object exists")

// ValidationError provides an error consumable by a client. This is safe to return to the API as the errorHandler is aware of this error
type ValidationError struct {
	error
	Message string
	// ValidationError represents the validation error that ocurred. See the API errorHandler for how this is returned to the client.
	ValidationError error
}

// NewValidationError creates a ValidationError conditional on the type of error passed in. It is expected that validationError
// is typically of type ozzo-validation.Errors object, containing a key/value pair of string/errors, but it can be of any type of error that is
// consumable by a client.
func NewValidationError(message string, validationError error) error {
	return &ValidationError{Message: message, ValidationError: validationError}
}

// NewValidationErrorMessage creates a ValidationError without an underlying error.
func NewValidationErrorMessage(message string) error {
	return &ValidationError{Message: message}
}

func (e *ValidationError) Error() string {
	if e.ValidationError != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.ValidationError.Error())
	}

	return e.Message
}
