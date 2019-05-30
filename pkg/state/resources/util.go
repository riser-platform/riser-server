package resources

import "github.com/riser-platform/riser-server/pkg/core"

// TODO: Move this to AppDefaults or something
var defaultRiserAppVersion = "v1"

func commonLabels(deployment *core.Deployment) map[string]string {
	return map[string]string{
		"deployment": deployment.Name,
		"stage":      deployment.Stage,
		"app":        deployment.App.Name,
		"riser-app":  defaultRiserAppVersion,
	}
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
