package postgres

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_DeploymentRepository_handleUpdateStatusResult(t *testing.T) {
	failedErr := errors.New("failed")
	tt := []struct {
		f        *fakeSqlResult
		expected error
	}{
		{
			f: &fakeSqlResult{
				RowsAffectedFn: func() (int64, error) {
					return 1, nil
				},
			},
			expected: nil,
		},
		{
			f: &fakeSqlResult{
				RowsAffectedFn: func() (int64, error) {
					return 1, failedErr
				},
			},
			expected: failedErr,
		},
		{
			f: &fakeSqlResult{
				RowsAffectedFn: func() (int64, error) {
					return 0, nil
				},
			},
			expected: core.ErrConflictNewerVersion,
		},
	}

	repository := deploymentRepository{}
	for idx, test := range tt {
		result := repository.handleUpdateStatusResult(test.f)
		assert.Equal(t, test.expected, result, "test %d", idx)
	}
}
