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

func CreateVirtualService(deployment *core.Deployment, gatewayHost string) (*VirtualService, error) {
	virtualService := &VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-%d", deployment.Namespace, deployment.Name, deployment.App.Expose.ContainerPort),
			Namespace: deployment.Namespace,
			Labels:    commonLabels(deployment),
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
				fmt.Sprintf("%s.%s.%s", deployment.Name, deployment.Namespace, gatewayHost),
				fmt.Sprintf("%s.%s.svc.cluster.local", deployment.Name, deployment.Namespace),
			},
			Http: []*istionet.HTTPRoute{
				&istionet.HTTPRoute{
					Route: []*istionet.HTTPRouteDestination{
						&istionet.HTTPRouteDestination{
							Destination: &istionet.Destination{
								Host: fmt.Sprintf("%s.%s.svc.cluster.local", deployment.Name, deployment.Namespace),
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
