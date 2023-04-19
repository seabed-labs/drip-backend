package main

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/ixparser"

	api "github.com/dcaf-labs/solana-go-retryable-http-client"

	"github.com/dcaf-labs/drip/internal/cli"
	"github.com/dcaf-labs/drip/pkg/service/alert"
	"github.com/dcaf-labs/drip/pkg/service/clients/coingecko"
	"github.com/dcaf-labs/drip/pkg/service/clients/orcawhirlpool"
	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/clients/tokenregistry"
	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/database"
	"github.com/dcaf-labs/drip/pkg/service/repository/query"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func main() {
	fxApp := fx.New(getDependencies()...)
	if err := fxApp.Start(context.Background()); err != nil {
		log.WithField("err", err.Error()).Fatalf("cli error")
	}
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
			repository.NewAccountUpdateQueue,
			repository.NewTransactionUpdateQueue,

			api.GetDefaultClientProvider,
			solana.NewSolanaClient,
			tokenregistry.NewTokenRegistry,
			orcawhirlpool.NewOrcaWhirlpoolClient,
			coingecko.NewCoinGeckoClient,
			ixparser.NewIxParser,

			processor.NewProcessor,
			alert.NewAlertService,
		),
		fx.Invoke(
			database.RunMigrations,
			cli.RunCLI,
		),
		fx.NopLogger,
	}
}
