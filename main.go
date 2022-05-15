package main

import (
	"context"

	"github.com/dcaf-protocol/drip/internal/scripts"

	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"

	"github.com/dcaf-protocol/drip/internal/pkg/api"
	"github.com/dcaf-protocol/drip/internal/server"

	"github.com/dcaf-protocol/drip/internal/pkg/repository"

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
		Infof("received exit signal, stoping server")
}

func getDependencies() []fx.Option {
	return []fx.Option{
		fx.Provide(
			configs.NewAppConfig,
			configs.NewPSQLConfig,
			psql.NewGORMDatabase,
			repository.Use,
			solana.CreateSolanaClient,
			api.NewHandler,
		),
		fx.Invoke(
			psql.RunMigrations,
			scripts.Backfill,
			// func() { log.SetFormatter(&log.JSONFormatter{}) },
			server.Run,
		),
		fx.NopLogger,
	}
}
