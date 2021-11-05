package resources

import (
	"fmt"

	"github.com/riser-platform/riser-server/pkg/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

// Create a knative serving spec pod spec

func CreateKNativeConfiguration(ctx *core.DeploymentContext) *servingv1.Configuration {
	podSpec := createPodSpec(ctx)
	// KNative does not allow setting this
	podSpec.EnableServiceLinks = nil

	revisionMeta := createRevisionMeta(ctx)

	// Not sure yet if we want this with KNative since KNative seems to handle readiness probes differently via the queue-proxy.

	return &servingv1.Configuration{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ctx.DeploymentConfig.Name,
			Namespace:   ctx.DeploymentConfig.Namespace,
			Labels:      deploymentLabels(ctx),
			Annotations: deploymentAnnotations(ctx),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Configuration",
			APIVersion: "serving.knative.dev/v1",
		},
		Spec: servingv1.ConfigurationSpec{
			Template: servingv1.RevisionTemplateSpec{
				ObjectMeta: revisionMeta,
				Spec: servingv1.RevisionSpec{
					PodSpec: podSpec,
				},
			},
		},
	}
}

func createRevisionMeta(ctx *core.DeploymentContext) metav1.ObjectMeta {
	revisionMeta := metav1.ObjectMeta{
		Name:        fmt.Sprintf("%s-%d", ctx.DeploymentConfig.Name, ctx.RiserRevision),
		Labels:      deploymentLabels(ctx),
		Annotations: deploymentAnnotations(ctx),
	}
	if ctx.DeploymentConfig.App.Autoscale != nil {
		if ctx.DeploymentConfig.App.Autoscale.Min != nil {
			revisionMeta.Annotations["autoscaling.knative.dev/minScale"] = fmt.Sprintf("%d", *ctx.DeploymentConfig.App.Autoscale.Min)
		}
		if ctx.DeploymentConfig.App.Autoscale.Max != nil {
			revisionMeta.Annotations["autoscaling.knative.dev/maxScale"] = fmt.Sprintf("%d", *ctx.DeploymentConfig.App.Autoscale.Max)
		}
	}

	return revisionMeta
}
