package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	middleware2 "github.com/dcaf-protocol/drip/pkg/api/middleware"
	controller "github.com/dcaf-protocol/drip/pkg/api/routes"

	"github.com/dcaf-protocol/drip/pkg/configs"
	swagger "github.com/dcaf-protocol/drip/pkg/swagger"
	oapiMiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func APIServer(
	lc fx.Lifecycle,
	config *configs.AppConfig,
	middlewareHandler *middleware2.Handler,
	apiHandler *controller.Handler,
) {
	var httpSrv *http.Server
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error
			httpSrv, err = listenAndServe(config, middlewareHandler, apiHandler)
			return err
		},
		OnStop: func(ctx context.Context) error {
			return shutdown(httpSrv)
		},
	})
}

func listenAndServe(
	config *configs.AppConfig,
	middlewareHandler *middleware2.Handler,
	apiHandler *controller.Handler,
) (*http.Server, error) {
	swaggerSpec, err := swagger.GetSwagger()
	if err != nil {
		return nil, err
	}
	swaggerSpec.Servers = nil
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middlewareHandler.ValidateAccessToken)
	e.Use(middlewareHandler.RateLimit)
	e.Use(oapiMiddleware.OapiRequestValidator(swaggerSpec))
	swagger.RegisterHandlers(e, apiHandler)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: cors.AllowAll().Handler(e),
	}
	log.WithField("port", config.Port).Infof("starting api")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithField("err", err.Error()).Fatalf("api listening")
		}
	}()
	return srv, nil
}

func shutdown(
	httpSrv *http.Server,
) error {
	log.Infof("stopping api")
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := httpSrv.Shutdown(ctxShutDown); err != nil {
		log.WithField("err", err.Error()).Fatalf("failed to shutdown api")
		return err
	}

	log.Infof("api exited")
	return nil
}
