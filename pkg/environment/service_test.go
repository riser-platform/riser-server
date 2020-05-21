package environment

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_Ping(t *testing.T) {
	unixNow := time.Now().Unix()
	environmentRepository := &core.FakeEnvironmentRepository{
		GetFn: func(name string) (*core.Environment, error) {
			assert.Equal(t, "myenv", name)
			return &core.Environment{Name: name}, nil
		},
		SaveFn: func(environment *core.Environment) error {
			assert.Equal(t, "myenv", environment.Name)
			assert.InDelta(t, unixNow, environment.Doc.LastPing.Unix(), 3)
			return nil
		},
	}

	service := service{environments: environmentRepository}

	err := service.Ping("myenv")

	assert.NoError(t, err)
	assert.Equal(t, 1, environmentRepository.GetCallCount)
	assert.Equal(t, 1, environmentRepository.SaveCallCount)
}
func Test_Ping_WhenEnvironmentDoesNotExist_SavesNewEnvironment(t *testing.T) {
	unixNow := time.Now().UTC().Unix()

	environmentRepository := &core.FakeEnvironmentRepository{
		GetFn: func(string) (*core.Environment, error) {
			return nil, core.ErrNotFound
		},
		SaveFn: func(envArg *core.Environment) error {
			assert.Equal(t, "myenv", envArg.Name)
			assert.InDelta(t, unixNow, envArg.Doc.LastPing.Unix(), 3)
			return nil
		},
	}

	service := service{environments: environmentRepository}

	err := service.Ping("myenv")

	assert.NoError(t, err)
	assert.Equal(t, 1, environmentRepository.SaveCallCount)
}

func Test_Ping_WhenGetEnvironmentError_ReturnsError(t *testing.T) {
	environmentRepository := &core.FakeEnvironmentRepository{
		GetFn: func(string) (*core.Environment, error) {
			return nil, errors.New("test")
		},
		SaveFn: func(*core.Environment) error {
			return nil
		},
	}

	service := service{environments: environmentRepository}

	err := service.Ping("myenv")

	assert.Equal(t, "Error retrieving environment \"myenv\": test", err.Error())
}

func Test_Ping_WhenPingEnvironmentError_ReturnsError(t *testing.T) {
	environmentRepository := &core.FakeEnvironmentRepository{
		GetFn: func(string) (*core.Environment, error) {
			return &core.Environment{}, nil
		},
		SaveFn: func(*core.Environment) error {
			return errors.New("test")
		},
	}

	service := service{environments: environmentRepository}

	err := service.Ping("myenv")

	assert.Equal(t, "Error saving environment \"myenv\": test", err.Error())
}

func Test_GetStatus_Healthy(t *testing.T) {
	environmentRepository := &core.FakeEnvironmentRepository{
		GetFn: func(envName string) (*core.Environment, error) {
			assert.Equal(t, "myenv", envName)
			return &core.Environment{
				Name: "myenv",
				Doc: core.EnvironmentDoc{
					LastPing: time.Now(),
				},
			}, nil
		},
	}

	service := service{environments: environmentRepository}

	status, err := service.GetStatus("myenv")

	assert.NoError(t, err)
	assert.Equal(t, "myenv", status.EnvironmentName)
	assert.True(t, status.Healthy)
	assert.Empty(t, status.Reason)
	assert.Equal(t, 1, environmentRepository.GetCallCount)
}

func Test_GetStatus_OldPing_ReturnsUnhealthy(t *testing.T) {
	environmentRepository := &core.FakeEnvironmentRepository{
		GetFn: func(envName string) (*core.Environment, error) {
			assert.Equal(t, "myenv", envName)
			return &core.Environment{
				Name: "myenv",
				Doc: core.EnvironmentDoc{
					LastPing: time.Now().Add(-UnhealthyAfter),
				},
			}, nil
		},
	}

	service := service{environments: environmentRepository}

	status, err := service.GetStatus("myenv")

	assert.NoError(t, err)
	assert.False(t, status.Healthy)
	assert.Equal(t, "The status may be stale. The last status was reported 30 seconds ago.", status.Reason)
}

