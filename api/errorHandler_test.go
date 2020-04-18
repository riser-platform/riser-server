package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/riser-platform/riser-server/pkg/core"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
)

func Test_ErrorHandler_LogsWithStackTrace(t *testing.T) {
	logBuf := &bytes.Buffer{}
	ctx, rec := errorHandlerTestSetup(logBuf)

	err := errors.WithStack(errors.New("test"))

	ErrorHandler(err, ctx)

	splitLogLines := strings.Split(logBuf.String(), "\n")

	// There's an extra \n that creates an extra line
	assert.Len(t, splitLogLines, 2)
	jsonLog := unmarshalErrorLogLine(t, splitLogLines[0])
	assert.Equal(t, jsonLog["level"], "ERROR")
	assert.Contains(t, jsonLog["message"], "/api/errorHandler_test.go:")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "{\"message\":\"Internal Server Error\"}\n", rec.Body.String())
}

func Test_ErrorHandler_WhenHttpErrorStatusNon500_DoesNotLog(t *testing.T) {
	logBuf := &bytes.Buffer{}
	ctx, rec := errorHandlerTestSetup(logBuf)

	err := echo.NewHTTPError(http.StatusConflict, "test")

	ErrorHandler(err, ctx)

	assert.Empty(t, logBuf)
	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Equal(t, "{\"message\":\"test\"}\n", rec.Body.String())
}

func Test_ErrorHandler_WhenHttpErrorInternal_AlwaysLogsInternal(t *testing.T) {
	logBuf := &bytes.Buffer{}
	ctx, rec := errorHandlerTestSetup(logBuf)

	err := echo.NewHTTPError(http.StatusConflict, "test")
	err = err.SetInternal(errors.New("log me"))

	ErrorHandler(err, ctx)

	splitLogLines := strings.Split(logBuf.String(), "\n")
	// There's an extra \n that creates an extra line
	assert.Len(t, splitLogLines, 2)
	jsonLog := unmarshalErrorLogLine(t, splitLogLines[0])
	assert.Equal(t, "ERROR", jsonLog["level"])
	assert.Contains(t, jsonLog["message"], "log me")
	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Equal(t, "{\"message\":\"test\"}\n", rec.Body.String())
}

func Test_ErrorHandler_WhenValidationErrorWithFields_FormatsResponse(t *testing.T) {
	logBuf := &bytes.Buffer{}
	ctx, rec := errorHandlerTestSetup(logBuf)

	err := core.NewValidationError("validation failed",
		validation.Errors{
			"field1": errors.New("field1 error"),
			"field2": errors.New("field2 error"),
			"field3": validation.Errors{
				"nested": errors.New("nested1 error"),
			},
		},
	)

	ErrorHandler(err, ctx)

	assert.Empty(t, logBuf)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	jsonResponse := ValidationErrorResponse{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &jsonResponse), rec.Body.String())
	assert.Equal(t, jsonResponse.Message, "validation failed")
	assert.Len(t, jsonResponse.ValidationErrors, 3)
	assert.Equal(t, "field1 error", jsonResponse.ValidationErrors["field1"])
	assert.Equal(t, "field2 error", jsonResponse.ValidationErrors["field2"])
	assert.Equal(t, "nested1 error", jsonResponse.ValidationErrors["field3.nested"])
}

func Test_ErrorHandler_WhenValidationErrorNoFields(t *testing.T) {
	logBuf := &bytes.Buffer{}
	ctx, rec := errorHandlerTestSetup(logBuf)

	err := core.NewValidationError("validation failed", errors.New("details"))

	ErrorHandler(err, ctx)

	assert.Empty(t, logBuf)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	jsonResponse := ValidationErrorResponse{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &jsonResponse), rec.Body.String())
	assert.Equal(t, jsonResponse.Message, "validation failed: details")
	assert.Len(t, jsonResponse.ValidationErrors, 0)
}

func Test_ErrorHandler_WhenValidationErrorIsNil(t *testing.T) {
	logBuf := &bytes.Buffer{}
	ctx, rec := errorHandlerTestSetup(logBuf)

	err := core.NewValidationError("validation failed", nil)

	ErrorHandler(err, ctx)

	assert.Empty(t, logBuf)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	jsonResponse := ValidationErrorResponse{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &jsonResponse), rec.Body.String())
	assert.Equal(t, jsonResponse.Message, "validation failed")
	assert.Len(t, jsonResponse.ValidationErrors, 0)
}

// An ozzo-validation Internal error means that something went wrong (e.g. a misconfigured validation rule).
func Test_ErrorHandler_WhenValidationErrorIsInternal_Returns500AndLogsInternal(t *testing.T) {
	logBuf := &bytes.Buffer{}
	ctx, rec := errorHandlerTestSetup(logBuf)

	err := &core.ValidationError{
		Message:         "validation failed",
		ValidationError: validation.NewInternalError(errors.New("log me")),
	}

	ErrorHandler(err, ctx)

	splitLogLines := strings.Split(logBuf.String(), "\n")
	// There's an extra \n that creates an extra line
	assert.Len(t, splitLogLines, 2)
	jsonLog := unmarshalErrorLogLine(t, splitLogLines[0])
	assert.Equal(t, "ERROR", jsonLog["level"])
	assert.Contains(t, jsonLog["message"], "log me")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	jsonResponse := ValidationErrorResponse{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &jsonResponse), rec.Body.String())
	assert.Equal(t, jsonResponse.Message, http.StatusText(http.StatusInternalServerError))
}

func errorHandlerTestSetup(logWriter io.Writer) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Logger.SetOutput(logWriter)
	e.HTTPErrorHandler = ErrorHandler
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func unmarshalErrorLogLine(t *testing.T, line string) map[string]interface{} {
	jsonLog := map[string]interface{}{}
	require.NoError(t, json.Unmarshal([]byte(line), &jsonLog), line)
	return jsonLog
}
