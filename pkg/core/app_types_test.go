package core

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_AppIdentifier_IdValue(t *testing.T) {
	validUUID := uuid.New()
	tt := []struct {
		v                AppIdOrName
		expectedValue    uuid.UUID
		expectedHasValue bool
	}{
		{AppIdOrName(validUUID.String()), validUUID, true},
		{AppIdOrName("foo"), uuid.Nil, false},
	}

	for _, test := range tt {
		value, hasValue := test.v.IdValue()
		assert.Equal(t, test.expectedValue, value, test.v)
		assert.Equal(t, test.expectedHasValue, hasValue, test.v)
	}
}

func Test_AppIdentifier_NameValue(t *testing.T) {
	tt := []struct {
		v                AppIdOrName
		expectedValue    *NamespacedName
		expectedHasValue bool
	}{
		{AppIdOrName("foo"), nil, false},
		{AppIdOrName(uuid.New().String()), nil, false},
		{AppIdOrName("myapp.myns"), &NamespacedName{"myapp", "myns"}, true},
	}
	for _, test := range tt {
		value, hasValue := test.v.NameValue()
		assert.Equal(t, test.expectedValue, value, test.v)
		assert.Equal(t, test.expectedHasValue, hasValue, test.v)
	}
}
