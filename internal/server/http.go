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
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: e,
	}
	log.WithField("port", config.Port).Infof("starting server")
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
