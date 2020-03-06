package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NamespacedName_String(t *testing.T) {
	name := NamespacedName{"naMe", "Ns"}

	assert.Equal(t, "name.ns", name.String())
}

func Test_NewNamespacedName_Default(t *testing.T) {
	name := NewNamespacedName("naMe", "")

	assert.Equal(t, "name", name.Name)
	assert.Equal(t, DefaultNamespace, name.Namespace)
}

func Test_ParseNamespacedName(t *testing.T) {
	result := ParseNamespacedName("myDep.mYns")

	assert.Equal(t, "mydep", result.Name)
	assert.Equal(t, "myns", result.Namespace)
}

func Test_ParseNamespacedName_NoNS(t *testing.T) {
	result := ParseNamespacedName("mydeP")

	assert.Equal(t, "mydep", result.Name)
	assert.Empty(t, result.Namespace)
}
