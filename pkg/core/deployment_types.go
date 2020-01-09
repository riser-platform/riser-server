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
	App           *model.AppConfig
	Traffic       TrafficConfig
	ManualRollout bool
}

type DeploymentDocker struct {
	Tag string `json:"tag"`
}

// Needed for serialization to postgres since we do partial updates on traffic
type TrafficConfig []TrafficConfigRule

type TrafficConfigRule struct {
	RiserGeneration int64
	RevisionName    string
	Percent         int64
}

type DeploymentDoc struct {
	Status  *DeploymentStatus   `json:"status,omitempty"`
	Traffic []TrafficConfigRule `json:"traffic"`
}

type DeploymentStatus struct {
	ObservedRiserGeneration   int64                      `json:"observedRiserGeneration"`
	LastUpdated               time.Time                  `json:"lastUpdated"`
	Revisions                 []DeploymentRevisionStatus `json:"revisions"`
	LatestReadyRevisionName   string                     `json:"latestReadyRevisionName"`
	LatestCreatedRevisionName string                     `json:"latestCreatedRevisionName"`
	Traffic                   []DeploymentTrafficStatus  `json:"traffic"`
}

type DeploymentTrafficStatus struct {
	// TODO: Consider removing Latest as we will no longer use it in traffic
	Latest       *bool  `json:"latest,omitempty"`
	Percent      *int64 `json:"percent,omitempty"`
	RevisionName string `json:"revisionName"`
}

type DeploymentRevisionStatus struct {
	Name                string          `json:"name"`
	AvailableReplicas   int32           `json:"availableReplicas"`
	DockerImage         string          `json:"dockerImage"`
	RiserGeneration     int64           `json:"riserGeneration"`
	RolloutStatus       string          `json:"rolloutStatus"`
	RolloutStatusReason string          `json:"rolloutStatusReason"`
	Problems            []StatusProblem `json:"problems"`
}

type StatusProblem struct {
	Count   int    `json:"count"`
	Message string `json:"message"`
}

type DeploymentContext struct {
	Deployment      *DeploymentConfig
	Stage           *StageConfig
	RiserGeneration int64
	SecretNames     []string
	ManualRollout   bool
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

// Needed for sql.Scanner interface. Normally this is only needed on the "Doc" object but we need this here since we do traffic only updates.
func (a TrafficConfig) Value() (driver.Value, error) {
	return json.Marshal(a)
}
