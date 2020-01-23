package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/stage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/assert"
)

func Test_PutRollout_ValidatesStage(t *testing.T) {
	rollout := model.RolloutRequest{}
	req := httptest.NewRequest(http.MethodPut, "/", testMarshal(rollout))
	req.Header.Add("CONTENT-TYPE", "application/json")
	ctx, _ := newContextWithRecorder(req)

	stageService := &stage.FakeService{
		ValidateDeployableFn: func(stageName string) error {
			return errors.New("test")
		},
	}

	err := PutRollout(ctx, nil, stageService, nil)

	assert.IsType(t, &echo.HTTPError{}, err)
	echoErr := err.(*echo.HTTPError)
	assert.Equal(t, http.StatusBadRequest, echoErr.Code)
	assert.Equal(t, "test", echoErr.Message)
}

func Test_PutRollout_ValidatesTraffic(t *testing.T) {
	rollout := model.RolloutRequest{}
	req := httptest.NewRequest(http.MethodPut, "/", testMarshal(rollout))
	req.Header.Add("CONTENT-TYPE", "application/json")
	ctx, _ := newContextWithRecorder(req)

	stageService := &stage.FakeService{
		ValidateDeployableFn: func(stageName string) error {
			return nil
		},
	}

	err := PutRollout(ctx, nil, stageService, nil)

	assert.Equal(t, "Invalid rollout request: traffic: must specify one or more traffic rules.", err.Error())
}

func Test_mapTrafficRulesToDomain(t *testing.T) {
	in := []model.TrafficRule{
		model.TrafficRule{
			RiserGeneration: 1,
			Percent:         10,
		},
		model.TrafficRule{
			RiserGeneration: 2,
			Percent:         90,
		},
	}

	result := mapTrafficRulesToDomain("myapp", in)

	assert.Len(t, result, 2)
	assert.EqualValues(t, 1, result[0].RiserGeneration)
	assert.Equal(t, "myapp-1", result[0].RevisionName)
	assert.Equal(t, 10, result[0].Percent)
	assert.EqualValues(t, 2, result[1].RiserGeneration)
	assert.Equal(t, "myapp-2", result[1].RevisionName)
	assert.Equal(t, 90, result[1].Percent)
}
