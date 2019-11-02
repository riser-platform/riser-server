package resources

import (
	"fmt"

	"github.com/riser-platform/riser-server/pkg/core"
)

func commonLabels(deployment *core.Deployment) map[string]string {
	return map[string]string{
		riserLabel("deployment"): deployment.Name,
		riserLabel("stage"):      deployment.Stage,
		riserLabel("app"):        deployment.App.Name,
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
