package main

import (
	"context"

	"github.com/dcaf-protocol/drip/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/pkg/configs"
	psql2 "github.com/dcaf-protocol/drip/pkg/database/psql"
	"github.com/dcaf-protocol/drip/pkg/processor"
	"github.com/dcaf-protocol/drip/pkg/repository"
	"github.com/dcaf-protocol/drip/pkg/repository/query"

	"github.com/dcaf-protocol/drip/internal/scripts"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func main() {
	fxApp := fx.New(getDependencies()...)
	if err := fxApp.Start(context.Background()); err != nil {
		log.WithError(err).Fatalf("failed to start codegen")
	}
	if err := fxApp.Start(context.Background()); err != nil {
		log.WithField("err", err.Error()).Fatalf("starting fx app for codegen")
	}
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
			processor.NewProcessor,
		),
		fx.Invoke(
			psql2.RunMigrations,
			scripts.Backfill,
		),
		fx.NopLogger,
	}
}
