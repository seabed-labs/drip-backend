package api

import (
	"net/http"

	Swagger "github.com/dcaf-protocol/drip/pkg/swagger"

	"github.com/labstack/echo/v4"
)

func (h Handler) Get(
	c echo.Context,
) error {
	return c.JSON(http.StatusOK, Swagger.PingResponse{
		Message: "pong",
	})
}
