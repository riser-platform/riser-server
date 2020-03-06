package core

import (
	"github.com/google/uuid"
)

// DeploymentReservation represents a reservation of a deployment name to an app
type DeploymentReservation struct {
	Id    uuid.UUID
	AppId uuid.UUID
	Name  string
	// At the time of writing the namespace should never be different from the app's namespace as we do not allow a deployment
	// to be deployed to a different namespace than the app. This may change in the future
	Namespace string
}
