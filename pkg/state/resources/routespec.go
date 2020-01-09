package resources

import (
	"fmt"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/util"
)

func createRouteSpec(trafficConfig core.TrafficConfig) RouteSpec {
	spec := RouteSpec{
		Traffic: []TrafficTarget{},
	}

	for _, rule := range trafficConfig {
		spec.Traffic = append(spec.Traffic, TrafficTarget{
			RevisionName: rule.RevisionName,
			Percent:      util.PtrInt64(rule.Percent),
			Tag:          fmt.Sprintf("r%d", rule.RiserGeneration),
		})
	}

	return spec
}
