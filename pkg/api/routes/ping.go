package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	Swagger "github.com/dcaf-labs/drip/pkg/swagger"
)

func (h Handler) Get(
	c echo.Context,
) error {
	return c.JSON(http.StatusOK, Swagger.PingResponse{
		Message: "pong",
	})
}
