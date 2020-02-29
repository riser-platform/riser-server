package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NamespacedName_String(t *testing.T) {
	name := NamespacedName{"name", "ns"}

	assert.Equal(t, "name.ns", name.String())
}

func Test_ParseNamespacedName(t *testing.T) {
	result := ParseNamespacedName("mydep.myns")

	assert.Equal(t, "mydep", result.Name)
	assert.Equal(t, "myns", result.Namespace)
}

func Test_ParseNamespacedName_NoNS(t *testing.T) {
	result := ParseNamespacedName("mydep")

	assert.Equal(t, "mydep", result.Name)
	assert.Empty(t, result.Namespace)
}
