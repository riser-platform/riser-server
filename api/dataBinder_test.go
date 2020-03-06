package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testValidationError = errors.New("test validation error")

type decoratedModel struct {
	val           int
	bindVal       int
	validationVal int
}

func (d *decoratedModel) ApplyDefaults() error {
	d.val = d.bindVal
	return nil
}

func (d *decoratedModel) Validate() error {
	if d.val != d.validationVal {
		return testValidationError
	}
	return nil
}

type plainModel struct {
	val int
}

func Test_Bind_MutatesFieldsAndValidates(t *testing.T) {
	model := &decoratedModel{
		val:           1,
		bindVal:       2,
		validationVal: 2,
	}

	ctx := setupDataBinderTest(model)

	err := ctx.Bind(model)

	assert.NoError(t, err)
	assert.Equal(t, 2, model.val)
}

func Test_Bind_WhenValidationFails(t *testing.T) {
	model := &decoratedModel{
		val:           1,
		bindVal:       2,
		validationVal: 1,
	}

	ctx := setupDataBinderTest(model)

	err := ctx.Bind(model)

	require.IsType(t, &core.ValidationError{}, err)
	cve := err.(*core.ValidationError)
	assert.Equal(t, testValidationError, cve.ValidationError)
}

func Test_Bind_NoBindingOrValidation(t *testing.T) {
	model := &plainModel{
		val: 1,
	}

	ctx := setupDataBinderTest(model)

	err := ctx.Bind(model)

	assert.NoError(t, err)
	assert.Equal(t, 1, model.val)
}

func setupDataBinderTest(model interface{}) echo.Context {
	e := echo.New()
	e.Binder = &DataBinder{}
	req := httptest.NewRequest(http.MethodPost, "/", safeMarshal(model))
	req.Header.Add("CONTENT-TYPE", "application/json")

	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

func safeMarshal(i interface{}) io.Reader {
	jsonBytes, _ := json.Marshal(i)
	return bytes.NewBuffer(jsonBytes)
}
