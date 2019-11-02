package resources

import (
	"github.com/riser-platform/riser-server/pkg/core"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateService(deployment *core.Deployment) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
			Labels:    commonLabels(deployment),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Port: deployment.App.Expose.ContainerPort,
					Name: deployment.App.Expose.Protocol,
				},
			},
			Selector: map[string]string{riserLabel("deployment"): deployment.Name},
		},
	}

	return service, nil
}
