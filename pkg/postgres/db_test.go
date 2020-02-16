package postgres

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_AddAuthToConnString(t *testing.T) {
	result, err := AddAuthToConnString(
		"postgres://myhost.local/riserdb?arg=val",
		"myuser",
		"mypass")

	assert.NoError(t, err)
	assert.Equal(t, "postgres://myuser:mypass@myhost.local/riserdb?arg=val", result)
}

func Test_AddAuthToConnString_BadUrl(t *testing.T) {
	result, err := AddAuthToConnString(
		"not@valid:test",
		"myuser",
		"mypass")

	assert.Empty(t, result)
	assert.Error(t, err)
}

type fakeSqlResult struct {
	LastInsertIdFn func() (int64, error)
	RowsAffectedFn func() (int64, error)
}

func (f *fakeSqlResult) LastInsertId() (int64, error) {
	return f.LastInsertIdFn()
}

func (f *fakeSqlResult) RowsAffected() (int64, error) {
	return f.RowsAffectedFn()
}

func Test_ResultHasRows(t *testing.T) {
	tt := []struct {
		f *fakeSqlResult
		e bool
	}{
		{
			f: &fakeSqlResult{
				RowsAffectedFn: func() (int64, error) {
					return 1, nil
				},
			},
			e: true,
		},
		{
			f: &fakeSqlResult{
				RowsAffectedFn: func() (int64, error) {
					return 0, nil
				},
			},
			e: false,
		},
		{
			f: &fakeSqlResult{
				RowsAffectedFn: func() (int64, error) {
					return 1, errors.New("err")
				},
			},
			e: false,
		},
	}

	for idx, test := range tt {
		assert.Equal(t, test.e, ResultHasRows(test.f), "test %d", idx)
	}

}
