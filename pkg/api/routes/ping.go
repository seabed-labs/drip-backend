package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/labstack/echo/v4"
)

func (h Handler) Get(
	c echo.Context,
) error {
	return c.JSON(http.StatusOK, apispec.PingResponse{
		Message: "pong",
	})
}
