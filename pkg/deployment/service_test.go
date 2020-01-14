package deployment

import (
	"github.com/pkg/errors"

	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
)

// Note: See snapshot_test for state based testing of deployment artifacts

func Test_prepareForDeployment_whenNewDeploymentCreates(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "myapp-mydep",
		Stage: "mystage",
		App: &model.AppConfig{
			Name: "myapp",
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetFn: func(deploymentNameArg string, stageNameArg string) (*core.Deployment, error) {
			assert.Equal(t, "myapp-mydep", deploymentNameArg)
			assert.Equal(t, "mystage", stageNameArg)
			return nil, core.ErrNotFound
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
	// Sanity check that defaults are tested. Exhaustive default tests are in defaults_test
	assert.NotNil(t, deployment.App.Expose)
	assert.Equal(t, int64(1), result)
	assert.Equal(t, 1, deploymentRepository.GetCallCount)
	assert.Equal(t, 1, deploymentRepository.CreateCallCount)
}

func Test_prepareForDeployment_whenExistingDeployment(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "myapp-mydep",
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
		UpdateTrafficFn: func(name string, stageName string, riserGeneration int64, traffic core.TrafficConfig) error {
			assert.Equal(t, "myapp-mydep", name)
			assert.Equal(t, "mystage", stageName)
			return nil
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
	assert.Equal(t, 1, deploymentRepository.UpdateTrafficCallCount)
	assert.Equal(t, 0, deploymentRepository.CreateCallCount)
}

func Test_prepareForDeployment_whenIncrementGenerationFails(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "myapp-mydep",
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

func Test_prepareForDeployment_whenUpdateTrafficFails(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "myapp-mydep",
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
			return 1, nil
		},
		UpdateTrafficFn: func(name string, stageName string, riserGeneration int64, traffic core.TrafficConfig) error {
			return errors.New("broke")
		},
	}

	service := service{deployments: deploymentRepository}
	result, err := service.prepareForDeployment(deployment)

	assert.Zero(t, result)
	assert.Equal(t, "Error updating traffic: broke", err.Error())
}

/*
	Scenario. Two apps with the following app and deployment names
	app
	app-foo

	app decides to deploy with the deployment name "app-foo", resulting in a duplicate deployment name of "app-foo" of the app "app-foo"
*/
func Test_prepareForDeployment_whenAppDoesNotOwnDeployment(t *testing.T) {
	deployment := &core.DeploymentConfig{
		Name:  "myapp-mydep",
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
		Name:  "myapp-mydep",
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
		Name:  "myapp-mydep",
		Stage: "mystage",
		App: &model.AppConfig{
			Name: "myapp",
		},
	}

	deploymentRepository := &core.FakeDeploymentRepository{
		GetFn: func(deploymentNameArg string, stageNameArg string) (*core.Deployment, error) {
			assert.Equal(t, "myapp-mydep", deploymentNameArg)
			assert.Equal(t, "mystage", stageNameArg)
			return nil, core.ErrNotFound
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

func Test_computeTraffic_NewDeployment(t *testing.T) {
	cfg := &core.DeploymentConfig{
		Name: "myapp",
	}

	result := computeTraffic(1, cfg, nil)

	assert.Len(t, result, 1)
	assert.EqualValues(t, result[0].RiserGeneration, 1)
	assert.Equal(t, result[0].RevisionName, "myapp-1")
	assert.EqualValues(t, result[0].Percent, 100)
}

// Manual rollout for a first time deployment is effectively not allowed
func Test_computeTraffic_NewDeployment_ManualRollout(t *testing.T) {
	cfg := &core.DeploymentConfig{
		Name:          "myapp",
		ManualRollout: true,
	}

	result := computeTraffic(1, cfg, nil)

	assert.Len(t, result, 1)
	assert.EqualValues(t, result[0].RiserGeneration, 1)
	assert.Equal(t, result[0].RevisionName, "myapp-1")
	assert.EqualValues(t, result[0].Percent, 100)
}

func Test_computeTraffic_ExistingDeployment_ManualRollout(t *testing.T) {
	cfg := &core.DeploymentConfig{
		Name:          "myapp",
		ManualRollout: true,
	}

	existingDeployment := &core.Deployment{
		Doc: core.DeploymentDoc{
			Traffic: core.TrafficConfig{
				core.TrafficConfigRule{
					RiserGeneration: 1,
					RevisionName:    "myapp-1",
					Percent:         100,
				},
			},
		},
	}

	result := computeTraffic(2, cfg, existingDeployment)

	assert.Len(t, result, 2)
	assert.EqualValues(t, result[0].RiserGeneration, 2)
	assert.Equal(t, result[0].RevisionName, "myapp-2")
	assert.EqualValues(t, result[0].Percent, 0)
	assert.EqualValues(t, result[1].RiserGeneration, 1)
	assert.Equal(t, result[1].RevisionName, "myapp-1")
	assert.EqualValues(t, result[1].Percent, 100)
}

func Test_computeTraffic_ExistingDeployment_ManualRollout_RemovesExistingZeroPercentRules(t *testing.T) {
	cfg := &core.DeploymentConfig{
		Name:          "myapp",
		ManualRollout: true,
	}

	existingDeployment := &core.Deployment{
		Doc: core.DeploymentDoc{
			Traffic: core.TrafficConfig{
				core.TrafficConfigRule{
					RiserGeneration: 1,
					RevisionName:    "myapp-1",
					Percent:         100,
				},
				core.TrafficConfigRule{
					RiserGeneration: 2,
					RevisionName:    "myapp-2",
					Percent:         0,
				},
			},
		},
	}

	result := computeTraffic(3, cfg, existingDeployment)

	assert.Len(t, result, 2)
	assert.EqualValues(t, result[0].RiserGeneration, 3)
	assert.Equal(t, result[0].RevisionName, "myapp-3")
	assert.EqualValues(t, result[0].Percent, 0)
	assert.EqualValues(t, result[1].RiserGeneration, 1)
	assert.Equal(t, result[1].RevisionName, "myapp-1")
	assert.EqualValues(t, result[1].Percent, 100)
}

func Test_validateDeploymentConfig_ValidatesName(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"app", nil},
		{"app-good", nil},
		{"app-", core.NewValidationError(`invalid deployment name "app-"`, errors.New(`must be either "app" or start with "app-"`))},
		{"mydep", core.NewValidationError(`invalid deployment name "mydep"`, errors.New(`must be either "app" or start with "app-"`))},
		{"app-b@d", core.NewValidationError(`invalid deployment name "app-b@d"`, errors.New("must be lowercase, alphanumeric, and start with a letter"))},
	}

	for _, tt := range tests {
		deployment := &core.DeploymentConfig{
			Name: tt.name,
			App: &model.AppConfig{
				Name: "app",
			},
		}

		result := validateDeploymentConfig(deployment)

		if tt.err == nil {
			assert.Nil(t, result, tt.name)
		} else {
			assert.IsType(t, &core.ValidationError{}, result, tt.name)
			ve := result.(*core.ValidationError)
			ttve := tt.err.(*core.ValidationError)
			assert.Equal(t, ttve.Error(), ve.Error(), tt.name)
		}
	}
}
