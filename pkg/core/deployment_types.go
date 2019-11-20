package core

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/riser-platform/riser-server/api/v1/model"
)

type Deployment struct {
	Name      string
	StageName string
	AppName   string
	// RiserGeneration is for tracking deployment changes and has no relation to a k8s Deployment generation
	RiserGeneration int64
	Doc             DeploymentDoc
}

type DeploymentConfig struct {
	Name      string
	Namespace string
	Stage     string
	Docker    DeploymentDocker
	// TODO: Move to core and remove api/v1/model dependency
	App *model.AppConfig
}

type DeploymentDocker struct {
	Tag string `json:"tag"`
}

type DeploymentDoc struct {
	Status *DeploymentStatus `json:"status,omitempty"`
}

type DeploymentStatus struct {
	ObservedRiserGeneration int64                      `json:"observedRiserGeneration"`
	Problems                []DeploymentStatusProblem  `json:"problems"`
	LastUpdated             time.Time                  `json:"lastUpdated"`
	Revisions               []DeploymentRevisionStatus `json:"revisions"`
	LatestReadyRevisionName string                     `json:"latestReadyRevisionName"`
	Traffic                 []DeploymentTrafficStatus  `json:"traffic"`
}

type DeploymentTrafficStatus struct {
	Latest       *bool  `json:"latest,omitempty"`
	Percent      *int64 `json:"percent,omitempty"`
	RevisionName string `json:"revisionName"`
}

type DeploymentRevisionStatus struct {
	Name              string `json:"name"`
	AvailableReplicas int32  `json:"availableReplicas"`
	DockerImage       string `json:"dockerImage"`
	RiserGeneration   int64  `json:"riserGeneration"`
	// TODO: Probably delete in favor of reporting traffic  at the status level
	RolloutStatus       string `json:"rolloutStatus"`
	RolloutStatusReason string `json:"rolloutStatusReason"`
}

type DeploymentStatusProblem struct {
	Count   int    `json:"count"`
	Message string `json:"message"`
}

type DeploymentContext struct {
	Deployment      *DeploymentConfig
	Stage           *StageConfig
	RiserGeneration int64
	SecretNames     []string
}

// Needed for sql.Scanner interface
func (a *DeploymentDoc) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Needed for sql.Scanner interface
func (a *DeploymentDoc) Scan(value interface{}) error {
	return jsonbSqlUnmarshal(value, &a)
}

// Needed for sql.Scanner interface. Normally this is only needed on the "Doc" object but we need this here since we do status only updates.
func (a *DeploymentStatus) Value() (driver.Value, error) {
	return json.Marshal(a)
}
