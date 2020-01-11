package model

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_mergeValidationErrors(t *testing.T) {
	errA := validation.Errors{}
	errA["field1"] = errors.New("field1 error")
	errB := validation.Errors{}
	errB["field2"] = errors.New("field2 error")

	result := mergeValidationErrors(errA, errB, "b")

	require.IsType(t, validation.Errors{}, result)
	validationErrors := result.(validation.Errors)
	assert.Len(t, validationErrors, 2)
	assert.Equal(t, "field1 error", validationErrors["field1"].Error())
	assert.Equal(t, "field2 error", validationErrors["b.field2"].Error())
}

func Test_mergeValidationErrors_BaseNotValidationError(t *testing.T) {
	errA := errors.New("internal error")
	errB := validation.Errors{}
	errB["field2"] = errors.New("field2 error")

	result := mergeValidationErrors(errA, errB, "b")

	assert.Equal(t, errA, result)
}

func Test_mergeValidationErrors_ToMergeNotValidationError(t *testing.T) {
	errA := validation.Errors{}
	errB := errors.New("internal error")

	result := mergeValidationErrors(errA, errB, "b")

	assert.Equal(t, errB, result)
}

func Test_mergeValidationErrors_NilBase(t *testing.T) {
	errB := validation.Errors{}
	errB["field2"] = errors.New("field2 error")

	result := mergeValidationErrors(nil, errB, "b")

	require.IsType(t, validation.Errors{}, result)
	validationErrors := result.(validation.Errors)
	assert.Len(t, validationErrors, 1)
	assert.Equal(t, "field2 error", validationErrors["b.field2"].Error())
}

func Test_mergeValidationErrors_NilToMerge(t *testing.T) {
	errA := validation.Errors{}
	errA["field1"] = errors.New("field1 error")

	result := mergeValidationErrors(errA, nil, "b")

	require.IsType(t, validation.Errors{}, result)
	validationErrors := result.(validation.Errors)
	assert.Len(t, validationErrors, 1)
	assert.Equal(t, "field1 error", validationErrors["field1"].Error())
}

func Test_mergeValidationErrors_EmptyPrefix(t *testing.T) {
	errA := validation.Errors{}
	errA["field1"] = errors.New("field1 error")
	errB := validation.Errors{}
	errB["field2"] = errors.New("field2 error")

	result := mergeValidationErrors(errA, errB, "")

	require.IsType(t, validation.Errors{}, result)
	validationErrors := result.(validation.Errors)
	assert.Len(t, validationErrors, 2)
	assert.Equal(t, "field1 error", validationErrors["field1"].Error())
	assert.Equal(t, "field2 error", validationErrors["field2"].Error())
}
