package model

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v3"
)

func mergeValidationErrors(baseError error, toMerge error, fieldPrefix string) error {
	if toMerge == nil {
		return baseError
	}

	var isValidationErrors bool
	baseValidationErrors := validation.Errors{}

	if baseError != nil {
		baseValidationErrors, isValidationErrors = baseError.(validation.Errors)
		if !isValidationErrors {
			return baseError
		}
	}

	toMergeValidationErrors, isValidationErrors := toMerge.(validation.Errors)
	if !isValidationErrors {
		return toMerge
	}

	for k, v := range toMergeValidationErrors {
		fieldName := k
		if fieldPrefix != "" {
			fieldName = fmt.Sprintf("%s.%s", fieldPrefix, k)
		}
		baseValidationErrors[fieldName] = v
	}

	return baseValidationErrors

}
