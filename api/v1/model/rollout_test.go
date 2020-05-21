package model

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RolloutRequest_Valid(t *testing.T) {
	rolloutRequest := &RolloutRequest{
		Traffic: []TrafficRule{
			{
				RiserRevision: 1,
				Percent:       10,
			},
			{
				RiserRevision: 2,
				Percent:       90,
			},
		},
	}

	err := rolloutRequest.Validate()

	assert.NoError(t, err)
}

func Test_RolloutRequest_ValidateTrafficRequired(t *testing.T) {
	rolloutRequest := &RolloutRequest{}

	err := rolloutRequest.Validate()

	assert.Equal(t, "traffic: must specify one or more traffic rules.", err.Error())
}

func Test_RolloutRequest_ValidateTraffic100Percent(t *testing.T) {
	rolloutRequest := &RolloutRequest{
		Traffic: []TrafficRule{
			{
				RiserRevision: 1,
				Percent:       10,
			},
			{
				RiserRevision: 2,
				Percent:       80,
			},
		},
	}

	err := rolloutRequest.Validate()

	assert.Equal(t, "traffic: rule percentages must add up to 100.", err.Error())
}

func Test_RolloutRequest_Validate_NoDupeRevisions(t *testing.T) {
	rolloutRequest := &RolloutRequest{
		Traffic: []TrafficRule{
			{
				RiserRevision: 1,
				Percent:       10,
			},
			{
				RiserRevision: 2,
				Percent:       70,
			},
			{
				RiserRevision: 1,
				Percent:       20,
			},
		},
	}

	err := rolloutRequest.Validate()
	require.IsType(t, validation.Errors{}, err)
	validationErrors := err.(validation.Errors)

	assert.Len(t, validationErrors, 1)
	assert.Equal(t, "revision \"1\" specified twice. You may only specify one rule per revision", validationErrors["traffic[2].riserRevision"].Error())
}

func Test_RolloutRequest_ValidateTrafficRule(t *testing.T) {
	rolloutRequest := &RolloutRequest{
		Traffic: []TrafficRule{
			{
				Percent: 105,
			},
			{
				RiserRevision: 2,
				Percent:       -5,
			},
		},
	}

	err := rolloutRequest.Validate()

	require.IsType(t, validation.Errors{}, err)
	validationErrors := err.(validation.Errors)

	assert.Len(t, validationErrors, 3)
	assert.Equal(t, "cannot be blank", validationErrors["traffic[0].riserRevision"].Error())
	assert.Equal(t, "must be no greater than 100", validationErrors["traffic[0].percent"].Error())
	assert.Equal(t, "must be no less than 0", validationErrors["traffic[1].percent"].Error())
}
