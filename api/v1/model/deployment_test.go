package model

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/jinzhu/copier"

	"github.com/stretchr/testify/assert"
)

var minimumValidDeploymentRequest = &DeploymentRequest{
	DeploymentMeta: DeploymentMeta{
		Name:        "mydep",
		Environment: "test",
		Docker:      DeploymentDocker{},
	},
	App: &AppConfigWithOverrides{AppConfig: *minimumValidAppConfig},
}

func Test_DeploymentRequest_ApplyDefaults(t *testing.T) {
	model := &DeploymentRequest{}

	err := model.ApplyDefaults()

	assert.NoError(t, err)
	// Not an exhaustive check of app defaults, just ensure that the app has a default set
	assert.EqualValues(t, "apps", model.App.Namespace)
}

func Test_DeploymentRequest_Validates(t *testing.T) {
	err := minimumValidDeploymentRequest.Validate()

	assert.NoError(t, err)
}

// This is just a sanity check that validation is happening. See Test_RulesNamingIdentifier for better coverage.
func Test_DeploymentRequest_ValidateName(t *testing.T) {
	model := createMinDeploymentRequest()
	model.Name = "5name"

	err := model.Validate()
	assert.IsType(t, validation.Errors{}, err)
	validationErrors := err.(validation.Errors)
	assert.Len(t, validationErrors, 1)
	assert.Equal(t, "must be lowercase, alphanumeric, and start with a letter", validationErrors["name"].Error())
}

func Test_DeploymentRequest_ValidateRequired(t *testing.T) {
	model := &DeploymentRequest{}
	err := model.Validate()

	assert.IsType(t, validation.Errors{}, err)
	validationErrors := err.(validation.Errors)
	assert.Len(t, validationErrors, 3)
	assertFieldsRequired(t, validationErrors, "app", "name", "environment")
}

func Test_DeploymentRequest_ValidateEmptyApp(t *testing.T) {
	model := &DeploymentRequest{
		App: &AppConfigWithOverrides{},
	}

	err := model.Validate()

	// No specifics, just ensure that an empty app is validated (validation coverage is in appconfig_test)
	assert.IsType(t, validation.Errors{}, err)
}

func createMinDeploymentRequest() *DeploymentRequest {
	model := &DeploymentRequest{}
	_ = copier.Copy(model, minimumValidDeploymentRequest)
	return model
}
