package environment

import (
	"fmt"
	"strings"
	"time"

	"github.com/imdario/mergo"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
)

// UnhealthyAfter indicates the duration that is used to calculate if the environment is unhealthy due to not receiving any type of communication from the environment
var UnhealthyAfter = time.Duration(30) * time.Second

type RepoSettings struct {
	URL         string
	LocalGitDir string
}

type Service interface {
	// Ping "pings" an environment, signaling that the environment has received some form of update
	// This is used to help identify when a cluster is no longer reporting status to the server
	// If the environment has not yet been provisioned this will automatically create the environment (this may change in the future)
	Ping(envName string) error
	GetConfig(envName string) (*core.EnvironmentConfig, error)
	SetConfig(envName string, environment *core.EnvironmentConfig) error
	GetStatus(envName string) (*core.EnvironmentStatus, error)
	ValidateDeployable(envName string) error
}

type service struct {
	environments core.EnvironmentRepository
}

func NewService(environments core.EnvironmentRepository) Service {
	return &service{environments}
}

func (s *service) GetConfig(envName string) (*core.EnvironmentConfig, error) {
	err := s.ValidateDeployable(envName)
	if err != nil {
		return nil, err
	}

	env, err := s.environments.Get(envName)
	if err != nil {
		return nil, err
	}

	return &env.Doc.Config, nil
}

// SetConfig merges any non zero value with the existing environment configuration
// DeleteConfig should be added if we ever need to clear a environment config value
func (s *service) SetConfig(envName string, environmentConfig *core.EnvironmentConfig) error {
	environment, err := s.environments.Get(envName)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error retrieving environment %q", envName))
	}

	err = mergo.MergeWithOverwrite(&environment.Doc.Config, environmentConfig)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error merging environment configuration for environment %q", envName))
	}

	err = s.environments.Save(environment)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error saving environment %q", envName))
	}

	return nil
}

func (s *service) GetStatus(envName string) (*core.EnvironmentStatus, error) {
	environment, err := s.environments.Get(envName)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error retrieving environment %q", envName))
	}

	status := &core.EnvironmentStatus{
		EnvironmentName: envName,
		Healthy:         true,
	}

	lastPingWithinThreshold := environment.Doc.LastPing.After(time.Now().Add(-UnhealthyAfter))
	if !lastPingWithinThreshold {
		status.Healthy = false
		status.Reason = fmt.Sprintf("The status may be stale. The last status was reported %s.", humanize.Time(environment.Doc.LastPing))
	}

	return status, nil
}

// ValidateDeployable validates the existence of a environment and returns a user friendly error with a list of valid environments
// In the future this may become more sophisticated to determine if it's deployable for a given app e.g. based on RBAC, teams, etc.
func (s *service) ValidateDeployable(envName string) error {
	environments, err := s.environments.List()
	if err != nil {
		return errors.Wrap(err, "Unable to validate environment")
	}

	envNames := []string{}
	for _, environment := range environments {
		if environment.Name == envName {
			return nil
		}
		envNames = append(envNames, environment.Name)
	}

	return core.NewValidationErrorMessage(fmt.Sprintf("Invalid environment. Must be one of: %s", strings.Join(envNames, ", ")))
}

func (s *service) Ping(envName string) error {
	environment, err := s.environments.Get(envName)
	if err != nil {
		// The controller should provision the environment with configuration. In the case that it does not exist though,
		// we create the environment here with no configuration.
		// TODO: Revisit after controller work as this may not be necessary in future. If this changes be sure to update comments on interface.
		if err == core.ErrNotFound {
			environment = &core.Environment{Name: envName}
		} else {
			return errors.Wrap(err, fmt.Sprintf("Error retrieving environment %q", envName))
		}
	}

	environment.Doc.LastPing = time.Now().UTC()

	err = s.environments.Save(environment)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error saving environment %q", envName))
	}

	return nil
}
