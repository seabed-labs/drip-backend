package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dcaf-protocol/drip/docs"
	"github.com/dcaf-protocol/drip/internal/api"
	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

// gin-swagger middleware
// swagger embed files

// @title           Drip Backend
// @version         1.0
// @description     Drip backend service.

// @contact.name   Dcaf Mocha
// @contact.email  dcafmocha@protonmail.com

// @host  localhost:8080
func Run(
	lc fx.Lifecycle,
	api *api.Handler,
	config *configs.Config,
) {
	var httpSrv http.Server
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			httpSrv = *listenAndServe(api, config)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdown(&httpSrv)
			return nil
		},
	})
}

func listenAndServe(
	api *api.Handler,
	config *configs.Config,
) *http.Server {

	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = getURL(config)
	r.Use(loggingMiddleware())
	r.GET("/", api.Ping)
	r.GET("/ping", api.Ping)
	if !configs.IsProd(config.Environment) {
		r.GET("/mint", api.Mint)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: r,
	}
	log.WithField("port", config.Port).Infof("starting server")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithField("err", err.Error()).Fatalf("server listening")
		}
	}()
	return srv
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
