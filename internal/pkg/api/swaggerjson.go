package api

import (
	"fmt"
	"net/http"

	"github.com/dcaf-protocol/drip/internal/configs"

	"github.com/labstack/echo/v4"

	Swagger "github.com/dcaf-protocol/drip/pkg/swagger"
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
