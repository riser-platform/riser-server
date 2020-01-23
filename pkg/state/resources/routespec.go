package resources

import (
	"fmt"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateKNativeRoute(ctx *core.DeploymentContext) *Route {
	return &Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ctx.Deployment.Name,
			Namespace:   ctx.Deployment.Namespace,
			Labels:      deploymentLabels(ctx),
			Annotations: deploymentAnnotations(ctx),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: "serving.knative.dev/v1",
		},
		Spec: createRouteSpec(ctx.Deployment.Traffic),
	}
}

func createRouteSpec(trafficConfig core.TrafficConfig) RouteSpec {
	spec := RouteSpec{
		Traffic: []TrafficTarget{},
	}

	for _, rule := range trafficConfig {
		spec.Traffic = append(spec.Traffic, TrafficTarget{
			RevisionName: rule.RevisionName,
			Percent:      util.PtrInt64(int64(rule.Percent)),
			Tag:          fmt.Sprintf("r%d", rule.RiserGeneration),
		})
	}

	return spec
}
