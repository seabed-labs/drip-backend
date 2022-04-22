package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dcaf-protocol/drip/internal/api"
	"github.com/dcaf-protocol/drip/internal/configs"
	swagger "github.com/dcaf-protocol/drip/pkg/swagger"
	oapiMiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func Run(
	lc fx.Lifecycle,
	api *api.Handler,
	config *configs.Config,
) {
	var httpSrv *http.Server
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error
			httpSrv, err = listenAndServe(api, config)
			return err
		},
		OnStop: func(ctx context.Context) error {
			return shutdown(httpSrv)
		},
	})
}

func listenAndServe(
	handler *api.Handler,
	config *configs.Config,
) (*http.Server, error) {
	swaggerSpec, err := swagger.GetSwagger()
	if err != nil {
		return nil, err
	}
	swaggerSpec.Servers = nil

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(oapiMiddleware.OapiRequestValidator(swaggerSpec))

	swagger.RegisterHandlers(e, handler)
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	//r.Use(middleware.OapiRequestValidator(swaggerSpec))
	//swagger.RegisterHandlers()
	// We now register our petStore above as the handler for the interface
	//return nil, nil
	//swagger, err := api.
	//docs.SwaggerInfo.BasePath = "/"
	//docs.SwaggerInfo.Host = getURL(config)
	//r.Use(loggingMiddleware())
	//r.GET("/", api.Ping)
	//r.GET("/ping", api.Ping)
	//if !configs.IsProd(config.Environment) {
	//	r.GET("/mint", api.Mint)
	//}
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: e,
	}
	//log.WithField("port", config.Port).Infof("starting server")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithField("err", err.Error()).Fatalf("server listening")
		}
	}()
	return srv, nil
}

func shutdown(
	httpSrv *http.Server,
) error {
	log.Infof("stopping server")
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := httpSrv.Shutdown(ctxShutDown); err != nil {
		log.WithField("err", err.Error()).Fatalf("failed to shutdown server")
		return err
	}

	log.Infof("server exited")
	return nil
}

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.WithField("params", c.Params).Infof("handling %s", c.Request.RequestURI)
		c.Next()
	}
}

func getURL(config *configs.Config) string {
	switch config.Environment {
	case configs.DevnetEnv:
		return "drip-backend-devnet.herokuapp.com"
	case configs.MainnetEnv:
		return "drip-backend-mainnet.herokuapp.com"
	case configs.NilEnv:
		fallthrough
	case configs.LocalnetEnv:
		fallthrough
	default:
		return fmt.Sprintf("localhost:%d", config.Port)
	}
}
