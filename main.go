package main

import (
	"context"

	"github.com/dcaf-protocol/drip/internal/pkg/repository"

	"github.com/dcaf-protocol/drip/internal/pkg/repository/query"

	"github.com/dcaf-protocol/drip/internal/pkg/processor"

	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"

	"github.com/dcaf-protocol/drip/internal/api"
	"github.com/dcaf-protocol/drip/internal/pkg/controller"

	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/dcaf-protocol/drip/internal/database/psql"

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
			controller.NewHandler,
			processor.NewProcessor,
		),
		fx.Invoke(
			// func() { log.SetFormatter(&log.JSONFormatter{}) },
			psql.RunMigrations,
			api.APIServer,
		),
		fx.NopLogger,
	}
}
