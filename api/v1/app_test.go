package v1

import (
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/pkg/core"
)

func Test_mapAppFromDomain(t *testing.T) {
	domain := core.App{
		Id:        uuid.New(),
		Name:      "myapp",
		Namespace: "myns",
	}

	result := mapAppFromDomain(domain)

	assert.Equal(t, domain.Id, result.Id)
	assert.EqualValues(t, "myapp", result.Name)
	assert.EqualValues(t, "myns", result.Namespace)
}

func Test_mapAppArrayFromDomain(t *testing.T) {
	domainArray := []core.App{
		core.App{
			Id:        uuid.New(),
			Name:      "myapp1",
			Namespace: "myns1",
		},
		core.App{
			Name: "myapp2",
		},
	}

	result := mapAppArrayFromDomain(domainArray)

	assert.Len(t, result, 2)
	assert.Equal(t, domainArray[0].Id, result[0].Id)
	assert.EqualValues(t, "myapp1", result[0].Name)
	assert.EqualValues(t, "myns1", result[0].Namespace)
	assert.EqualValues(t, "myapp2", result[1].Name)
}
