package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DeploymentRequest_ApplyDefaults(t *testing.T) {
	model := &DeploymentRequest{}

	err := model.ApplyDefaults()

	assert.NoError(t, err)
	// Not an exhaustive check of app defaults, just ensure that the app has a default set
	assert.Equal(t, "apps", model.App.Namespace)
}
