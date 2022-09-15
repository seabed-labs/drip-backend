package controller

import (
	"fmt"
	"net/http"

	"github.com/dcaf-labs/drip/pkg/configs"

	"github.com/labstack/echo/v4"

	Swagger "github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/getkin/kin-openapi/openapi3"
)

func (h Handler) GetSwaggerJson(
	c echo.Context,
) error {
	swaggerSpec, err := Swagger.GetSwagger()
	if err != nil {
		return err
	}
	swaggerSpec.Servers = openapi3.Servers{
		&openapi3.Server{URL: getURL(h.network, h.env, h.port)},
	}
	return c.JSON(http.StatusOK, swaggerSpec)
}

func getURL(network configs.Network, env configs.Environment, port int) string {
	if configs.IsMainnet(network) {
		return "drip-backend-mainnet.herokuapp.com"
	} else if configs.IsDevnet(network) {
		if configs.IsStaging(env) {
			return "drip-backend-devnet-staging.herokuapp.com"
		} else if configs.IsProd(env) {
			return "drip-backend-devnet.herokuapp.com"
		}
	}
	return fmt.Sprintf("http://localhost:%d", port)
}
