package resources

import (
	"github.com/riser-platform/riser-server/pkg/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateKNativeService(deployment *core.Deployment, secretsForEnv []string) *Service {
	labels := commonLabels(deployment)
	// TODO: This won't be necessary after we move to using riser.dev/app.
	delete(labels, "app")

	podSpec := createPodSpec(deployment, secretsForEnv)
	// KNative does not allow setting this
	podSpec.EnableServiceLinks = nil

	podMeta := createPodObjectMeta(deployment, labels)
	// We should consider exposing this in the app config. We don't want to disable scale-to-zero cluster wide as we
	// want to eventually support that on an app by app basis.
	podMeta.Annotations["autoscaling.knative.dev/minScale"] = "1"
	// Not sure yet if we want this with KNative since KNative seems to handle readiness probes differently via the queue-proxy.
	// TODO: Test KNative w/ Istio mTLS to see if we still need this attribute.
	delete(podMeta.Annotations, "sidecar.istio.io/rewriteAppHTTPProbers")

	return &Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
			Labels:    labels,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "serving.knative.dev/v1",
		},
		ServiceSpec: ServiceSpec{
			ConfigurationSpec: ConfigurationSpec{
				Template: RevisionTemplateSpec{
					ObjectMeta: podMeta,
					Spec: RevisionSpec{
						PodSpec: podSpec,
					},
				},
			},
		},
	}
}
