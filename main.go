package main

import (
	"context"

	"github.com/dcaf-protocol/drip/internal/api"
	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/dcaf-protocol/drip/internal/pkg/solanaclient"
	"github.com/dcaf-protocol/drip/internal/server/http"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	fxApp := fx.New(
		fx.Provide(
			configs.NewConfig,
			solanaclient.NewSolanaClient,
			api.NewHandler,
		),
		fx.Invoke(
			func() { log.SetFormatter(&log.JSONFormatter{}) },
			http.Run,
		),
		fx.NopLogger,
	)
	if err := fxApp.Start(context.Background()); err != nil {
		log.WithError(err).Fatalf("failed to start drip backend")
	}
	log.Info("starting drip backend")
	sig := <-fxApp.Done()
	log.WithFields(log.Fields{"signal": sig}).
		Infof("received exit signal, stoping server")

}
