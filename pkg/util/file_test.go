package util

import (
	"testing"

	"gotest.tools/assert"
)

func Test_EnsureTrailingSlash(t *testing.T) {
	tt := []struct {
		in       string
		expected string
	}{
		{"/test", "/test/"},
		{"/test/", "/test/"},
	}

	for _, test := range tt {
		result := EnsureTrailingSlash(test.in)
		assert.Equal(t, test.expected, result)
	}
}
