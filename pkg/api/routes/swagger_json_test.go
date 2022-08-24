package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dcaf-labs/drip/pkg/configs"

	"github.com/labstack/echo/v4"

	"github.com/test-go/testify/assert"
)

func TestHandler_GetSwaggerJson(t *testing.T) {

	t.Run("getURL should return correct api URL", func(t *testing.T) {
		assert.True(t, strings.Contains(getURL(configs.NilEnv, 0), "localhost"))
		assert.True(t, strings.Contains(getURL(configs.LocalnetEnv, 0), "localhost"))
		assert.True(t, strings.Contains(getURL(configs.DevnetEnv, 0), "devnet"))
		assert.True(t, strings.Contains(getURL(configs.MainnetEnv, 0), "mainnet"))
	})

	t.Run("should return json from handler", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h := Handler{}
		assert.NoError(t, h.GetSwaggerJson(c))
		assert.Equal(t, rec.Code, 200)
		assert.NotEmpty(t, rec.Body.String())
	})
}
