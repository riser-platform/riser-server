package state

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/riser-platform/riser-server/pkg/state/resources"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/util"
)

const riserManagedStatePath = "state/riser-managed"

type getResourcePathFunc func(resource KubeResource) string

func RenderDeleteDeployment(deploymentName, namespace string) []core.ResourceFile {
	return []core.ResourceFile{
		{
			Name:   getDeploymentScmDir(deploymentName, namespace),
			Delete: true,
		},
		{
			Name:   getAppConfigScmPath(deploymentName, namespace),
			Delete: true,
		},
	}
}

// RenderGeneric is used for generic resources (e.g. riser app namespaces). They will be placed in the root of the namespaced folder.
func RenderGeneric(environmentName string, resources ...KubeResource) ([]core.ResourceFile, error) {
	return renderKubeResources(func(resource KubeResource) string {
		return getGenericStatePath(resource)
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
	}, filterNilResources(deploymentResources...)...)

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
		Name:     getAppConfigScmPath(deployment.Name, deployment.Namespace),
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

func getDeploymentScmDir(deploymentName, namespace string) string {
	return strings.ToLower(filepath.Join(riserManagedStatePath,
		namespace,
		"deployments",
		deploymentName))
}

func getDeploymentScmPath(deploymentName, namespace, environmentName string, resource KubeResource) string {
	return strings.ToLower(filepath.Join(
		getDeploymentScmDir(deploymentName, namespace),
		getFileNameFromResource(resource)))
}

func getSecretScmPath(app string, environmentName string, sealedSecret KubeResource) string {
	return strings.ToLower(filepath.Join(
		riserManagedStatePath,
		sealedSecret.GetNamespace(),
		"secrets",
		app,
		getFileNameFromResource(sealedSecret)))
}

func getAppConfigScmPath(deploymentName, namespace string) string {
	return strings.ToLower(filepath.Join(
		"riser-config",
		namespace,
		fmt.Sprintf("%s.yaml", deploymentName)))
}

func getGenericStatePath(resource KubeResource) string {
	return strings.ToLower(filepath.Join(
		riserManagedStatePath,
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

// TODO: Consider refactoring the KubeResource interface or adding another type to determine if a resource should be rendered or not
// instead of using nil
func filterNilResources(in ...KubeResource) []KubeResource {
	out := []KubeResource{}
	for _, r := range in {
		if r == nil || reflect.ValueOf(r).IsNil() {
			continue
		}
		out = append(out, r)
	}
	return out
}
