package resources

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/util"

	corev1 "k8s.io/api/core/v1"
)

func createPodSpec(ctx *core.DeploymentContext) corev1.PodSpec {
	return corev1.PodSpec{
		EnableServiceLinks: util.PtrBool(false),
		Containers: []corev1.Container{
			{
				Name:           ctx.DeploymentConfig.Name,
				Image:          fmt.Sprintf("%s:%s", ctx.DeploymentConfig.App.Image, ctx.DeploymentConfig.Docker.Tag),
				Resources:      resources(ctx.DeploymentConfig.App),
				ReadinessProbe: readinessProbe(ctx.DeploymentConfig.App),
				Env:            k8sEnvVars(ctx),
				Ports:          createPodPorts(ctx.DeploymentConfig.App.Expose),
			},
		},
	}
}

func createPodPorts(expose *model.AppConfigExpose) []corev1.ContainerPort {
	containerPortName := ""
	// See https://github.com/knative/serving/blob/master/docs/runtime-contract.md#protocols-and-ports
	if expose.Protocol == "http2" {
		containerPortName = "h2c"
	}
	ports := []corev1.ContainerPort{
		{
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: expose.ContainerPort,
			Name:          containerPortName,
		},
	}
	return ports
}

func readinessProbe(appConfig *model.AppConfig) *corev1.Probe {
	if appConfig.HealthCheck == nil {
		return nil
	}

	probe := &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: appConfig.HealthCheck.Path,
			},
		},
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
