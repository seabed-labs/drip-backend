package main

import (
	"context"

	repository2 "github.com/dcaf-labs/drip/pkg/service/repository/analytics"

	"github.com/dcaf-labs/drip/pkg/service/ixparser"

	api2 "github.com/dcaf-labs/solana-go-retryable-http-client"

	"github.com/dcaf-labs/drip/pkg/api"
	"github.com/dcaf-labs/drip/pkg/api/middleware"
	controller "github.com/dcaf-labs/drip/pkg/api/routes"
	"github.com/dcaf-labs/drip/pkg/service/base"

	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/database"
	"github.com/dcaf-labs/drip/pkg/service/repository/query"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	fxApp := fx.New(getDependencies()...)
	if err := fxApp.Start(context.Background()); err != nil {
		log.WithError(err).Fatalf("failed to start drip consumer processor")
	}
	log.Info("starting drip consumer processor")
	sig := <-fxApp.Done()
	log.WithFields(log.Fields{"signal": sig}).
		Infof("received exit signal, stoping consumer processor")
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
			repository2.NewAnalyticsRepository,
			api2.GetDefaultClientProvider,
			solana.NewSolanaClient,
			base.NewBase,
			middleware.NewHandler,
			controller.NewHandler,
			ixparser.NewIxParser,
		),
		fx.Invoke(
			api.StartServer,
		),
		fx.NopLogger,
	}
}
