package deployment

import (
	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/riser-platform/riser-server/pkg/core"
)

const DefaultExposeProtocol = "http"

func applyDefaults(deployment *core.DeploymentConfig) {

	// Hard coded until we implement namespace support
	deployment.Namespace = "apps"

	if deployment.App.Expose == nil {
		deployment.App.Expose = &model.AppConfigExpose{}
	}
	if deployment.App.Expose.Protocol == "" {
		deployment.App.Expose.Protocol = DefaultExposeProtocol
	}
}
