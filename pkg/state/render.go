package state

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/riser-platform/riser-server/pkg/state/resources"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/util"
)

type getResourcePathFunc func(resource KubeResource) string

func RenderDeleteDeployment(deploymentName, namespace, environmentName string) []core.ResourceFile {
	return []core.ResourceFile{
		{
			Name:   getDeploymentScmDir(deploymentName, namespace, environmentName),
			Delete: true,
		},
		{
			Name:   getAppConfigScmPath(deploymentName, namespace, environmentName),
			Delete: true,
		},
	}
}

// RenderGeneric is used for generic resources (e.g. riser app namespaces). They will be placed in the root of the namespaced folder.
func RenderGeneric(environmentName string, resources ...KubeResource) ([]core.ResourceFile, error) {
	return renderKubeResources(func(resource KubeResource) string {
		return getGenericStatePath(environmentName, resource)
	}, resources...)
}

func RenderSealedSecret(app, environmentName string, sealedSecret *resources.SealedSecret) ([]core.ResourceFile, error) {
	return renderKubeResources(func(resource KubeResource) string {
		return getSecretScmPath(app, environmentName, sealedSecret)
	}, sealedSecret)
}

// RenderDeployment renders resources that target a deployment's git folder
func RenderDeployment(deployment *core.DeploymentConfig, deploymentResources ...KubeResource) ([]core.ResourceFile, error) {
	files, err := renderKubeResources(func(resource KubeResource) string {
		return getDeploymentScmPath(deployment.Name, deployment.Namespace, deployment.EnvironmentName, resource)
	}, deploymentResources...)

	if err != nil {
		return nil, err
	}

	appConfigFile, err := renderAppConfig(deployment)
	if err != nil {
		return nil, err
	}
	files = append(files, *appConfigFile)

	return files, nil
}

// RenderRoute renders just the route resource.
func RenderRoute(deploymentName, namespace, environmentName string, resource KubeResource) ([]core.ResourceFile, error) {
	files, err := renderKubeResources(func(resource KubeResource) string {
		return getDeploymentScmPath(deploymentName, namespace, environmentName, resource)
	}, resource)

	if err != nil {
		return nil, err
	}

	return files, nil
}

func renderAppConfig(deployment *core.DeploymentConfig) (*core.ResourceFile, error) {
	serialized, err := util.ToYaml(deployment.App)
	if err != nil {
		return nil, errors.Wrap(err, "Error serializing app config")
	}
	return &core.ResourceFile{
		Name:     getAppConfigScmPath(deployment.Name, deployment.Namespace, deployment.EnvironmentName),
		Contents: serialized,
	}, nil
}

func renderKubeResources(pathFunc getResourcePathFunc, resources ...KubeResource) ([]core.ResourceFile, error) {
	files := []core.ResourceFile{}
	for _, resource := range resources {
		serialized, err := util.ToYaml(resource)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error serializing resource %q", resource.GetObjectKind()))
		}
		files = append(files, core.ResourceFile{
			Name:     pathFunc(resource),
			Contents: serialized,
		})
	}
	return files, nil
}

func getDeploymentScmDir(deploymentName, namespace, environmentName string) string {
	return strings.ToLower(filepath.Join(getRiserManagedStatePath(environmentName),
		namespace,
		"deployments",
		deploymentName))
}

func getDeploymentScmPath(deploymentName, namespace, environmentName string, resource KubeResource) string {
	return strings.ToLower(filepath.Join(
		getDeploymentScmDir(deploymentName, namespace, environmentName),
		getFileNameFromResource(resource)))
}

func getSecretScmPath(app string, environmentName string, sealedSecret KubeResource) string {
	return strings.ToLower(filepath.Join(
		getRiserManagedStatePath(environmentName),
		sealedSecret.GetNamespace(),
		"secrets",
		app,
		getFileNameFromResource(sealedSecret)))
}

func getAppConfigScmPath(deploymentName, namespace, environmentName string) string {
	return strings.ToLower(filepath.Join(
		"riser-config",
		environmentName,
		namespace,
		fmt.Sprintf("%s.yaml", deploymentName)))
}

func getRiserManagedStatePath(envName string) string {
	return strings.ToLower(filepath.Join(
		"state",
		envName,
		"riser-managed",
	))
}

func getGenericStatePath(envName string, resource KubeResource) string {
	return strings.ToLower(filepath.Join(
		getRiserManagedStatePath(envName),
		resource.GetNamespace(),
		getFileNameFromResource(resource),
	))
}

func getFileNameFromResource(resource KubeResource) string {
	group := resource.GetObjectKind().GroupVersionKind().GroupVersion().Group
	if group == "" {
		return strings.ToLower(fmt.Sprintf("%s.%s.yaml", resource.GetObjectKind().GroupVersionKind().Kind, resource.GetName()))
	}

	return strings.ToLower(fmt.Sprintf("%s.%s.%s.yaml", group, resource.GetObjectKind().GroupVersionKind().Kind, resource.GetName()))
}
