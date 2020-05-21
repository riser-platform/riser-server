package model

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UnsealedSecret_ValidateRequired(t *testing.T) {
	secret := UnsealedSecret{}

	err := secret.Validate()

	require.IsType(t, validation.Errors{}, err)
	validationErrors := err.(validation.Errors)
	assert.Len(t, validationErrors, 5)
	assertFieldsRequired(t, validationErrors, "name", "plainTextValue", "app", "namespace", "environment")
}
