package model

const (
	RolloutStatusInProgress = "InProgress"
	RolloutStatusComplete   = "Completed"
	RolloutStatusFailed     = "Failed"
	RolloutStatusUnknown    = "Unknown"
)

type AppStatus struct {
	Stages      []StageStatus      `json:"stages"`
	Deployments []DeploymentStatus `json:"deployments"`
}

type DeploymentStatus struct {
	DeploymentName          string `json:"deployment"`
	StageName               string `json:"stage"`
	RiserGeneration         int64  `json:"riserGeneration"`
	DeploymentStatusMutable `json:",inline"`
}

type DeploymentStatusMutable struct {
	ObservedRiserGeneration   int64                      `json:"observedRiserGeneration"`
	Revisions                 []DeploymentRevisionStatus `json:"revisions,omitempty"`
	Traffic                   []DeploymentTrafficStatus  `json:"traffic,omitempty"`
	LatestCreatedRevisionName string                     `json:"latestCreatedRevisionName"`
	LatestReadyRevisionName   string                     `json:"latestReadyRevisionName"`
}

type DeploymentTrafficStatus struct {
	Latest       *bool  `json:"latest,omitempty"`
	Percent      *int64 `json:"percent,omitempty"`
	RevisionName string `json:"revisionName"`
}

type DeploymentRevisionStatus struct {
	Name              string          `json:"name"`
	AvailableReplicas int32           `json:"availableReplicas"`
	DockerImage       string          `json:"dockerImage"`
	RiserGeneration   int64           `json:"riserGeneration"`
	Problems          []StatusProblem `json:"problems,omitempty"`
	// TODO: Probably rename and use different
	RolloutStatus string `json:"rolloutStatus"`
	// TODO: Probably deprecate RolloutStatusReason
	RolloutStatusReason string `json:"rolloutStatusReason"`
}

type StatusProblem struct {
	Count   int    `json:"count"`
	Message string `json:"message"`
}

type StageStatus struct {
	StageName string `json:"stageName"`
	Healthy   bool   `json:"healthy"`
	Reason    string `json:"string"`
}
