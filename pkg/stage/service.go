package stage

import (
	"fmt"
	"strings"
	"time"

	"github.com/imdario/mergo"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
)

// UnhealthyAfter indicates the duration that is used to calculate if the stage is "unhealthy" due to not receiving a "ping" of any kind from it
var UnhealthyAfter = time.Duration(30) * time.Second

type Service interface {
	// Ping "pings" a stage, signaling that the stage has received some form of update
	// This is used to help identify when a cluster is no longer reporting status to the server
	// If the stage has not yet been provisioned this will automatically create the stage (this may change in the future)
	Ping(stageName string) error
	SetConfig(stageName string, stage *core.StageConfig) error
	GetStatus(stageName string) (*core.StageStatus, error)
	ValidateDeployable(stageName string) error
}

type service struct {
	stages core.StageRepository
}

func NewService(stages core.StageRepository) Service {
	return &service{stages}
}

// SetConfig merges any non zero value with the existing stage configuration
// DeleteConfig should be added if we ever need to clear a stage config value
func (s *service) SetConfig(stageName string, stageConfig *core.StageConfig) error {
	stage, err := s.stages.Get(stageName)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error retrieving stage %q", stageName))
	}

	err = mergo.MergeWithOverwrite(&stage.Doc.Config, stageConfig)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error merging stage configuration for stage %q", stageName))
	}

	err = s.stages.Save(stage)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error saving stage %q", stageName))
	}

	return nil
}

func (s *service) GetStatus(stageName string) (*core.StageStatus, error) {
	stage, err := s.stages.Get(stageName)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error retrieving stage %q", stageName))
	}

	status := &core.StageStatus{
		StageName: stageName,
		Healthy:   true,
	}

	lastPingWithinThreshold := stage.Doc.LastPing.After(time.Now().Add(-UnhealthyAfter))
	if !lastPingWithinThreshold {
		status.Healthy = false
		status.Reason = fmt.Sprintf("The status may be stale. The last status was reported %s.", humanize.Time(stage.Doc.LastPing))
	}

	return status, nil
}

// ValidateDeployable validates the existence of a stage and returns a user friendly error with a list of valid stages
// In the future this may become more sophisticated to determine if it's deployable for a given app e.g. based on RBAC, teams, etc.
func (s *service) ValidateDeployable(stageName string) error {
	stages, err := s.stages.List()
	if err != nil {
		return errors.Wrap(err, "Unable to validate stage")
	}

	stageNames := []string{}
	for _, stage := range stages {
		if stage.Name == stageName {
			return nil
		}
		stageNames = append(stageNames, stage.Name)
	}

	return core.NewValidationErrorMessage(fmt.Sprintf("Invalid stage. Must be one of: %s", strings.Join(stageNames, ", ")))
}

func (s *service) Ping(stageName string) error {
	stage, err := s.stages.Get(stageName)
	if err != nil {
		// The controller should provision the stage with configuration. In the case that it does not exist though,
		// we create the stage here with no configuration.
		// TODO: Revisit after controller work as this may not be necessary in future. If this changes be sure to update comments on interface.
		if err == core.ErrNotFound {
			stage = &core.Stage{Name: stageName}
		} else {
			return errors.Wrap(err, fmt.Sprintf("Error retrieving stage %q", stageName))
		}
	}

	stage.Doc.LastPing = time.Now().UTC()

	err = s.stages.Save(stage)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error saving stage %q", stageName))
	}

	return nil
}
