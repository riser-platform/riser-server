package model

import "github.com/google/uuid"

const (
	RevisionStatusWaiting   = "Waiting"
	RevisionStatusReady     = "Ready"
	RevisionStatusUnhealthy = "Unhealthy"
	RevisionStatusUnknown   = "Unknown"
)

type AppStatus struct {
	Stages      []StageStatus      `json:"stages"`
	Deployments []DeploymentStatus `json:"deployments"`
}

type DeploymentStatus struct {
	AppId                   uuid.UUID `json:"appId"`
	DeploymentName          string    `json:"deployment"`
	Namespace               string    `json:"namespace"`
	StageName               string    `json:"stage"`
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
	RevisionStatus       string `json:"rolloutStatus"`
	RevisionStatusReason string `json:"revisionStatusReason,omitempty"`
}

type StageStatus struct {
	StageName string `json:"stageName"`
	Healthy   bool   `json:"healthy"`
	Reason    string `json:"string"`
}
