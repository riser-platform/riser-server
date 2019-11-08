package deployment

import (
	"database/sql"

	"github.com/pkg/errors"

	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
)

// Note: See snapshot_test for state based testing of deployment artifacts

func Test_prepareForDeployment_whenNewDeploymentCreates(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "mydep",
		Stage: "mystage",
		App: &model.AppConfig{
			Name: "myapp",
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetFn: func(deploymentNameArg string, stageNameArg string) (*core.Deployment, error) {
			assert.Equal(t, "myapp-mydep", deploymentNameArg)
			assert.Equal(t, "mystage", stageNameArg)
			return nil, sql.ErrNoRows
		},
		CreateFn: func(deploymentArg *core.Deployment) error {
			assert.Equal(t, "myapp-mydep", deploymentArg.Name)
			assert.Equal(t, "mystage", deploymentArg.StageName)
			assert.Equal(t, "myapp", deploymentArg.AppName)
			assert.Equal(t, int64(1), deploymentArg.RiserGeneration)
			return nil
		},
	}

	service := service{deployments: deploymentRepository}
	result, err := service.prepareForDeployment(deployment)

	assert.NoError(t, err)
	// Sanity check that defaults are tested. Exhaustive default tests are in util_test
	assert.NotNil(t, deployment.App.Expose)
	assert.Equal(t, int64(1), result)
	assert.Equal(t, 1, deploymentRepository.GetCallCount)
	assert.Equal(t, 1, deploymentRepository.CreateCallCount)
}

func Test_prepareForDeployment_whenExistingDeployment(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "mydep",
		Stage: "mystage",
		App: &model.AppConfig{
			Name: "myapp",
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetFn: func(deploymentNameArg string, stageNameArg string) (*core.Deployment, error) {
			assert.Equal(t, "myapp-mydep", deploymentNameArg)
			assert.Equal(t, "mystage", stageNameArg)
			return &core.Deployment{Name: "myapp-mydep", StageName: "mystage", AppName: "myapp"}, nil
		},
		IncrementGenerationFn: func(name string, stageName string) (int64, error) {
			assert.Equal(t, "myapp-mydep", name)
			assert.Equal(t, "mystage", stageName)
			return 3, nil
		},
	}

	service := service{deployments: deploymentRepository}
	result, err := service.prepareForDeployment(deployment)

	assert.NoError(t, err)
	// Sanity check that defaults are tested. Exhaustive default tests are in util_test
	assert.NotNil(t, deployment.App.Expose)
	assert.Equal(t, int64(3), result)
	assert.Equal(t, 1, deploymentRepository.GetCallCount)
	assert.Equal(t, 1, deploymentRepository.IncrementGenerationCallCount)
	assert.Equal(t, 0, deploymentRepository.CreateCallCount)
}

func Test_prepareForDeployment_whenIncrementGenerationFails(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "mydep",
		Stage: "mystage",
		App: &model.AppConfig{
			Name: "myapp",
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetFn: func(deploymentNameArg string, stageNameArg string) (*core.Deployment, error) {
			return &core.Deployment{Name: "myapp-mydep", StageName: "mystage", AppName: "myapp"}, nil
		},
		IncrementGenerationFn: func(name string, stageName string) (int64, error) {
			return 0, errors.New("test")
		},
	}

	service := service{deployments: deploymentRepository}
	result, err := service.prepareForDeployment(deployment)

	assert.Zero(t, result)
	assert.Equal(t, "Error incrementing deployment generation: test", err.Error())
}

/*
	Scenario. Two apps with the following app and deployment names
	app
	app-foo

	app decides to deploy with suffix "foo", resulting in a duplicate deployment name of "app-foo"
*/
func Test_prepareForDeployment_whenAppDoesNotOwnDeployment(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "mydep",
		Stage: "mystage",
		App: &model.AppConfig{
			Name: "myapp",
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetFn: func(deploymentNameArg string, stageNameArg string) (*core.Deployment, error) {
			return &core.Deployment{Name: "myapp-mydep", StageName: "mystage", AppName: "myapp-mydep"}, nil
		},
	}

	service := service{deployments: deploymentRepository}
	result, err := service.prepareForDeployment(deployment)

	assert.Zero(t, result)
	assert.Equal(t, `A deployment with the name "myapp-mydep" is owned by app "myapp-mydep"`, err.Error())
	assert.IsType(t, &core.ValidationError{}, err)
}

func Test_prepareForDeployment_whenGetFails(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "mydep",
		Stage: "mystage",
		App: &model.AppConfig{
			Name: "myapp",
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetFn: func(string, string) (*core.Deployment, error) {
			return nil, errors.New("test")
		},
	}

	service := service{deployments: deploymentRepository}
	result, err := service.prepareForDeployment(deployment)

	assert.Zero(t, result)
	assert.Equal(t, `Error retrieving deployment "myapp-mydep" in stage "mystage": test`, err.Error())
}

func Test_prepareForDeployment_whenCreateFails(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "mydep",
		Stage: "mystage",
		App: &model.AppConfig{
			Name: "myapp",
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetFn: func(deploymentNameArg string, stageNameArg string) (*core.Deployment, error) {
			assert.Equal(t, "myapp-mydep", deploymentNameArg)
			assert.Equal(t, "mystage", stageNameArg)
			return nil, sql.ErrNoRows
		},
		CreateFn: func(newDeploymentArg *core.Deployment) error {
			return errors.New("test")
		},
	}

	service := service{deployments: deploymentRepository}
	result, err := service.prepareForDeployment(deployment)

	assert.Zero(t, result)
	assert.Equal(t, `Error creating deployment "myapp-mydep" in stage "mystage": test`, err.Error())
}

func Test_deploy_ValidatesName(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name: "app-b@d",
		App: &model.AppConfig{
			Name: "app",
		},
	}

	result := deploy(&core.DeploymentContext{Deployment: deployment}, nil)

	assert.IsType(t, &core.ValidationError{}, result)
	ve := result.(*core.ValidationError)
	assert.Equal(t, "invalid deployment name \"app-b@d\": must be lowercase, alphanumeric, and start with a letter", ve.Error())
}
