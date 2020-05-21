package v1

import (
	"testing"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_mapNamespaceFromDomain(t *testing.T) {
	domain := core.Namespace{Name: "myns"}

	result := mapNamespaceFromDomain(domain)

	assert.EqualValues(t, "myns", result.Name)
}

func Test_mapNamespaceArrayFromDomain(t *testing.T) {
	domainArray := []core.Namespace{
		{Name: "myns1"},
		{Name: "myns2"},
	}

	result := mapNamespaceArrayFromDomain(domainArray)

	assert.Len(t, result, 2)
	assert.EqualValues(t, "myns1", result[0].Name)
	assert.EqualValues(t, "myns2", result[1].Name)
}
