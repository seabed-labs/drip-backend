package main

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/api"
	"github.com/dcaf-labs/drip/pkg/api/middleware"
	controller "github.com/dcaf-labs/drip/pkg/api/routes"
	"github.com/dcaf-labs/drip/pkg/clients/solana"
	"github.com/dcaf-labs/drip/pkg/configs"
	"github.com/dcaf-labs/drip/pkg/database/psql"
	"github.com/dcaf-labs/drip/pkg/repository"
	"github.com/dcaf-labs/drip/pkg/repository/query"
	"github.com/dcaf-labs/drip/pkg/service/processor"
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
			psql.NewDatabase,
			psql.NewGORMDatabase,
			query.Use,
			solana.NewSolanaClient,
			repository.NewRepository,
			middleware.NewHandler,
			controller.NewHandler,
			processor.NewProcessor,
		),
		fx.Invoke(
			// func() { log.SetFormatter(&log.JSONFormatter{}) },
			psql.RunMigrations,
			api.StartServer,
		),
		fx.NopLogger,
	}
}
