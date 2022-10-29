package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dcaf-labs/drip/pkg/service/config"

	"github.com/labstack/echo/v4"

	"github.com/test-go/testify/assert"
)

func TestHandler_GetSwaggerJson(t *testing.T) {

	t.Run("getServerURL should return correct api URL", func(t *testing.T) {
		assert.True(t, strings.Contains(getServerURL(config.NilNetwork, config.StagingEnv, 0), "localhost"))
		assert.True(t, strings.Contains(getServerURL(config.LocalNetwork, config.StagingEnv, 0), "localhost"))
		assert.True(t, strings.Contains(getServerURL(config.DevnetNetwork, config.StagingEnv, 0), "devnet"))
		assert.True(t, strings.Contains(getServerURL(config.MainnetNetwork, config.StagingEnv, 0), "mainnet"))
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
