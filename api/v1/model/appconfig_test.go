package model

import (
	"testing"

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
		copier.Copy(appConfig, minimumValidAppConfig)
		appConfig.Image = tt.image
		err := appConfig.Validate()

		if tt.valid {
			assert.NoError(t, err, tt.image)
		} else {
			assert.IsType(t, validation.Errors{}, err, tt.image)
			validationErrors := err.(validation.Errors)
			assert.Len(t, validationErrors, 1, tt.image)
			require.Contains(t, validationErrors, "image", tt.image)
			assert.Equal(t, "must not contain a tag or digest", validationErrors["image"].Error(), tt.image)
		}
	}
}

func assertFieldsRequired(t *testing.T, errors validation.Errors, fieldNames ...string) {
	for _, fieldName := range fieldNames {
		require.Contains(t, errors, fieldName, "missing required field %q", fieldName)
		assert.Equal(t, "cannot be blank", errors[fieldName].Error())
	}
}
