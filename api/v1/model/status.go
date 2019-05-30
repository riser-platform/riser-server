package model

const (
	RolloutStatusInProgress = "InProgress"
	RolloutStatusComplete   = "Completed"
	RolloutStatusFailed     = "Failed"
)

type Status struct {
	Stages      []StageStatus      `json:"stages"`
	Deployments []DeploymentStatus `json:"deployments"`
}

// DeploymentStatus contains a summary of human friendly status for a deployment
type DeploymentStatus struct {
	AppName             string                    `json:"app"`
	DeploymentName      string                    `json:"deployment"`
	StageName           string                    `json:"stage"`
	RolloutStatus       string                    `json:"rolloutStatus"`
	RolloutRevision     int64                     `json:"rolloutRevision"`
	RolloutStatusReason string                    `json:"rolloutStatusReason"`
	DockerImage         string                    `json:"dockerImage"`
	Problems            []DeploymentStatusProblem `json:"problems"`
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
