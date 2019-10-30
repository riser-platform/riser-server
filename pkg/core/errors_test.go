package core

import (
	"testing"

	"github.com/pkg/errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/stretchr/testify/assert"
)

func Test_NewValidationError_WhenValidationErrors(t *testing.T) {
	err := validation.Errors{}

	result := NewValidationError("msg", err)

	assert.IsType(t, &ValidationError{}, result)
	validationError := result.(*ValidationError)
	assert.Equal(t, "msg", validationError.Message)
	assert.Equal(t, err, validationError.ValidationErrors)
	assert.Nil(t, validationError.Internal)
}

func Test_NewValidationError_WhenNotValidationError(t *testing.T) {
	err := errors.New("int")

	result := NewValidationError("msg", err)

	assert.IsType(t, &ValidationError{}, result)
	validationError := result.(*ValidationError)
	assert.Equal(t, "msg", validationError.Message)
	assert.Equal(t, err, validationError.Internal)
	assert.Nil(t, validationError.ValidationErrors)
}
