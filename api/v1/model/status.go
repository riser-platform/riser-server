package model

import "github.com/google/uuid"

const (
	RevisionStatusWaiting   = "Waiting"
	RevisionStatusReady     = "Ready"
	RevisionStatusUnhealthy = "Unhealthy"
	RevisionStatusUnknown   = "Unknown"
)

type AppStatus struct {
	Environments []EnvironmentStatus `json:"environments"`
	Deployments  []DeploymentStatus  `json:"deployments"`
}

type DeploymentStatus struct {
	AppId                   uuid.UUID `json:"appId"`
	DeploymentName          string    `json:"deployment"`
	Namespace               string    `json:"namespace"`
	EnvironmentName         string    `json:"environment"`
	RiserRevision           int64     `json:"riserRevision"`
	DeploymentStatusMutable `json:",inline"`
}

type DeploymentStatusMutable struct {
	ObservedRiserRevision     int64                      `json:"observedRiserRevision"`
	Revisions                 []DeploymentRevisionStatus `json:"revisions,omitempty"`
	Traffic                   []DeploymentTrafficStatus  `json:"traffic,omitempty"`
	LatestCreatedRevisionName string                     `json:"latestCreatedRevisionName"`
	LatestReadyRevisionName   string                     `json:"latestReadyRevisionName"`
}

type DeploymentTrafficStatus struct {
	Percent      *int64 `json:"percent,omitempty"`
	RevisionName string `json:"revisionName"`
	Tag          string `json:"tag,omitempty"`
}

type DeploymentRevisionStatus struct {
	Name                 string `json:"name"`
	AvailableReplicas    int32  `json:"availableReplicas"`
	DockerImage          string `json:"dockerImage"`
	RiserRevision        int64  `json:"riserRevision"`
	RevisionStatus       string `json:"revisionStatus"`
	RevisionStatusReason string `json:"revisionStatusReason,omitempty"`
}

type EnvironmentStatus struct {
	EnvironmentName string `json:"environmentName"`
	Healthy         bool   `json:"healthy"`
	Reason          string `json:"string"`
}
