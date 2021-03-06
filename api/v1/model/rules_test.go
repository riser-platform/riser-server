package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	validation "github.com/go-ozzo/ozzo-validation/v3"
)

func Test_RulesAppName(t *testing.T) {
	var tests = []struct {
		app      string
		expected string
	}{
		{"good-app", ""},
		{"a", "the length must be between 3 and 47"},
		{"123456789012345678901234567890123456789012345678901", "the length must be between 3 and 47"},
		{"1abc", "must be lowercase, alphanumeric, and start with a letter"},
		{"ABCD", "must be lowercase, alphanumeric, and start with a letter"},
		{"bad!", "must be lowercase, alphanumeric, and start with a letter"},
		{"", "cannot be blank"},
	}

	for _, tt := range tests {
		err := validation.Validate(tt.app, RulesAppName()...)
		assertValidationMessage(t, tt.app, tt.expected, err)
	}
}

func Test_RulesNamingIdentifier(t *testing.T) {
	var tests = []struct {
		v string
		e string
	}{
		// Good
		{"valid1", ""},
		{"is-valid", ""},
		// Bad
		{"va", "the length must be between 3 and 63"},
		{"1234567890123456789012345678901234567890123456789012345678901234", "the length must be between 3 and 63"},
		{"1abc", "must be lowercase, alphanumeric, and start with a letter"},
		{"ABCD", "must be lowercase, alphanumeric, and start with a letter"},
		{"bad!", "must be lowercase, alphanumeric, and start with a letter"},
		{"bad-", "must be lowercase, alphanumeric, and start with a letter"},
	}

	for _, tt := range tests {
		err := validation.Validate(tt.v, RulesNamingIdentifier()...)
		assertValidationMessage(t, tt.v, tt.e, err)
	}
}

func assertValidationMessage(t *testing.T, value string, expected string, err error) {
	if expected == "" {
		assert.NoError(t, err, value)
	} else {
		require.Error(t, err, expected, value)
		assert.Equal(t, expected, err.Error(), value)
	}

}
