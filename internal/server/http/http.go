package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dcaf-protocol/drip/internal/api"
	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

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
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/", http.HandlerFunc(api.Ping)).Methods("GET")
	r.HandleFunc("/mint/{publicKey}", http.HandlerFunc(api.GetMint)).Methods("GET")
	// mux := http.NewServeMux()
	// mux.Handle("/", logMiddleware(http.HandlerFunc(api.Ping)))
	// mux.Handle("/mint/:address", logMiddleware(http.HandlerFunc(api.GetMint)))
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

type MiddlewareFunc func(http.Handler) http.Handler

func loggingMiddleware(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithField("params", r.URL.Query()).Infof("handling %s", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
