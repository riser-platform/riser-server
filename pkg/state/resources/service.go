package resources

import (
	"github.com/riser-platform/riser-server/pkg/core"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateService(ctx *core.DeploymentContext) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ctx.Deployment.Name,
			Namespace:   ctx.Deployment.Namespace,
			Annotations: commonAnnotations(ctx),
			Labels:      commonLabels(ctx),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Port: ctx.Deployment.App.Expose.ContainerPort,
					Name: ctx.Deployment.App.Expose.Protocol,
				},
			},
			Selector: map[string]string{riserLabel("deployment"): ctx.Deployment.Name},
		},
	}

	return service, nil
}
