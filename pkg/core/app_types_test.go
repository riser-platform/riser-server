package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DecodeAppId(t *testing.T) {
	result, err := DecodeAppId("ff")

	assert.NoError(t, err)
	assert.Equal(t, []byte{255}, []byte(result))
}

func Test_DecodeAppId_WhenInvalid_ReturnsErr(t *testing.T) {
	result, err := DecodeAppId("Z")

	assert.Nil(t, result)
	assert.Equal(t, err.Error(), "app ID must be a valid hex string")

}
