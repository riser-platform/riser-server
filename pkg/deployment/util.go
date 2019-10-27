package deployment

import (
	"fmt"
	"strings"

	"github.com/riser-platform/riser-server/api/v1/model"

	"github.com/imdario/mergo"

	"github.com/riser-platform/riser-server/pkg/core"
)

const DefaultExposeProtocol = "http"

func ApplyDefaults(deployment *core.NewDeployment) *core.NewDeployment {
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

func ApplyOverrides(in *core.NewDeployment) (*core.Deployment, error) {
	out := &core.Deployment{
		DeploymentMeta: in.DeploymentMeta,
	}
	app := in.App.AppConfig
	if overrideApp, ok := in.App.Overrides[in.Stage]; ok {
		err := mergo.Merge(&overrideApp, app)
		if err != nil {
			return nil, err
		}

		app = overrideApp
	}

	out.App = &app

	return out, nil
}
