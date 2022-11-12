package main

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/repository/database"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	fxApp := fx.New(getDependencies()...)
	if err := fxApp.Start(context.Background()); err != nil {
		logrus.WithError(err).Fatalf("failed to start migrations")
	}
	if err := fxApp.Start(context.Background()); err != nil {
		logrus.WithField("err", err.Error()).Fatalf("starting fx app for migrations")
	}
}

func getDependencies() []fx.Option {
	return []fx.Option{
		fx.Provide(
			config.NewPSQLConfig,
			database.NewDatabase,
		),
		fx.Invoke(
			database.RunMigrations,
		),
		fx.NopLogger,
	}
}
