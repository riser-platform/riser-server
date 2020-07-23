package resources

import (
	"fmt"

	"github.com/riser-platform/riser-server/pkg/core"
	securityv1beta1 "istio.io/api/security/v1beta1"
	typev1beta1 "istio.io/api/type/v1beta1"
	"istio.io/client-go/pkg/apis/security/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateHealthcheckDenyPolicy(dCtx *core.DeploymentContext) *v1beta1.AuthorizationPolicy {
	if dCtx.DeploymentConfig.App.HealthCheck == nil {
		return nil
	}

	return &v1beta1.AuthorizationPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-healthcheck-deny", dCtx.DeploymentConfig.Name),
			Namespace:   dCtx.DeploymentConfig.Namespace,
			Labels:      deploymentLabels(dCtx),
			Annotations: deploymentAnnotations(dCtx),
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "AuthorizationPolicy",
			APIVersion: "security.istio.io/v1beta1",
		},
		Spec: securityv1beta1.AuthorizationPolicy{
			Action: securityv1beta1.AuthorizationPolicy_DENY,
			Selector: &typev1beta1.WorkloadSelector{
				MatchLabels: map[string]string{
					riserLabel("deployment"): dCtx.DeploymentConfig.Name,
				},
			},
			Rules: []*securityv1beta1.Rule{
				{
					To: []*securityv1beta1.Rule_To{
						{
							Operation: &securityv1beta1.Operation{
								Paths: []string{dCtx.DeploymentConfig.App.HealthCheck.Path},
							},
						},
					},
				},
			},
		},
	}
}
