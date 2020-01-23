package model

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/pkg/errors"
)

type RolloutRequest struct {
	Traffic []TrafficRule `json:"traffic"`
}

type TrafficRule struct {
	RiserGeneration int64 `json:"riserGeneration"`
	Percent         int   `json:"percent"`
}

func (rolloutRequest *RolloutRequest) Validate() error {
	var err error
	percentage := 0
	revisions := map[int64]bool{}
	for idx, rule := range rolloutRequest.Traffic {
		percentage += rule.Percent
		if _, ok := revisions[rule.RiserGeneration]; ok {
			err = mergeValidationErrors(err,
				validation.Errors{"riserGeneration": fmt.Errorf("revision \"%d\" specified twice. You may only specify one rule per revision", rule.RiserGeneration)},
				fmt.Sprintf("traffic[%d]", idx))
		}
		revisions[rule.RiserGeneration] = true
		ruleErr := rule.Validate()
		if ruleErr != nil {
			err = mergeValidationErrors(err, ruleErr, fmt.Sprintf("traffic[%d]", idx))
		}
	}

	rolloutErr := validation.ValidateStruct(rolloutRequest, validation.Field(&rolloutRequest.Traffic,
		validation.Required.Error("must specify one or more traffic rules"),
		validation.By(func(interface{}) error {
			if percentage != 100 {
				return errors.New("rule percentages must add up to 100")
			}
			return nil
		}),
	))

	if rolloutErr != nil {
		err = mergeValidationErrors(err, rolloutErr, "")
	}
	return err
}

func (trafficRule *TrafficRule) Validate() error {
	return validation.ValidateStruct(trafficRule,
		validation.Field(&trafficRule.RiserGeneration, validation.Required, validation.Min(0)),
		validation.Field(&trafficRule.Percent, validation.Min(0), validation.Max(100)),
	)
}
