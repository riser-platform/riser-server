package stage

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_Ping(t *testing.T) {
	unixNow := time.Now().Unix()
	stageRepository := &core.FakeStageRepository{
		GetFn: func(name string) (*core.Stage, error) {
			assert.Equal(t, "mystage", name)
			return &core.Stage{Name: name}, nil
		},
		SaveFn: func(stage *core.Stage) error {
			assert.Equal(t, "mystage", stage.Name)
			assert.InDelta(t, unixNow, stage.Doc.LastPing.Unix(), 3)
			return nil
		},
	}

	service := service{stages: stageRepository}

	err := service.Ping("mystage")

	assert.NoError(t, err)
	assert.Equal(t, 1, stageRepository.GetCallCount)
	assert.Equal(t, 1, stageRepository.SaveCallCount)
}
func Test_Ping_WhenStageDoesNotExist_SavesNewStage(t *testing.T) {
	unixNow := time.Now().UTC().Unix()

	stageRepository := &core.FakeStageRepository{
		GetFn: func(string) (*core.Stage, error) {
			return nil, core.ErrNotFound
		},
		SaveFn: func(stageArg *core.Stage) error {
			assert.Equal(t, "mystage", stageArg.Name)
			assert.InDelta(t, unixNow, stageArg.Doc.LastPing.Unix(), 3)
			return nil
		},
	}

	service := service{stages: stageRepository}

	err := service.Ping("mystage")

	assert.NoError(t, err)
	assert.Equal(t, 1, stageRepository.SaveCallCount)
}

func Test_Ping_WhenGetStageError_ReturnsError(t *testing.T) {
	stageRepository := &core.FakeStageRepository{
		GetFn: func(string) (*core.Stage, error) {
			return nil, errors.New("test")
		},
		SaveFn: func(*core.Stage) error {
			return nil
		},
	}

	service := service{stages: stageRepository}

	err := service.Ping("mystage")

	assert.Equal(t, "Error retrieving stage \"mystage\": test", err.Error())
}

func Test_Ping_WhenPingStageError_ReturnsError(t *testing.T) {
	stageRepository := &core.FakeStageRepository{
		GetFn: func(string) (*core.Stage, error) {
			return &core.Stage{}, nil
		},
		SaveFn: func(*core.Stage) error {
			return errors.New("test")
		},
	}

	service := service{stages: stageRepository}

	err := service.Ping("mystage")

	assert.Equal(t, "Error saving stage \"mystage\": test", err.Error())
}

func Test_GetStatus_Healthy(t *testing.T) {
	stageRepository := &core.FakeStageRepository{
		GetFn: func(stageName string) (*core.Stage, error) {
			assert.Equal(t, "mystage", stageName)
			return &core.Stage{
				Name: "mystage",
				Doc: core.StageDoc{
					LastPing: time.Now(),
				},
			}, nil
		},
	}

	service := service{stages: stageRepository}

	status, err := service.GetStatus("mystage")

	assert.NoError(t, err)
	assert.Equal(t, "mystage", status.StageName)
	assert.True(t, status.Healthy)
	assert.Empty(t, status.Reason)
	assert.Equal(t, 1, stageRepository.GetCallCount)
}

func Test_GetStatus_OldPing_ReturnsUnhealthy(t *testing.T) {
	stageRepository := &core.FakeStageRepository{
		GetFn: func(stageName string) (*core.Stage, error) {
			assert.Equal(t, "mystage", stageName)
			return &core.Stage{
				Name: "mystage",
				Doc: core.StageDoc{
					LastPing: time.Now().Add(-UnhealthyAfter),
				},
			}, nil
		},
	}

	service := service{stages: stageRepository}

	status, err := service.GetStatus("mystage")

	assert.NoError(t, err)
	assert.False(t, status.Healthy)
	assert.Equal(t, "The status may be stale. The last status was reported 30 seconds ago.", status.Reason)
}

func Test_GetStatus_GetError_ReturnsError(t *testing.T) {
	stageRepository := &core.FakeStageRepository{
		GetFn: func(stageName string) (*core.Stage, error) {
			return nil, errors.New("test")
		},
	}

	service := service{stages: stageRepository}

	status, err := service.GetStatus("mystage")

	assert.Equal(t, "Error retrieving stage \"mystage\": test", err.Error())
	assert.Nil(t, status)
}

func Test_SetConfig(t *testing.T) {
	lastPing := time.Now().UTC()
	config := &core.StageConfig{
		SealedSecretCert: []byte{0x1},
	}

	stageRepository := &core.FakeStageRepository{
		GetFn: func(stageName string) (*core.Stage, error) {
			assert.Equal(t, "mystage", stageName)
			return &core.Stage{
				Name: "mystage",
				Doc: core.StageDoc{
					LastPing: lastPing,
					Config: core.StageConfig{
						SealedSecretCert: []byte{0x2},
					},
				},
			}, nil
		},
		SaveFn: func(stage *core.Stage) error {
			assert.Equal(t, "mystage", stage.Name)
			assert.Equal(t, lastPing, stage.Doc.LastPing)
			assert.Equal(t, config.SealedSecretCert, stage.Doc.Config.SealedSecretCert)
			return nil
		},
	}

	service := service{stageRepository}

	err := service.SetConfig("mystage", config)

	assert.NoError(t, err)
	assert.Equal(t, 1, stageRepository.SaveCallCount)
}

func Test_SetConfig_IgnoresEmptyValues(t *testing.T) {
	stageRepository := &core.FakeStageRepository{
		GetFn: func(stageName string) (*core.Stage, error) {
			assert.Equal(t, "mystage", stageName)
			return &core.Stage{
				Name: "mystage",
				Doc: core.StageDoc{
					Config: core.StageConfig{
						SealedSecretCert: []byte{0x2},
					},
				},
			}, nil
		},
		SaveFn: func(stage *core.Stage) error {
			assert.Equal(t, "mystage", stage.Name)
			assert.Equal(t, []byte{0x2}, stage.Doc.Config.SealedSecretCert)
			return nil
		},
	}

	service := service{stageRepository}

	err := service.SetConfig("mystage", &core.StageConfig{})

	assert.NoError(t, err)
	assert.Equal(t, 1, stageRepository.SaveCallCount)
}
