package resources

import (
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func Test_CreateDeployment(t *testing.T) {
	replicas := int32(2)
	deployment := &core.DeploymentConfig{
		Name:      "myapp-deployment",
		Namespace: "apps",
		Stage:     "dev",
		Docker: core.DeploymentDocker{
			Tag: "myTag",
		},
		App: &model.AppConfig{
			Name:     "myapp",
			Image:    "hashicorp/http-echo",
			Replicas: &replicas,
			HealthCheck: &model.AppConfigHealthCheck{
				Path: "/health",
			},
			Environment: map[string]intstr.IntOrString{
				"envKey": intstr.Parse("envVal"),
			},
			Expose: &model.AppConfigExpose{
				ContainerPort: 1337,
			},
		},
	}

	secretsForEnv := []string{
		"mysecret",
	}

	result, err := CreateDeployment(&core.DeploymentContext{Deployment: deployment, SecretNames: secretsForEnv, RiserGeneration: 3})

	assert.Nil(t, err)
	assert.NotNil(t, result)
	// Metadata
	assert.Equal(t, "myapp-deployment", result.Name)
	assert.Equal(t, "apps", result.Namespace)
	assert.Equal(t, "Deployment", result.Kind)
	assert.Equal(t, "apps/v1", result.APIVersion)
	assert.Equal(t, 3, len(result.Labels))
	assert.Equal(t, "dev", result.Labels[riserLabel("stage")])
	assert.Equal(t, "myapp", result.Labels[riserLabel("app")])
	assert.Equal(t, "myapp-deployment", result.Labels[riserLabel("deployment")])
	assert.Equal(t, 1, len(result.Annotations))
	assert.Equal(t, "3", result.Annotations[riserLabel("generation")])

	// Pod
	assert.Len(t, result.Spec.Template.Labels, 3)
	assert.Equal(t, "dev", result.Spec.Template.Labels[riserLabel("stage")])
	assert.Equal(t, "myapp", result.Spec.Template.Labels[riserLabel("app")])
	assert.Equal(t, "myapp-deployment", result.Spec.Template.Labels[riserLabel("deployment")])
	assert.Len(t, result.Spec.Template.Annotations, 2)
	assert.Equal(t, "3", result.Spec.Template.Annotations[riserLabel("generation")])
	assert.Equal(t, result.Spec.Template.Annotations["sidecar.istio.io/rewriteAppHTTPProbers"], "true")
	assert.Equal(t, 1, len(result.Spec.Template.Spec.Containers))
	assert.Equal(t, &replicas, result.Spec.Replicas)
	assert.False(t, *result.Spec.Template.Spec.EnableServiceLinks)

	container := result.Spec.Template.Spec.Containers[0]
	assert.Equal(t, "hashicorp/http-echo:myTag", container.Image)

	// Env
	assert.Equal(t, 2, len(container.Env))
	assert.Equal(t, "ENVKEY", container.Env[0].Name)
	assert.Equal(t, "envVal", container.Env[0].Value)
	assert.Equal(t, "MYSECRET", container.Env[1].Name)
	assert.Equal(t, "myapp-mysecret", container.Env[1].ValueFrom.SecretKeyRef.LocalObjectReference.Name)
	assert.Equal(t, "data", container.Env[1].ValueFrom.SecretKeyRef.Key)

	// Ports
	assert.Equal(t, 1, len(container.Ports))
	assert.EqualValues(t, 1337, container.Ports[0].ContainerPort)
	assert.Equal(t, corev1.ProtocolTCP, container.Ports[0].Protocol)

	// Health (Readiness Probe)
	assert.Equal(t, "/health", container.ReadinessProbe.HTTPGet.Path)
	assert.Empty(t, container.ReadinessProbe.HTTPGet.Port)

	// Resource Defaults
	assert.EqualValues(t, 500, container.Resources.Limits.Cpu().MilliValue(), "millicores")
	assert.EqualValues(t, 256000000, container.Resources.Limits.Memory().Value(), "bytes")
	assert.EqualValues(t, 125, container.Resources.Requests.Cpu().MilliValue(), "millicores")
	assert.EqualValues(t, 128000000, container.Resources.Requests.Memory().Value(), "bytes")
}

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
			Port: int32Ptr(8080),
		},
	}

	result := readinessProbe(app)

	assert.Equal(t, "/health", result.HTTPGet.Path)
	assert.EqualValues(t, 8080, result.HTTPGet.Port.IntVal)
}

func Test_resources(t *testing.T) {
	app := &model.AppConfig{
		Resources: &model.AppConfigResources{
			CpuCores: float32Ptr(1.5),
			MemoryMB: int32Ptr(4096),
		},
	}

	result := resources(app)

	assert.EqualValues(t, 1500, result.Limits.Cpu().MilliValue(), "millicores")
	assert.EqualValues(t, 4096000000, result.Limits.Memory().Value(), "bytes")
}
