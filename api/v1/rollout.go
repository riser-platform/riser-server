package v1

import (
	"fmt"
	"github.com/riser-platform/riser-server/pkg/stage"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/rollout"
)

func PutRollout(c echo.Context, rolloutService rollout.Service, stageService stage.Service) error {
	rolloutRequest := &model.RolloutRequest{}

	err := c.Bind(rolloutRequest)
	if err != nil {
		return err
	}

	deploymentName := c.Param("deploymentName")
	stageName := c.Param("stageName")

	err = stageService.ValidateDeployable(stageName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = validation.Validate(&rolloutRequest)
	if err != nil {
		return core.NewValidationError("Invalid rollout request", err)
	}

	return rolloutService.UpdateTraffic(deploymentName, stageName, mapTrafficRulesToDomain(deploymentName, rolloutRequest.Traffic))
}

func mapTrafficRulesToDomain(deploymentName string, traffic []model.TrafficRule) core.TrafficConfig {
	out := core.TrafficConfig{}
	for _, rule := range traffic {
		out = append(out, core.TrafficConfigRule{
			RiserGeneration: rule.RiserGeneration,
			RevisionName:    fmt.Sprintf("%s-%d", deploymentName, rule.RiserGeneration),
			Percent:         rule.Percent,
		})
	}
	return out
}
