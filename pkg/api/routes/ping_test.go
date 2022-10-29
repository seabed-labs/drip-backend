package controller

import (
	"testing"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/dcaf-labs/drip/pkg/unittest"

	"github.com/labstack/echo/v4"

	"github.com/test-go/testify/assert"
)

func TestHandler_Get(t *testing.T) {
	e := echo.New()

	t.Run("should return pong", func(t *testing.T) {
		c, rec := unittest.GetTestRequestRecorder(e, nil)

		h := Handler{}

		assert.NoError(t, h.Get(c))
		assert.Equal(t, rec.Code, 200)
		pingRes, err := apispec.ParseGetResponse(rec.Result())
		assert.NoError(t, err)
		assert.NotNil(t, pingRes.JSON200)
		assert.Equal(t, pingRes.JSON200.Message, "pong")
	})
}
