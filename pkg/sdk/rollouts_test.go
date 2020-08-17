package sdk

import (
	"errors"
	"net/http"
	"testing"

	"github.com/riser-platform/riser-server/api/v1/model"
	"github.com/stretchr/testify/assert"
)

func Test_Rollouts_Save(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/rollout/dev/myns/myapp", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		rollout := &model.RolloutRequest{}
		mustUnmarshalR(r.Body, rollout)
		assert.Len(t, rollout.Traffic, 2)
		assert.EqualValues(t, 1, rollout.Traffic[0].RiserRevision)
		assert.EqualValues(t, 10, rollout.Traffic[0].Percent)
		assert.EqualValues(t, 2, rollout.Traffic[1].RiserRevision)
		assert.EqualValues(t, 90, rollout.Traffic[1].Percent)

	})

	err := client.Rollouts.Save("myapp", "myns", "dev", "r1:10", "r2:90")

	assert.NoError(t, err)
}

func Test_parseTrafficRules(t *testing.T) {
	tests := []struct {
		trafficRules  []string
		expectedRules []model.TrafficRule
		expectedErr   error
	}{
		{
			trafficRules: []string{"r1:10", "r2:90"},
			expectedRules: []model.TrafficRule{
				{RiserRevision: 1, Percent: 10},
				{RiserRevision: 2, Percent: 90},
			},
		},
		{
			trafficRules: []string{"r1:10", "bad:90"},
			expectedErr:  errors.New(`Rules must be in the format of "r(rev):(percentage)" e.g. "r1:100" routes 100% of traffic to rev 1`),
		},
	}

	for _, test := range tests {
		parsedRules, err := parseTrafficRules(test.trafficRules...)

		assert.ElementsMatch(t, test.expectedRules, parsedRules)
		assert.Equal(t, test.expectedErr, err)
	}
}
