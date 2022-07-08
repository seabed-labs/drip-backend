package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/test-go/testify/assert"
)

func TestHandler_Get(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := Handler{}
	assert.NoError(t, h.Get(c))
	assert.Equal(t, rec.Code, 200)
	assert.Equal(t, "{\"message\":\"pong\"}\n", rec.Body.String())
}
