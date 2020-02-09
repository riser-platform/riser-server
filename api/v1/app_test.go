package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/pkg/core"
)

func Test_PostApp_Validates(t *testing.T) {
	newApp := model.NewApp{
		Name: "0notvalid",
	}

	req := httptest.NewRequest(http.MethodPost, "/", safeMarshal(newApp))
	req.Header.Add("CONTENT-TYPE", "application/json")
	ctx, _ := newContextWithRecorder(req)

	result := PostApp(ctx, nil)

	assert.Equal(t, "invalid app name: must be lowercase, alphanumeric, and start with a letter", result.Error())
}

func Test_mapAppFromDomain(t *testing.T) {
	appId, _ := core.DecodeAppId("aaaa")
	domain := core.App{
		Name:   "myapp",
		Hashid: appId,
	}

	result := mapAppFromDomain(domain)

	assert.Equal(t, "myapp", result.Name)
	assert.Equal(t, "aaaa", result.Id)
}

func Test_mapAppArrayFromDomain(t *testing.T) {
	domainArray := []core.App{
		core.App{
			Name: "myapp1",
		},
		core.App{
			Name: "myapp2",
		},
	}

	result := mapAppArrayFromDomain(domainArray)

	assert.Len(t, result, 2)
	assert.Equal(t, "myapp1", result[0].Name)
	assert.Equal(t, "myapp2", result[1].Name)
}
