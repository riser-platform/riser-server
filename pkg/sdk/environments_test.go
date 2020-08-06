package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/stretchr/testify/assert"
)

func Test_Environments_Ping(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/environments/dev/ping", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.Environments.Ping("dev")

	assert.NoError(t, err)
}

func Test_Environments_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/environments", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, `[{"name":"dev"},{"name":"prod"}]`)
	})

	environments, err := client.Environments.List()

	assert.NoError(t, err)
	assert.Len(t, environments, 2)
	assert.Equal(t, "dev", environments[0].Name)
	assert.Equal(t, "prod", environments[1].Name)
}

func Test_Environments_SetConfig(t *testing.T) {
	setup()
	defer teardown()

	config := &model.EnvironmentConfig{PublicGatewayHost: "tempuri.org"}

	mux.HandleFunc("/api/v1/environments/dev/config", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		actualConfig := &model.EnvironmentConfig{}
		mustUnmarshalR(r.Body, actualConfig)
		assert.Equal(t, config, actualConfig)
		w.WriteHeader(http.StatusAccepted)
	})

	err := client.Environments.SetConfig("dev", config)

	assert.NoError(t, err)
}

func Test_Environments_GetConfig(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/environments/dev/config", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, `{"publicGatewayHost": "tempuri.org"}`)
		w.WriteHeader(http.StatusAccepted)
	})

	result, err := client.Environments.GetConfig("dev")

	assert.NoError(t, err)
	assert.Equal(t, "tempuri.org", result.PublicGatewayHost)
}
