package main

import (
	"context"

	"github.com/dcaf-protocol/drip/pkg/middleware"

	"github.com/dcaf-protocol/drip/pkg/api"
	"github.com/dcaf-protocol/drip/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/pkg/configs"
	"github.com/dcaf-protocol/drip/pkg/controller"
	psql2 "github.com/dcaf-protocol/drip/pkg/database/psql"
	"github.com/dcaf-protocol/drip/pkg/processor"
	"github.com/dcaf-protocol/drip/pkg/repository"
	"github.com/dcaf-protocol/drip/pkg/repository/query"

	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func main() {
	fxApp := fx.New(getDependencies()...)
	if err := fxApp.Start(context.Background()); err != nil {
		log.WithError(err).Fatalf("failed to start drip backend")
	}
	log.Info("starting drip backend")
	sig := <-fxApp.Done()
	log.WithFields(log.Fields{"signal": sig}).
		Infof("received exit signal, stoping api")
}

func getDependencies() []fx.Option {
	return []fx.Option{
		fx.Provide(
			configs.NewAppConfig,
			configs.NewPSQLConfig,
			psql2.NewDatabase,
			psql2.NewGORMDatabase,
			query.Use,
			solana.NewSolanaClient,
			repository.NewRepository,
			middleware.NewHandler,
			controller.NewHandler,
			processor.NewProcessor,
		),
		fx.Invoke(
			// func() { log.SetFormatter(&log.JSONFormatter{}) },
			psql2.RunMigrations,
			api.APIServer,
		),
		fx.NopLogger,
	}
}
