package core

import (
	"github.com/riser-platform/riser-server/api/v1/model"
)

type NewDeployment struct {
	DeploymentMeta `json:",inline"`
	// TODO: Move to core and remove api/v1/model dependency
	App *model.AppConfigWithOverrides
}

type Deployment struct {
	DeploymentMeta `json:",inline"`
	// TODO: Move to core and remove api/v1/model dependency
	App *model.AppConfig
}

type DeploymentMeta struct {
	Name      string
	Namespace string
	Stage     string
	Docker    DeploymentDocker
}

type DeploymentDocker struct {
	Tag string `json:"tag"`
}
