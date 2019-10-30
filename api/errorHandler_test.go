package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
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

	err := echo.NewHTTPError(http.StatusBadRequest, "test")

	ErrorHandler(err, ctx)

	assert.Empty(t, logBuf)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "{\"message\":\"test\"}\n", rec.Body.String())
}

func Test_ErrorHandler_WhenHttpErrorInternal_AlwaysLogsInternal(t *testing.T) {
	logBuf := &bytes.Buffer{}
	ctx, rec := errorHandlerTestSetup(logBuf)

	err := echo.NewHTTPError(http.StatusBadRequest, "test")
	err = err.SetInternal(errors.New("test"))

	ErrorHandler(err, ctx)

	splitLogLines := strings.Split(logBuf.String(), "\n")
	// There's an extra \n that creates an extra line
	assert.Len(t, splitLogLines, 2)
	jsonLog := unmarshalErrorLogLine(t, splitLogLines[0])
	assert.Equal(t, "ERROR", jsonLog["level"])
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "{\"message\":\"test\"}\n", rec.Body.String())
}

func Test_ErrorHandler_WhenValidationError_ReturnsValidationResponse(t *testing.T) {
	logBuf := &bytes.Buffer{}
	ctx, rec := errorHandlerTestSetup(logBuf)

	err := &core.ValidationError{
		Message: "validation failed",
		ValidationErrors: validation.Errors{
			"field": errors.New("field error"),
		},
	}

	ErrorHandler(err, ctx)

	assert.Empty(t, logBuf)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	jsonResponse := ValidationErrorResponse{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &jsonResponse), rec.Body.String())
	assert.Equal(t, jsonResponse.Message, "validation failed")
	assert.Len(t, jsonResponse.ValidationErrors, 1)
	assert.Equal(t, "field error", jsonResponse.ValidationErrors["field"])
}

func Test_ErrorHandler_WhenValidationErrorWithInternal_LogsInternal(t *testing.T) {
	logBuf := &bytes.Buffer{}
	ctx, _ := errorHandlerTestSetup(logBuf)

	err := &core.ValidationError{
		Message:  "validation failed",
		Internal: errors.New("test error"),
	}

	ErrorHandler(err, ctx)

	splitLogLines := strings.Split(logBuf.String(), "\n")
	// There's an extra \n that creates an extra line
	assert.Len(t, splitLogLines, 2)
	jsonLog := unmarshalErrorLogLine(t, splitLogLines[0])
	assert.Equal(t, "ERROR", jsonLog["level"])
	assert.Contains(t, jsonLog["message"], "test error")
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
