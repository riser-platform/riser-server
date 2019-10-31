package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	validation "github.com/go-ozzo/ozzo-validation"
)

var appNameTests = []struct {
	app      string
	expected string
}{
	{"good-app", ""},
	{"a", "the length must be between 3 and 50"},
	{"123456789012345678901234567890123456789012345678901", "the length must be between 3 and 50"},
	{"1abc", "must be lowercase, alphanumeric, and start with a letter"},
	{"ABCD", "must be lowercase, alphanumeric, and start with a letter"},
	{"bad!", "must be lowercase, alphanumeric, and start with a letter"},
	{"", "cannot be blank"},
}

func Test_RulesAppName(t *testing.T) {
	for _, tt := range appNameTests {
		err := validation.Validate(tt.app, RulesAppName()...)
		assertValidationMessage(t, tt.expected, err)
	}
}

func assertValidationMessage(t *testing.T, expected string, err error) {
	if expected == "" {
		assert.NoError(t, err)
	} else {
		assert.Equal(t, expected, err.Error())
	}

}
