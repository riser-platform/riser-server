package resources

import (
	"fmt"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

func CreateKNativeRoute(ctx *core.DeploymentContext) *servingv1.Route {
	labels := deploymentLabels(ctx)
	if ctx.DeploymentConfig.App.Expose != nil && ctx.DeploymentConfig.App.Expose.Scope != model.AppExposeScope_External {
		labels["serving.knative.dev/visibility"] = "cluster-local"
	}
	return &servingv1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ctx.DeploymentConfig.Name,
			Namespace:   ctx.DeploymentConfig.Namespace,
			Labels:      labels,
			Annotations: deploymentAnnotations(ctx),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: "serving.knative.dev/v1",
		},
		Spec: createRouteSpec(ctx.DeploymentConfig.Traffic),
	}
}

func createRouteSpec(trafficConfig core.TrafficConfig) servingv1.RouteSpec {
	spec := servingv1.RouteSpec{
		Traffic: []servingv1.TrafficTarget{},
	}

	for _, rule := range trafficConfig {
		spec.Traffic = append(spec.Traffic, servingv1.TrafficTarget{
			RevisionName: rule.RevisionName,
			Percent:      util.PtrInt64(int64(rule.Percent)),
			Tag:          fmt.Sprintf("r%d", rule.RiserRevision),
		})
	}

	return spec
}
