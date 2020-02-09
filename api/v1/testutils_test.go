package v1

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

func safeMarshal(i interface{}) io.Reader {
	jsonBytes, _ := json.Marshal(i)
	return bytes.NewBuffer(jsonBytes)
}

func newContextWithRecorder(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}
