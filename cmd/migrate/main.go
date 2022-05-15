package main

import (
	"context"

	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/dcaf-protocol/drip/internal/database/psql"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func main() {
	fxApp := fx.New(getDependencies()...)
	if err := fxApp.Start(context.Background()); err != nil {
		log.WithError(err).Fatalf("failed to start migrations")
	}
	if err := fxApp.Start(context.Background()); err != nil {
		log.WithField("err", err.Error()).Fatalf("starting fx app for migrations")
	}
}

func getDependencies() []fx.Option {
	return []fx.Option{
		fx.Provide(
			configs.NewPSQLConfig,
			psql.NewDatabase,
		),
		fx.Invoke(
			psql.RunMigrations,
		),
		fx.NopLogger,
	}
}