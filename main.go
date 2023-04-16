package main

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/ixparser"

	"github.com/dcaf-labs/drip/pkg/producer"

	api2 "github.com/dcaf-labs/solana-go-retryable-http-client"

	"github.com/dcaf-labs/drip/pkg/api"
	"github.com/dcaf-labs/drip/pkg/consumer"
	"github.com/dcaf-labs/drip/pkg/job/token"
	"github.com/dcaf-labs/drip/pkg/job/tokenaccount"

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
	_ "github.com/heroku/x/hmetrics/onload"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	fxApp := fx.New(getDependencies()...)
	if err := fxApp.Start(context.Background()); err != nil {
		log.WithError(err).Fatalf("failed to start drip backend")
	}
	log.Info("starting drip backend")
	sig := <-fxApp.Done()
	log.WithField("signal", sig).Infof("received exit signal, stoping api")
}

func getDependencies() []fx.Option {
	return []fx.Option{
		fx.Provide(
			config.NewAppConfig,
			config.NewPSQLConfig,
			database.NewDatabase,
			database.NewGORMDatabase,
			query.Use,
			repository.NewRepository,
			repository.NewTransactionProcessingCheckpointRepository,
			repository.NewTransactionUpdateQueue,
			repository.NewAccountUpdateQueue,
			api2.GetDefaultClientProvider,
			solana.NewSolanaClient,
			tokenregistry.NewTokenRegistry,
			orcawhirlpool.NewOrcaWhirlpoolClient,
			coingecko.NewCoinGeckoClient,
			processor.NewProcessor,
			alert.NewAlertService,
			base.NewBase,
			middleware.NewHandler,
			controller.NewHandler,
			ixparser.NewIxParser,
		),
		fx.Invoke(
			database.RunMigrations,
			producer.Server,
			consumer.Server,
			api.StartServer,
			token.NewTokenJob,
			tokenaccount.NewTokenAccountJob,
		),
		fx.NopLogger,
	}
}
