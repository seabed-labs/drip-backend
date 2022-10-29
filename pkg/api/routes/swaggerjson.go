package controller

import (
	"net/http"

	"github.com/dcaf-labs/drip/pkg/api/apispec"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

func (h Handler) GetSwaggerJson(
	c echo.Context,
) error {
	swaggerSpec, err := apispec.GetSwagger()
	if err != nil {
		return err
	}
	swaggerSpec.Servers = openapi3.Servers{
		&openapi3.Server{URL: getServerURL(h.network, h.env, h.port)},
	}
	return c.JSON(http.StatusOK, swaggerSpec)
}
