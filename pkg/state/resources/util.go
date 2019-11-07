package resources

import (
	"fmt"
	"strconv"

	"github.com/riser-platform/riser-server/pkg/core"
)

func commonLabels(ctx *core.DeploymentContext) map[string]string {
	return map[string]string{
		riserLabel("deployment"): ctx.Deployment.Name,
		riserLabel("stage"):      ctx.Deployment.Stage,
		riserLabel("app"):        ctx.Deployment.App.Name,
	}
}

func commonAnnotations(ctx *core.DeploymentContext) map[string]string {
	return map[string]string{
		riserLabel("generation"): strconv.FormatInt(ctx.RiserGeneration, 10),
	}
}

// riserLabel returns a fully qualified riser label or annotation (e.g. riser.dev/your-label)
func riserLabel(labelName string) string {
	return fmt.Sprintf("riser.dev/%s", labelName)
}

func int32Ptr(val int32) *int32 {
	return &val
}

func float32Ptr(val float32) *float32 {
	return &val
}

func boolPtr(val bool) *bool {
	return &val
}
