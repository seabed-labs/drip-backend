package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/dcaf-protocol/drip/internal/pkg/controller"
	swagger "github.com/dcaf-protocol/drip/pkg/swagger"
	oapiMiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"google.golang.org/api/idtoken"
)

func APIServer(
	lc fx.Lifecycle,
	api *controller.Handler,
	config *configs.AppConfig,
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
	handler *controller.Handler,
	config *configs.AppConfig,
) (*http.Server, error) {
	swaggerSpec, err := swagger.GetSwagger()
	if err != nil {
		return nil, err
	}
	swaggerSpec.Servers = nil
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(validateAccessToken(config.GoogleClientID))
	e.Use(oapiMiddleware.OapiRequestValidator(swaggerSpec))
	swagger.RegisterHandlers(e, handler)
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

// TODO(Mocha): Move this to a middleware folder
func validateAccessToken(googleClientID string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !strings.Contains(c.Request().RequestURI, "admin") {
				return next(c)
			}
			accessToken := c.Request().Header.Get("token-id")
			payload, err := idtoken.Validate(context.Background(), accessToken, googleClientID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, swagger.ErrorResponse{Error: "invalid token-id"})
			}
			log.
				WithField("email", payload.Claims["email"]).
				WithField("name", payload.Claims["name"]).
				Info("authorized")
			return next(c)
		}
	}
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
