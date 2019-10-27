package model

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/jinzhu/copier"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var minimumValidAppConfig = &AppConfig{
	Name:  "myapp",
	Id:    "myid",
	Image: "myimage",
}

func Test_AppConfig_ValidateRequired(t *testing.T) {
	appConfig := AppConfig{}
	err := appConfig.Validate()

	assert.IsType(t, validation.Errors{}, err)
	validationErrors := err.(validation.Errors)
	assert.Len(t, validationErrors, 3)
	assertFieldsRequired(t, validationErrors, "name", "id", "image")
}

// Note: We may not allow registry to be set here - it may be dictated by an admin on a per stage basis instead.
var imageTests = []struct {
	image string
	valid bool
}{
	{"image", true},
	{"my/image", true},
	{"registry.io/image", true},
	{"registry.io/my/image", true},
	{"image:tag", false},
	{"image@sha256:b876dd4c32a96067ab22201e521d4fe3724f6e5af7d48f50b0059ae253359a4c", false},
	{"my/image:tag", false},
	{"registry.io/my/image:tag", false},
}

func Test_AppConfig_ValidateImage(t *testing.T) {
	for _, tt := range imageTests {
		appConfig := &AppConfig{}
		_ = copier.Copy(appConfig, minimumValidAppConfig)
		appConfig.Image = tt.image
		err := appConfig.Validate()

		if tt.valid {
			assert.NoError(t, err, tt.image)
		} else {
			require.IsType(t, validation.Errors{}, err, tt.image)
			validationErrors := err.(validation.Errors)
			assert.Len(t, validationErrors, 1, tt.image)
			require.Contains(t, validationErrors, "image", tt.image)
			assert.Equal(t, "must not contain a tag or digest", validationErrors["image"].Error(), tt.image)
		}
	}
}

var protocolTests = []struct {
	protocol string
	valid    bool
}{
	{"http", true},
	{"grpc", true},
	{"", true},
	{"redis", false},
}

func Test_AppConfig_ValidateProtocol(t *testing.T) {
	for _, tt := range protocolTests {
		appConfig := &AppConfig{}
		_ = copier.Copy(appConfig, minimumValidAppConfig)
		appConfig.Expose = &AppConfigExpose{Protocol: tt.protocol}
		err := appConfig.Validate()

		if tt.valid {
			assert.NoError(t, err, tt.protocol)
		} else {
			require.IsType(t, validation.Errors{}, err, tt.protocol)
			validationErrors := err.(validation.Errors)
			assert.Len(t, validationErrors, 1, tt.protocol)
			require.Contains(t, validationErrors, "expose.protocol", tt.protocol)
			assert.Equal(t, "must be one of: http, grpc", validationErrors["expose.protocol"].Error(), tt.protocol)
		}
	}
}

func Test_mergeValidationErrors(t *testing.T) {
	errA := validation.Errors{}
	errA["field1"] = errors.New("field1 error")
	errB := validation.Errors{}
	errB["field2"] = errors.New("field2 error")

	result := mergeValidationErrors(errA, errB, "b")

	require.IsType(t, validation.Errors{}, result)
	validationErrors := result.(validation.Errors)
	assert.Len(t, validationErrors, 2)
	assert.Equal(t, "field1 error", validationErrors["field1"].Error())
	assert.Equal(t, "field2 error", validationErrors["b.field2"].Error())
}

func Test_mergeValidationErrors_BaseNotValidationError(t *testing.T) {
	errA := errors.New("internal error")
	errB := validation.Errors{}
	errB["field2"] = errors.New("field2 error")

	result := mergeValidationErrors(errA, errB, "b")

	assert.Equal(t, errA, result)
}

func Test_mergeValidationErrors_ToMergeNotValidationError(t *testing.T) {
	errA := validation.Errors{}
	errB := errors.New("internal error")

	result := mergeValidationErrors(errA, errB, "b")

	assert.Equal(t, errB, result)
}

func Test_mergeValidationErrors_NilBase(t *testing.T) {
	errB := validation.Errors{}
	errB["field2"] = errors.New("field2 error")

	result := mergeValidationErrors(nil, errB, "b")

	require.IsType(t, validation.Errors{}, result)
	validationErrors := result.(validation.Errors)
	assert.Len(t, validationErrors, 1)
	assert.Equal(t, "field2 error", validationErrors["b.field2"].Error())
}

func Test_mergeValidationErrors_NilToMerge(t *testing.T) {
	errA := validation.Errors{}
	errA["field1"] = errors.New("field1 error")

	result := mergeValidationErrors(errA, nil, "b")

	require.IsType(t, validation.Errors{}, result)
	validationErrors := result.(validation.Errors)
	assert.Len(t, validationErrors, 1)
	assert.Equal(t, "field1 error", validationErrors["field1"].Error())
}

func assertFieldsRequired(t *testing.T, errors validation.Errors, fieldNames ...string) {
	for _, fieldName := range fieldNames {
		require.Contains(t, errors, fieldName, "missing required field %q", fieldName)
		assert.Equal(t, "cannot be blank", errors[fieldName].Error())
	}
}
