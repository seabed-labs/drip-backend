package main

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/repository/database"

	"github.com/dcaf-labs/drip/pkg/service/configs"

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
			database.NewDatabase,
			database.NewGORMDatabase,
		),
		fx.Invoke(
			database.RunMigrations,
			database.GenerateModels,
		),
		fx.NopLogger,
	}
}
