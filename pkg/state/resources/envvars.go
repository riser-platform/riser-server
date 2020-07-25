package resources

import (
	"fmt"
	"sort"
	"strings"

	"github.com/riser-platform/riser-server/pkg/util"

	"github.com/riser-platform/riser-server/pkg/core"
	corev1 "k8s.io/api/core/v1"
)

func k8sEnvVars(ctx *core.DeploymentContext) []corev1.EnvVar {
	envVars := []corev1.EnvVar{}
	// User defined  vars
	for key, val := range ctx.DeploymentConfig.App.Environment {
		envVars = append(envVars, corev1.EnvVar{
			Name:  strings.ToUpper(key),
			Value: val.String(),
		})
	}

	// Secret vars
	for _, secret := range ctx.Secrets {
		secretEnv := corev1.EnvVar{
			Name: strings.ToUpper(secret.Name),
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key:      "data",
					Optional: util.PtrBool(false),
					LocalObjectReference: corev1.LocalObjectReference{
						Name: fmt.Sprintf("%s-%s-%d", ctx.DeploymentConfig.App.Name, secret.Name, secret.Revision),
					},
				},
			},
		}

		envVars = append(envVars, secretEnv)
	}

	// Platform vars
	envVars = append(envVars,
		corev1.EnvVar{Name: "RISER_APP", Value: string(ctx.DeploymentConfig.App.Name)},
		corev1.EnvVar{Name: "RISER_DEPLOYMENT", Value: string(ctx.DeploymentConfig.Name)},
		corev1.EnvVar{Name: "RISER_DEPLOYMENT_REVISION", Value: fmt.Sprintf("%d", ctx.RiserRevision)},
		corev1.EnvVar{Name: "RISER_ENVIRONMENT", Value: string(ctx.DeploymentConfig.EnvironmentName)},
		corev1.EnvVar{Name: "RISER_NAMESPACE", Value: string(ctx.DeploymentConfig.Namespace)})
	sort.Sort(&envVarSorter{items: envVars})
	return envVars
}

type envVarSorter struct {
	items []corev1.EnvVar
}

func (s *envVarSorter) Len() int {
	return len(s.items)
}

func (s *envVarSorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

func (s *envVarSorter) Less(i, j int) bool {
	return strings.Compare(s.items[i].Name, s.items[j].Name) < 0
}
