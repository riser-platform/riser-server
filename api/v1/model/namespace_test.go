package model

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Namespace_Validate_Name(t *testing.T) {
	tt := []struct {
		name        string
		expectedErr string
	}{
		{"myns", ""},
		{"!bad", "must be lowercase, alphanumeric, and start with a letter"},
		{"riser-nope", `namespace names may not begin with "riser-"`},
		{"knative-nope", `namespace names may not begin with "knative-"`},
		{"kube-nope", `namespace names may not begin with "kube-"`},
		{"istio-nope", `namespace names may not begin with "istio-"`},
	}

	for _, test := range tt {
		namespace := &Namespace{Name: NamespaceName(test.name)}
		err := namespace.Validate()
		if test.expectedErr == "" {
			assert.NoError(t, err)
		} else {
			require.IsType(t, validation.Errors{}, err, test.name)
			validationErrors := err.(validation.Errors)
			require.Len(t, validationErrors, 1, test.name)
			assert.Equal(t, test.expectedErr, validationErrors["name"].Error(), test.name)
		}
	}
}
