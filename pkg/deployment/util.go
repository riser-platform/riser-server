package deployment

import (
	"fmt"
	"strings"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/riser-platform/riser-server/pkg/core"
)

const DefaultExposeProtocol = "http"

func ApplyDefaults(deployment *core.Deployment) *core.Deployment {
	if deployment.Name == "" {
		deployment.Name = deployment.App.Name
	} else if !strings.EqualFold(deployment.Name, deployment.App.Name) {
		deployment.Name = fmt.Sprintf("%s-%s", deployment.App.Name, deployment.Name)
	}

	// Hard coded until we implement namespace support
	deployment.Namespace = "apps"

	if deployment.App.Expose == nil {
		deployment.App.Expose = &model.AppConfigExpose{}
	}
	if deployment.App.Expose.Protocol == "" {
		deployment.App.Expose.Protocol = DefaultExposeProtocol
	}
	return deployment
}
