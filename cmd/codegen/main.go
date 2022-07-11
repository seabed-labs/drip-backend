package main

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/configs"
	psql2 "github.com/dcaf-labs/drip/pkg/database/psql"

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
			configs.NewPSQLConfig,
			psql2.NewDatabase,
			psql2.NewGORMDatabase,
		),
		fx.Invoke(
			psql2.RunMigrations,
			psql2.GenerateModels,
		),
		fx.NopLogger,
	}
}
