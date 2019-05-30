package resources

import (
	"fmt"
	"sort"
	"strings"

	"github.com/riser-platform/riser-server/pkg/core"
	corev1 "k8s.io/api/core/v1"
)

func k8sEnvVars(deployment *core.Deployment, secretNames []string) []corev1.EnvVar {
	envVars := []corev1.EnvVar{}
	for key, val := range deployment.App.Environment {
		envVars = append(envVars, corev1.EnvVar{
			Name:  strings.ToUpper(key),
			Value: val.String(),
		})
	}

	for _, secretName := range secretNames {
		secretEnv := corev1.EnvVar{
			Name: strings.ToUpper(secretName),
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					// TODO: Instaed of hard coding this use SecretName or somehting and add GetKeySelector() (even though this will likely always return "data")
					Key:      "data",
					Optional: boolPtr(false),
					LocalObjectReference: corev1.LocalObjectReference{
						// TODO: Instead of using a string use a SecretName or something with a GetEnvName() and GetSecretObjectName() functions on it
						Name: fmt.Sprintf("%s-%s", deployment.App.Name, secretName),
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
