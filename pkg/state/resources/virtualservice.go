package resources

import (
	"fmt"

	"github.com/riser-platform/riser-server/pkg/core"
	istionet "istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: Make configurable
const defaultGateway = "riser-system/public-default"

// VirtualService represents an istio VirtualService
type VirtualService struct {
	metav1.ObjectMeta `json:"metadata"`
	metav1.TypeMeta   `json:",inline"`
	Spec              *istionet.VirtualService `json:"spec"`
}

func CreateVirtualService(ctx *core.DeploymentContext) (*VirtualService, error) {
	virtualService := &VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-%s-%d", ctx.Deployment.Namespace, ctx.Deployment.Name, ctx.Deployment.App.Expose.ContainerPort),
			Namespace:   ctx.Deployment.Namespace,
			Annotations: commonAnnotations(ctx),
			Labels:      commonLabels(ctx),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "VirtualService",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		Spec: &istionet.VirtualService{
			Gateways: []string{
				defaultGateway,
				// TODO: Validate mesh requests are routing via gateway and virtualService
				"mesh",
			},
			Hosts: []string{
				fmt.Sprintf("%s.%s.%s", ctx.Deployment.Name, ctx.Deployment.Namespace, ctx.Stage.PublicGatewayHost),
				fmt.Sprintf("%s.%s.svc.cluster.local", ctx.Deployment.Name, ctx.Deployment.Namespace),
			},
			Http: []*istionet.HTTPRoute{
				&istionet.HTTPRoute{
					Route: []*istionet.HTTPRouteDestination{
						&istionet.HTTPRouteDestination{
							Destination: &istionet.Destination{
								Host: fmt.Sprintf("%s.%s.svc.cluster.local", ctx.Deployment.Name, ctx.Deployment.Namespace),
							},
						},
					},
					// Prevents downtime during a deployment (https://github.com/istio/istio/issues/13616)
					// But now breaks on Istio 1.3... what a joke
					// Retries: &istionet.HTTPRetry{
					// 	PerTryTimeout: types.DurationProto(20 * time.Millisecond),
					// 	Attempts:      3,
					// 	RetryOn:       "gateway-error,connect-failure,refused-stream",
					// },
				},
			},
		},
	}

	return virtualService, nil
}
