package resources

import (
	"github.com/riser-platform/riser-server/pkg/core"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNamespace(namespace *core.Namespace) (*corev1.Namespace, error) {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace.Name,
			Labels: map[string]string{
				"istio-injection": "enabled",
			},
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
	}, nil
}
