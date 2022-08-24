package controller

import (
	"fmt"
	"net/http"

	"github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/dcaf-labs/drip/pkg/configs"
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
		&openapi3.Server{URL: getURL(h.env, h.port)},
	}
	return c.JSON(http.StatusOK, swaggerSpec)
}

func getURL(env configs.Environment, port int) string {
	switch env {
	case configs.DevnetEnv:
		return "drip-backend-devnet.herokuapp.com"
	case configs.MainnetEnv:
		return "drip-backend-mainnet.herokuapp.com"
	case configs.NilEnv:
		fallthrough
	case configs.LocalnetEnv:
		fallthrough
	default:
		return fmt.Sprintf("http://localhost:%d", port)
	}
}
