package resources

import (
	"fmt"
	"strconv"

	"github.com/riser-platform/riser-server/pkg/core"
)

// deploymentLabels are labels common to Riser deployment resources
func deploymentLabels(ctx *core.DeploymentContext) map[string]string {
	return map[string]string{
		riserLabel("deployment"): ctx.Deployment.Name,
		riserLabel("stage"):      ctx.Deployment.Stage,
		riserLabel("app"):        ctx.Deployment.App.Name,
	}
}

// deploymentAnnotations are annotations common to Riser deployment resources
func deploymentAnnotations(ctx *core.DeploymentContext) map[string]string {
	return map[string]string{
		riserLabel("generation"): strconv.FormatInt(ctx.RiserGeneration, 10),
	}
}

// riserLabel returns a fully qualified riser label or annotation (e.g. riser.dev/your-label)
func riserLabel(labelName string) string {
	return fmt.Sprintf("riser.dev/%s", labelName)
}
