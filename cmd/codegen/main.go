package main

import (
	"context"

	"github.com/dcaf-labs/drip/internal/codegen"
	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/repository/database"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	fxApp := fx.New(getDependencies()...)
	if err := fxApp.Start(context.Background()); err != nil {
		logrus.WithField("err", err.Error()).Fatalf("codegen failed")
	}
}

func getDependencies() []fx.Option {
	return []fx.Option{
		fx.Provide(
			config.NewPSQLConfig,
			database.NewDatabase,
			database.NewGORMDatabase,
		),
		fx.Invoke(
			database.RunMigrations,
			codegen.GenerateModels,
			codegen.GenerateAPIServer,
		),
		fx.NopLogger,
	}
}
