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

// CreateDeployment creates a kubernetes Deployment from a riser deployment
func CreateDeployment(ctx *core.DeploymentContext) (*appsv1.Deployment, error) {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ctx.Deployment.Name,
			Namespace:   ctx.Deployment.Namespace,
			Labels:      deploymentLabels(ctx),
			Annotations: deploymentAnnotations(ctx),
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
				ObjectMeta: createPodObjectMeta(ctx),
				Spec:       createPodSpec(ctx),
			},
		},
	}, nil
}

func createPodObjectMeta(ctx *core.DeploymentContext) metav1.ObjectMeta {
	annotations := deploymentAnnotations(ctx)
	// This is required for probes to work w/mTLS on teh same port as the service port
	annotations["sidecar.istio.io/rewriteAppHTTPProbers"] = "true"
	return metav1.ObjectMeta{
		Labels:      deploymentLabels(ctx),
		Annotations: annotations,
	}
}

func createPodSpec(ctx *core.DeploymentContext) corev1.PodSpec {
	return corev1.PodSpec{
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
	}
}

func readinessProbe(appConfig *model.AppConfig) *corev1.Probe {
	if appConfig.HealthCheck == nil {
		return nil
	}

	port := appConfig.HealthCheck.Port

	probe := &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: appConfig.HealthCheck.Path,
			},
		},
	}

	// TODO: KNative doesn't appear to allow setting this, need to investigate further
	if port != nil {
		probe.Handler.HTTPGet.Port = intstr.FromInt(int(*port))
	}

	return probe
}

func resources(appConfig *model.AppConfig) corev1.ResourceRequirements {
	res := corev1.ResourceRequirements{}
	if appConfig.Resources != nil {
		res.Limits = corev1.ResourceList{}
		if appConfig.Resources.CpuCores != nil {
			res.Limits[corev1.ResourceCPU] = *resource.NewScaledQuantity(int64(*appConfig.Resources.CpuCores*float32(1000)), resource.Milli)
		}
		if appConfig.Resources.MemoryMB != nil {
			res.Limits[corev1.ResourceMemory] = *resource.NewScaledQuantity(int64(*appConfig.Resources.MemoryMB), resource.Mega)
		}
	}
	return res
}
