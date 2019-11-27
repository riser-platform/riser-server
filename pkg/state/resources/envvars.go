package resources

import (
	"fmt"
	"github.com/riser-platform/riser-server/pkg/util"
	"sort"
	"strings"

	"github.com/riser-platform/riser-server/pkg/core"
	corev1 "k8s.io/api/core/v1"
)

func k8sEnvVars(ctx *core.DeploymentContext) []corev1.EnvVar {
	envVars := []corev1.EnvVar{}
	for key, val := range ctx.Deployment.App.Environment {
		envVars = append(envVars, corev1.EnvVar{
			Name:  strings.ToUpper(key),
			Value: val.String(),
		})
	}

	for _, secretName := range ctx.SecretNames {
		secretEnv := corev1.EnvVar{
			Name: strings.ToUpper(secretName),
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key:      "data",
					Optional: util.PtrBool(false),
					LocalObjectReference: corev1.LocalObjectReference{
						Name: fmt.Sprintf("%s-%s", ctx.Deployment.App.Name, secretName),
					},
				},
			},
		}

		envVars = append(envVars, secretEnv)
	}

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
