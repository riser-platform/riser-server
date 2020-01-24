package resources

import (
	"fmt"

	"github.com/riser-platform/riser-server/pkg/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateKNativeConfiguration(ctx *core.DeploymentContext) *Configuration {
	podSpec := createPodSpec(ctx)
	// KNative does not allow setting this
	podSpec.EnableServiceLinks = nil

	revisionMeta := createRevisionMeta(ctx)

	// Not sure yet if we want this with KNative since KNative seems to handle readiness probes differently via the queue-proxy.

	return &Configuration{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ctx.Deployment.Name,
			Namespace:   ctx.Deployment.Namespace,
			Labels:      deploymentLabels(ctx),
			Annotations: deploymentAnnotations(ctx),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Configuration",
			APIVersion: "serving.knative.dev/v1",
		},
		Spec: ConfigurationSpec{
			Template: RevisionTemplateSpec{
				ObjectMeta: revisionMeta,
				Spec: RevisionSpec{
					PodSpec: podSpec,
				},
			},
		},
	}
}

func createRevisionMeta(ctx *core.DeploymentContext) metav1.ObjectMeta {
	revisionMeta := metav1.ObjectMeta{
		Name:        fmt.Sprintf("%s-%d", ctx.Deployment.Name, ctx.RiserRevision),
		Labels:      deploymentLabels(ctx),
		Annotations: deploymentAnnotations(ctx),
	}
	if ctx.Deployment.App.Autoscale != nil {
		if ctx.Deployment.App.Autoscale.Min != nil {
			revisionMeta.Annotations["autoscaling.knative.dev/minScale"] = fmt.Sprintf("%d", *ctx.Deployment.App.Autoscale.Min)
		}
		if ctx.Deployment.App.Autoscale.Max != nil {
			revisionMeta.Annotations["autoscaling.knative.dev/maxScale"] = fmt.Sprintf("%d", *ctx.Deployment.App.Autoscale.Max)
		}
	}

	return revisionMeta
}
