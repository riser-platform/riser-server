package core

import (
	"github.com/google/uuid"
)

type App struct {
	Id   uuid.UUID
	Name string
}

type AppStatus struct {
	StageStatuses []StageStatus
	// Deployments returns the whole deployment. We should probably use a different type here with less data, but we can't just pass
	// Deployment.Doc.Status as we also need the DeploymentName and the Stage.
	Deployments []Deployment
}
