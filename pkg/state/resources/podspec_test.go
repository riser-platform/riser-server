package resources

import (
	"testing"

	"github.com/riser-platform/riser-server/pkg/util"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/stretchr/testify/assert"
)

// Basic podspec tests are covered in knativeservice_test.go.

func Test_readinessProbe_nilDeploy(t *testing.T) {
	app := &model.AppConfig{}

	result := readinessProbe(app)

	assert.Nil(t, result)
}

func Test_readinessProbe_nilHealth(t *testing.T) {
	app := &model.AppConfig{}

	result := readinessProbe(app)

	assert.Nil(t, result)
}

func Test_readinessProbe_httpGet(t *testing.T) {
	app := &model.AppConfig{
		HealthCheck: &model.AppConfigHealthCheck{
			Path: "/health",
			Port: util.PtrInt32(8080),
		},
	}

	result := readinessProbe(app)

	assert.Equal(t, "/health", result.HTTPGet.Path)
	assert.EqualValues(t, 8080, result.HTTPGet.Port.IntVal)
}

func Test_resources(t *testing.T) {
	app := &model.AppConfig{
		Resources: &model.AppConfigResources{
			CpuCores: util.PtrFloat32(1.5),
			MemoryMB: util.PtrInt32(4096),
		},
	}

	result := resources(app)

	assert.EqualValues(t, 1500, result.Limits.Cpu().MilliValue(), "millicores")
	assert.EqualValues(t, 4096000000, result.Limits.Memory().Value(), "bytes")
}
