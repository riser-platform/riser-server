package resources

import (
	"testing"

	"github.com/riser-platform/riser-server/pkg/util"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
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
		},
	}

	result := readinessProbe(app)

	assert.Equal(t, "/health", result.HTTPGet.Path)
	// KNative does not allow setting the port on a probe
	assert.Empty(t, result.HTTPGet.Port)
}

func Test_resources(t *testing.T) {
	app := &model.AppConfig{
		OverrideableAppConfig: model.OverrideableAppConfig{
			Resources: &model.AppConfigResources{
				CpuCores: util.PtrFloat32(1.5),
				MemoryMB: util.PtrInt32(4096),
			},
		},
	}

	result := resources(app)

	assert.EqualValues(t, 1500, result.Limits.Cpu().MilliValue(), "millicores")
	assert.EqualValues(t, 4096000000, result.Limits.Memory().Value(), "bytes")
}

func Test_createPodPorts_http(t *testing.T) {
	expose := &model.AppConfigExpose{
		Protocol:      "http",
		ContainerPort: 80,
	}

	result := createPodPorts(expose)

	assert.Len(t, result, 1)
	assert.EqualValues(t, 80, result[0].ContainerPort)
	assert.Equal(t, corev1.ProtocolTCP, result[0].Protocol)
	assert.Empty(t, result[0].Name)
}

func Test_createPodPorts_http2(t *testing.T) {
	expose := &model.AppConfigExpose{
		Protocol:      "http2",
		ContainerPort: 80,
	}

	result := createPodPorts(expose)

	assert.Len(t, result, 1)
	assert.EqualValues(t, 80, result[0].ContainerPort)
	assert.Equal(t, corev1.ProtocolTCP, result[0].Protocol)
	assert.Equal(t, "h2c", result[0].Name)
}
