package test

import (
	"context"
	"os"

	"github.com/dcaf-labs/drip/pkg/service/repository/database"

	"github.com/dcaf-labs/drip/pkg/service/configs"

	"github.com/dcaf-labs/drip/pkg/service/clients/solana"

	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/query"

	controller "github.com/dcaf-labs/drip/pkg/api/routes"
	"github.com/dcaf-labs/drip/pkg/service/processor"

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
			database.NewDatabase,
			database.NewGORMDatabase,
			query.Use,
			solana.NewSolanaClient,
			repository.NewRepository,
			controller.NewHandler,
			processor.NewProcessor,
		),
		fx.Invoke(
			database.RunMigrations,
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
