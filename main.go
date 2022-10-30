package main

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/api"
	"github.com/dcaf-labs/drip/pkg/api/middleware"
	controller "github.com/dcaf-labs/drip/pkg/api/routes"
	"github.com/dcaf-labs/drip/pkg/event"
	"github.com/dcaf-labs/drip/pkg/service/alert"
	"github.com/dcaf-labs/drip/pkg/service/base"
	"github.com/dcaf-labs/drip/pkg/service/clients"
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
	fxApp := fx.New(getDependencies()...)
	if err := fxApp.Start(context.Background()); err != nil {
		log.WithError(err).Fatalf("failed to start drip backend")
	}
	log.Info("starting drip backend")
	sig := <-fxApp.Done()
	log.WithFields(log.Fields{"signal": sig}).
		Infof("received exit signal, stoping api")
}

func getDependencies() []fx.Option {
	//config, _ := config.NewAppConfig()
	// Hack to save on dyno costs, this will run  the event server and api server in the same dyno for staging
	//if config.IsStagingEnvironment(config.Environment) {
	return []fx.Option{
		fx.Provide(
			config.NewAppConfig,
			config.NewPSQLConfig,

			database.NewDatabase,
			database.NewGORMDatabase,
			query.Use,
			repository.NewRepository,
			repository.NewAccountUpdateQueue,

			clients.DefaultClientProvider,
			solana.NewSolanaClient,
			tokenregistry.NewTokenRegistry,
			orcawhirlpool.NewOrcaWhirlpoolClient,
			coingecko.NewCoinGeckoClient,

			processor.NewProcessor,
			alert.NewAlertService,
			base.NewBase,

			middleware.NewHandler,
			controller.NewHandler,
		),
		fx.Invoke(
			func() { log.SetFormatter(&log.JSONFormatter{}) },
			database.RunMigrations,
			api.StartServer,
			event.Server,
		),
		fx.NopLogger,
	}
	//} else {
	//	return []fx.Option{
	//		fx.Provide(
	//			config.NewAppConfig,
	//			config.NewPSQLConfig,
	//			database.NewDatabase,
	//			database.NewGORMDatabase,
	//			query.Use,
	//			solana.NewSolanaClient,
	//			repository.NewRepository,
	//			middleware.NewHandler,
	//			controller.NewHandler,
	//			processor.NewProcessor,
	//			base.NewBase,
	//		),
	//		fx.Invoke(
	//			// func() { log.SetFormatter(&log.JSONFormatter{}) },
	//			database.RunMigrations,
	//			api.StartServer,
	//		),
	//		fx.NopLogger,
	//	}
	//}
}
