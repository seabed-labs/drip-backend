package main

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/api"
	"github.com/dcaf-labs/drip/pkg/api/middleware"
	controller "github.com/dcaf-labs/drip/pkg/api/routes"
	"github.com/dcaf-labs/drip/pkg/service/base"

	"github.com/dcaf-labs/drip/pkg/service/clients"
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
		log.WithError(err).Fatalf("failed to start drip event processor")
	}
	log.Info("starting drip event processor")
	sig := <-fxApp.Done()
	log.WithFields(log.Fields{"signal": sig}).
		Infof("received exit signal, stoping event processor")
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
			clients.DefaultClientProvider,
			solana.NewSolanaClient,
			base.NewBase,
			middleware.NewHandler,
			controller.NewHandler,
		),
		fx.Invoke(
			api.StartServer,
		),
		fx.NopLogger,
	}
}
