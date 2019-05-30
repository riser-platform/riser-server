package core

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type DeploymentStatus struct {
	AppName        string
	DeploymentName string
	StageName      string
	Doc            *DeploymentStatusDoc
}

type DeploymentStatusDoc struct {
	RolloutStatus       string                    `json:"rolloutStatus"`
	RolloutStatusReason string                    `json:"rolloutStatusReason"`
	RolloutRevision     int64                     `json:"revision"`
	DockerImage         string                    `json:"dockerImage"`
	Problems            []DeploymentStatusProblem `json:"problems"`
	LastUpdated         time.Time                 `json:"lastUpdated"`
}

type DeploymentStatusProblem struct {
	Count   int    `json:"count"`
	Message string `json:"message"`
}

type DeploymentStatusSummary struct {
	StageStatuses      []StageStatus
	DeploymentStatuses []DeploymentStatus
}

type StageStatus struct {
	StageName string
	Healthy   bool
	Reason    string
}

// Needed for sql.Scanner interface
func (a *DeploymentStatusDoc) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Needed for sql.Scanner interface
func (a *DeploymentStatusDoc) Scan(value interface{}) error {
	return jsonbSqlUnmarshal(value, &a)
}