func Test_GetStatus_GetError_ReturnsError(t *testing.T) {
	environmentRepository := &core.FakeEnvironmentRepository{
		GetFn: func(envName string) (*core.Environment, error) {
			return nil, errors.New("test")
		},
	}

	service := service{environments: environmentRepository}

	status, err := service.GetStatus("myenv")

	assert.Equal(t, "Error retrieving environment \"myenv\": test", err.Error())
	assert.Nil(t, status)
}

func Test_SetConfig(t *testing.T) {
	lastPing := time.Now().UTC()
	config := &core.EnvironmentConfig{
		SealedSecretCert: []byte{0x1},
	}

	environmentRepository := &core.FakeEnvironmentRepository{
		GetFn: func(envName string) (*core.Environment, error) {
			assert.Equal(t, "myenv", envName)
			return &core.Environment{
				Name: "myenv",
				Doc: core.EnvironmentDoc{
					LastPing: lastPing,
					Config: core.EnvironmentConfig{
						SealedSecretCert: []byte{0x2},
					},
				},
			}, nil
		},
		SaveFn: func(environment *core.Environment) error {
			assert.Equal(t, "myenv", environment.Name)
			assert.Equal(t, lastPing, environment.Doc.LastPing)
			assert.Equal(t, config.SealedSecretCert, environment.Doc.Config.SealedSecretCert)
			return nil
		},
	}

	service := service{environmentRepository}

	err := service.SetConfig("myenv", config)

	assert.NoError(t, err)
	assert.Equal(t, 1, environmentRepository.SaveCallCount)
}

func Test_SetConfig_IgnoresEmptyValues(t *testing.T) {
	environmentRepository := &core.FakeEnvironmentRepository{
		GetFn: func(envName string) (*core.Environment, error) {
			assert.Equal(t, "myenv", envName)
			return &core.Environment{
				Name: "myenv",
				Doc: core.EnvironmentDoc{
					Config: core.EnvironmentConfig{
						SealedSecretCert: []byte{0x2},
					},
				},
			}, nil
		},
		SaveFn: func(environment *core.Environment) error {
			assert.Equal(t, "myenv", environment.Name)
			assert.Equal(t, []byte{0x2}, environment.Doc.Config.SealedSecretCert)
			return nil
		},
	}

	service := service{environmentRepository}

	err := service.SetConfig("myenv", &core.EnvironmentConfig{})

	assert.NoError(t, err)
	assert.Equal(t, 1, environmentRepository.SaveCallCount)
}

func Test_ValidateDeployable(t *testing.T) {
	environmentRepository := &core.FakeEnvironmentRepository{
		ListFn: func() ([]core.Environment, error) {
			return []core.Environment{{Name: "myenv1"}}, nil
		},
	}

	service := service{environmentRepository}

	err := service.ValidateDeployable("myenv1")

	assert.NoError(t, err)
}

func Test_ValidateDeployable_WhenEnvironmentDoesNotExist(t *testing.T) {
	environmentRepository := &core.FakeEnvironmentRepository{
		ListFn: func() ([]core.Environment, error) {
			return []core.Environment{{Name: "myenv1"}, {Name: "myenv2"}}, nil
		},
	}

	service := service{environmentRepository}

	err := service.ValidateDeployable("myenv3")

	assert.IsType(t, &core.ValidationError{}, err)
	assert.Equal(t, "Invalid environment. Must be one of: myenv1, myenv2", err.Error())
}

func Test_ValidateDeployable_ReturnsError(t *testing.T) {
	environmentRepository := &core.FakeEnvironmentRepository{
		ListFn: func() ([]core.Environment, error) {
			return nil, errors.New("failed")
		},
	}

	service := service{environmentRepository}

	err := service.ValidateDeployable("myenv")

	assert.Equal(t, "Unable to validate environment: failed", err.Error())
}
