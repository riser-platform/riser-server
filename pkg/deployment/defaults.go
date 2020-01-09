package deployment

import (
	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/riser-platform/riser-server/pkg/core"
)

const DefaultExposeProtocol = "http"

func applyDefaults(deploymentConfig *core.DeploymentConfig) {
	// Hard coded until we implement namespace support
	deploymentConfig.Namespace = "apps"

	if deploymentConfig.App.Expose == nil {
		deploymentConfig.App.Expose = &model.AppConfigExpose{}
	}
	if deploymentConfig.App.Expose.Protocol == "" {
		deploymentConfig.App.Expose.Protocol = DefaultExposeProtocol
	}
}
