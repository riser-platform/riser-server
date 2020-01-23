package core

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/stretchr/testify/assert"
)

func Test_NewValidationError(t *testing.T) {
	err := &validation.Errors{}

	result := NewValidationError("msg", err)

	assert.IsType(t, &ValidationError{}, result)
	validationError := result.(*ValidationError)
	assert.Equal(t, "msg", validationError.Message)
	assert.Equal(t, err, validationError.ValidationError)
}
