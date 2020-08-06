package core

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/riser-platform/riser-server/api/v1/model"
)

// Deployment represents a deployment in a particular environment
type Deployment struct {
	DeploymentReservation
	DeploymentRecord
}

// DeploymentRecord represents the database fields specific for a deployment
type DeploymentRecord struct {
	Id              uuid.UUID
	ReservationId   uuid.UUID
	EnvironmentName string
	// RiserRevision is for tracking deployment changes and has no relation to a k8s deployment revision
	RiserRevision int64
	DeletedAt     *time.Time
	Doc           DeploymentDoc
}

type DeploymentConfig struct {
	Name            string
	Namespace       string
	EnvironmentName string
	Docker          DeploymentDocker
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
	RiserRevision int64  `json:"riserRevision"`
	RevisionName  string `json:"revisionName"`
	Percent       int    `json:"percent"`
}

type DeploymentDoc struct {
	Status  *DeploymentStatus   `json:"status,omitempty"`
	Traffic []TrafficConfigRule `json:"traffic"`
}

type DeploymentStatus struct {
	ObservedRiserRevision     int64                      `json:"observedRiserRevision"`
	LastUpdated               time.Time                  `json:"lastUpdated"`
	Revisions                 []DeploymentRevisionStatus `json:"revisions"`
	LatestReadyRevisionName   string                     `json:"latestReadyRevisionName"`
	LatestCreatedRevisionName string                     `json:"latestCreatedRevisionName"`
	Traffic                   []DeploymentTrafficStatus  `json:"traffic"`
}

type DeploymentTrafficStatus struct {
	Percent      *int64 `json:"percent,omitempty"`
	RevisionName string `json:"revisionName"`
	Tag          string `json:"tag,omitempty"`
}

type DeploymentRevisionStatus struct {
	Name                 string `json:"name"`
	DockerImage          string `json:"dockerImage"`
	RiserRevision        int64  `json:"riserRevision"`
	RevisionStatus       string `json:"revisionStatus"`
	RevisionStatusReason string `json:"revisionStatusReason"`
}

type DeploymentContext struct {
	DeploymentConfig  *DeploymentConfig
	EnvironmentConfig *EnvironmentConfig
	RiserRevision     int64
	Secrets           []SecretMeta
	ManualRollout     bool
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
