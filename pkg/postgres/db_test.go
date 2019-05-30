package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddAuthToConnString(t *testing.T) {
	result, err := AddAuthToConnString(
		"postgres://myhost.local/riserdb?arg=val",
		"myuser",
		"mypass")

	assert.NoError(t, err)
	assert.Equal(t, "postgres://myuser:mypass@myhost.local/riserdb?arg=val", result)
}

func Test_AddAuthToConnString_BadUrl(t *testing.T) {
	result, err := AddAuthToConnString(
		"not@valid:test",
		"myuser",
		"mypass")

	assert.Empty(t, result)
	assert.Error(t, err)
}
