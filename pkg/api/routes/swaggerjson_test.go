package controller

import (
	"testing"

	"github.com/dcaf-labs/drip/pkg/unittest"

	"github.com/labstack/echo/v4"

	"github.com/test-go/testify/assert"
)

func TestHandler_GetSwaggerJson(t *testing.T) {
	e := echo.New()

	t.Run("should return json from handler", func(t *testing.T) {
		c, rec := unittest.GetTestRequestRecorder(e, nil)

		h := Handler{}

		assert.NoError(t, h.GetSwaggerJson(c))
		assert.Equal(t, rec.Code, 200)
		assert.NotEmpty(t, rec.Body.String())
	})
}
