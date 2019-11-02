package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/riser-platform/riser-server/pkg/core"
)

func Test_deploy_ValidatesName(t *testing.T) {
	deployment := &core.Deployment{
		DeploymentMeta: core.DeploymentMeta{
			Name: "b@d",
		},
		App: &model.AppConfig{
			Name: "app",
		},
	}

	result := deploy(deployment, core.StageConfig{}, nil, nil)

	assert.IsType(t, &core.ValidationError{}, result)
	ve := result.(*core.ValidationError)
	assert.Equal(t, "invalid deployment name \"app-b@d\": must be lowercase, alphanumeric, and start with a letter", ve.Error())
}
