package resources

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// TODO: Configure these defaults per stage
const (
	DefaultResourceCPUCores            = float32(0.5)
	DefaultResourceMemoryMB            = int32(256)
	DefaultResourceRequestFactorCPU    = 0.25
	DefaultResourceRequestFactorMemory = 0.5
)

// CreateDeployment creates a kubernetes Deployment from a riser deployment
func CreateDeployment(deployment *core.Deployment, secretsForEnv []string) (*appsv1.Deployment, error) {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
			Labels:    commonLabels(deployment),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: deployment.App.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					riserLabel("deployment"): deployment.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: createPodObjectMeta(deployment, commonLabels(deployment)),
				Spec:       createPodSpec(deployment, secretsForEnv),
			},
		},
	}, nil
}

func createPodObjectMeta(deployment *core.Deployment, labels map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Labels: labels,
		Annotations: map[string]string{
			"sidecar.istio.io/rewriteAppHTTPProbers": "true",
		},
	}
}

func createPodSpec(deployment *core.Deployment, secretsForEnv []string) corev1.PodSpec {
	return corev1.PodSpec{
		EnableServiceLinks: boolPtr(false),
		Containers: []corev1.Container{
			corev1.Container{
				Name:           deployment.Name,
				Image:          fmt.Sprintf("%s:%s", deployment.App.Image, deployment.Docker.Tag),
				Resources:      resources(deployment.App),
				ReadinessProbe: readinessProbe(deployment.App),
				Env:            k8sEnvVars(deployment, secretsForEnv),
				Ports: []corev1.ContainerPort{
					corev1.ContainerPort{
						Protocol:      corev1.ProtocolTCP,
						ContainerPort: deployment.App.Expose.ContainerPort,
					},
				},
			},
		},
	}
}

func readinessProbe(appConfig *model.AppConfig) *corev1.Probe {
	if appConfig.HealthCheck == nil {
		return nil
	}

	// TODO: Do not commit (KNative does not allow setting this)
	// port := appConfig.HealthCheck.Port
	// if port == nil && appConfig.Expose != nil {
	// 	port = &appConfig.Expose.ContainerPort
	// }

	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: appConfig.HealthCheck.Path,
				// TODO: Do not commit (KNative does not allow setting this)
				// Port: intstr.FromInt(int(*port)),
			},
		},
	}
}

func resources(appConfig *model.AppConfig) corev1.ResourceRequirements {
	cpuCores := DefaultResourceCPUCores
	memoryMB := DefaultResourceMemoryMB
	if appConfig.Resources != nil {
		if appConfig.Resources.CpuCores != nil {
			cpuCores = *appConfig.Resources.CpuCores
		}
		if appConfig.Resources.MemoryMB != nil {
			memoryMB = *appConfig.Resources.MemoryMB
		}
	}
	return corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    *resource.NewScaledQuantity(int64(cpuCores*float32(1000)), resource.Milli),
			corev1.ResourceMemory: *resource.NewScaledQuantity(int64(memoryMB), resource.Mega),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    *resource.NewScaledQuantity(int64(cpuCores*float32(1000)*DefaultResourceRequestFactorCPU), resource.Milli),
			corev1.ResourceMemory: *resource.NewScaledQuantity(int64(float32(memoryMB)*DefaultResourceRequestFactorMemory), resource.Mega),
		},
	}
}
