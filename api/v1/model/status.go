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
	ObservedRiserGeneration int64                      `json:"observedRiserGeneration"`
	Problems                []DeploymentStatusProblem  `json:"problems"`
	Revisions               []DeploymentRevisionStatus `json:"revisions"`
	LatestReadyRevisionName string                     `json:"latestReadyRevisionName"`
}

type DeploymentRevisionStatus struct {
	Name                string `json:"name"`
	DockerImage         string `json:"dockerImage"`
	RiserGeneration     int64  `json:"riserGeneration"`
	RolloutStatus       string `json:"rolloutStatus"`
	RolloutStatusReason string `json:"rolloutStatusReason"`
}

type DeploymentStatusProblem struct {
	Count   int    `json:"count"`
	Message string `json:"message"`
}

type StageStatus struct {
	StageName string `json:"stageName"`
	Healthy   bool   `json:"healthy"`
	Reason    string `json:"string"`
}
