package controller

import (
	"fmt"
	"net/http"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	"github.com/dcaf-labs/drip/pkg/service/configs"

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
		&openapi3.Server{URL: getURL(h.network, h.env, h.port)},
	}
	return c.JSON(http.StatusOK, swaggerSpec)
}

func getURL(network configs.Network, env configs.Environment, port int) string {
	if configs.IsMainnetNetwork(network) {
		return "drip-backend-mainnet.herokuapp.com"
	} else if configs.IsDevnetNetwork(network) {
		if configs.IsStagingEnvironment(env) {
			return "drip-backend-devnet-staging.herokuapp.com"
		} else if configs.IsProductionEnvironment(env) {
			return "drip-backend-devnet.herokuapp.com"
		}
	}
	return fmt.Sprintf("http://localhost:%d", port)
}
