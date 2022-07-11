package test

import (
	"context"
	"os"

	controller "github.com/dcaf-protocol/drip/pkg/api/routes"
	"github.com/dcaf-protocol/drip/pkg/service/processor"

	"github.com/dcaf-protocol/drip/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/pkg/configs"
	"github.com/dcaf-protocol/drip/pkg/database/psql"
	"github.com/dcaf-protocol/drip/pkg/repository"
	"github.com/dcaf-protocol/drip/pkg/repository/query"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func InjectDependencies(
	testCase interface{},
) {
	err := os.Setenv("IS_TEST_DB", "true")
	if err != nil {
		panic("could not set IS_TEST_DB env var")
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	opts := []fx.Option{
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
			psql.RunMigrations,
			testCase,
		),
		fx.NopLogger,
	}
	app := fx.New(opts...)
	defer func() {
		if err := app.Stop(context.Background()); err != nil {
			logrus.WithError(err).Errorf("failed to stop test app")
		}
	}()
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}
