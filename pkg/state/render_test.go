package state

import (
	"testing"

	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/riser-platform/riser-server/pkg/state/resources"

	"github.com/stretchr/testify/require"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_getDeploymentScmPath(t *testing.T) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myapp01",
			Namespace: "apps",
		},
		TypeMeta: metav1.TypeMeta{
			Kind: "Deployment",
		},
	}
	meta := core.DeploymentMeta{
		Namespace: "apps",
		Name:      "myapp01",
		Stage:     "dev",
	}

	result := getDeploymentScmPath(meta, deployment)

	assert.Equal(t, "stages/dev/kube-resources/riser-managed/apps/deployments/myapp01/deployment.myapp01.yaml", result)
}

func Test_getAppConfigScmPath(t *testing.T) {
	deployment := &core.Deployment{
		DeploymentMeta: core.DeploymentMeta{
			Name:      "myapp01-test",
			Namespace: "apps",
			Stage:     "dev",
		},
		App: &model.AppConfig{
			Name: "myapp01",
		},
	}

	result := getAppConfigScmPath(deployment)

	assert.Equal(t, "stages/dev/configs/apps/myapp01/myapp01-test.yaml", result)
}

func Test_getSecretScmPath(t *testing.T) {
	secret := &resources.SealedSecret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "apps",
			Name:      "myapp-mysecret",
		},
		TypeMeta: metav1.TypeMeta{
			Kind: "SealedSecret",
		},
	}

	result := getSecretScmPath("myapp", "dev", secret)

	assert.Equal(t, "stages/dev/kube-resources/riser-managed/apps/secrets/myapp/sealedsecret.myapp-mysecret.yaml", result)
}

func Test_renderDeploymentResources(t *testing.T) {
	deployment := &core.Deployment{
		DeploymentMeta: core.DeploymentMeta{
			Name:      "mydeployment",
			Namespace: "apps",
			Stage:     "dev",
		},
		App: &model.AppConfig{
			Name: "myapp01",
		},
	}
	resource1 := &resources.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mydeployment-public-ingress-443",
		},
		TypeMeta: metav1.TypeMeta{
			Kind: "VirtualService",
		},
	}

	resource2 := &resources.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mydeployment-private-ingress-443",
		},
		TypeMeta: metav1.TypeMeta{
			Kind: "VirtualService",
		},
	}

	result, err := RenderDeployment(deployment, resource1, resource2)

	require.NoError(t, err)
	// Sanity check output - we'll use snapshot testing for exhaustive serialization and file system tests
	assert.Len(t, result, 3)
	assert.Equal(t, "stages/dev/kube-resources/riser-managed/apps/deployments/mydeployment/virtualservice.mydeployment-public-ingress-443.yaml", result[0].Name)
	assert.Contains(t, string(result[0].Contents), "name: mydeployment-public-ingress-443")
	assert.Equal(t, "stages/dev/kube-resources/riser-managed/apps/deployments/mydeployment/virtualservice.mydeployment-private-ingress-443.yaml", result[1].Name)
	assert.Contains(t, string(result[1].Contents), "name: mydeployment-private-ingress-443")
	assert.Equal(t, "stages/dev/configs/apps/myapp01/mydeployment.yaml", result[2].Name)
	assert.Contains(t, string(result[2].Contents), "name: myapp01")
}
