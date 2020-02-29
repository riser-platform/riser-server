package core

import (
	"strings"

	"github.com/google/uuid"
)

type App struct {
	Id        uuid.UUID
	Name      string
	Namespace string
}

type AppStatus struct {
	StageStatuses []StageStatus
	// Deployments returns the whole deployment. We should probably use a different type here with less data, but we can't just pass
	// Deployment.Doc.Status as we also need the DeploymentName and the Stage.
	Deployments []Deployment
}

type AppIdOrName string

func (v *AppIdOrName) IdValue() (idValue uuid.UUID, hasValue bool) {
	idValue, _ = uuid.Parse(string(*v))
	return idValue, idValue != uuid.Nil
}

func (v *AppIdOrName) NameValue() (nameValue *NamespacedName, hasValue bool) {
	parts := strings.Split(string(*v), ".")
	if len(parts) == 2 {
		return &NamespacedName{parts[0], parts[1]}, true
	}

	return nil, false
}
