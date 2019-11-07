package resources

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"

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
func CreateDeployment(ctx *core.DeploymentContext) (*appsv1.Deployment, error) {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ctx.Deployment.Name,
			Namespace:   ctx.Deployment.Namespace,
			Labels:      commonLabels(ctx),
			Annotations: commonAnnotations(ctx),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ctx.Deployment.App.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					riserLabel("deployment"): ctx.Deployment.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: commonLabels(ctx),
					Annotations: map[string]string{
						"sidecar.istio.io/rewriteAppHTTPProbers": "true",
					},
				},
				Spec: corev1.PodSpec{
					EnableServiceLinks: boolPtr(false),
					Containers: []corev1.Container{
						corev1.Container{
							Name:           ctx.Deployment.Name,
							Image:          fmt.Sprintf("%s:%s", ctx.Deployment.App.Image, ctx.Deployment.Docker.Tag),
							Resources:      resources(ctx.Deployment.App),
							ReadinessProbe: readinessProbe(ctx.Deployment.App),
							Env:            k8sEnvVars(ctx),
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: ctx.Deployment.App.Expose.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

func readinessProbe(appConfig *model.AppConfig) *corev1.Probe {
	if appConfig.HealthCheck == nil {
		return nil
	}

	port := appConfig.HealthCheck.Port
	if port == nil && appConfig.Expose != nil {
		port = &appConfig.Expose.ContainerPort
	}

	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: appConfig.HealthCheck.Path,
				Port: intstr.FromInt(int(*port)),
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
