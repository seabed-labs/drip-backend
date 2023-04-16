package integrationtest

import (
	"context"
	"os"

	"github.com/dcaf-labs/drip/pkg/api/middleware"
	controller "github.com/dcaf-labs/drip/pkg/api/routes"
	"github.com/dcaf-labs/drip/pkg/service/alert"
	"github.com/dcaf-labs/drip/pkg/service/base"
	"github.com/dcaf-labs/drip/pkg/service/clients/coingecko"
	"github.com/dcaf-labs/drip/pkg/service/clients/orcawhirlpool"
	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/clients/tokenregistry"
	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/database"
	"github.com/dcaf-labs/drip/pkg/service/repository/query"
	"github.com/dcaf-labs/drip/pkg/unittest"
	api2 "github.com/dcaf-labs/solana-go-retryable-http-client"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type TestOptions struct {
	FixturePath string
	AppConfig   config.AppConfig
	PSQLConfig  config.PSQLConfig
}

func TestWithInjectedDependencies(
	testOptions *TestOptions,
	testCase interface{},
) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// only set if it is not explicitly set already
	if os.Getenv("SHOULD_USE_EMBEDDED_DB") == "" {
		if err := os.Setenv("SHOULD_USE_EMBEDDED_DB", "true"); err != nil {
			logrus.WithError(err).Error("could not set SHOULD_USE_EMBEDDED_DB env var")
			os.Exit(1)
		}
	}

	// test http recorder
	httpClientProvider := api2.GetDefaultClientProvider()
	if testOptions != nil {
		recorderProvider, recorderTeardown := unittest.GetHTTPRecorderClientProvider("./fixtures/drip_1")
		defer recorderTeardown()
		httpClientProvider = func(options api2.RateLimitHTTPClientOptions) api2.RetryableHTTPClient {
			return recorderProvider()(options)
		}
	}
	providers := []interface{}{
		// Data access
		database.NewDatabase,
		database.NewGORMDatabase,
		query.Use,
		repository.NewRepository,
		repository.NewAccountUpdateQueue,
		// API Clients
		httpClientProvider,
		solana.NewSolanaClient,
		tokenregistry.NewTokenRegistry,
		orcawhirlpool.NewOrcaWhirlpoolClient,
		coingecko.NewCoinGeckoClient,
		// services
		base.NewBase,
		processor.NewProcessor,
		alert.NewAlertService,
		// server
		middleware.NewHandler,
		controller.NewHandler,
	}

	if testOptions != nil && testOptions.AppConfig != nil {
		providers = append(providers, func() config.AppConfig {
			return testOptions.AppConfig
		})
	} else {
		providers = append(providers, config.NewAppConfig)
	}
	if testOptions != nil && testOptions.PSQLConfig != nil {
		providers = append(providers, func() config.PSQLConfig {
			return testOptions.PSQLConfig
		})
	} else {
		providers = append(providers, config.NewPSQLConfig)
	}
	// comment out below for logs
	//logrus.SetOutput(ioutil.Discard)
	opts := []fx.Option{
		fx.Provide(providers...),
		fx.Invoke(
			database.RunMigrations,
			testCase,
		),
		fx.NopLogger,
	}
	app := fx.New(opts...)
	defer func() {
		if err := app.Stop(context.Background()); err != nil {
			logrus.WithError(err).Errorf("failed to stop integrationtest app")
		}
	}()
	if err := app.Start(context.Background()); err != nil {
		logrus.WithError(err).Error("failed to run integration test")
		panic(err)
	}
}
