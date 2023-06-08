package main

import (
	"context"

	queuerepository "github.com/dcaf-labs/drip/pkg/service/repository/queue"

	"github.com/dcaf-labs/drip/pkg/service/ixparser"

	api "github.com/dcaf-labs/solana-go-retryable-http-client"

	"github.com/dcaf-labs/drip/pkg/job/token"
	"github.com/dcaf-labs/drip/pkg/job/tokenaccount"
	"github.com/dcaf-labs/drip/pkg/service/clients/coingecko"

	"github.com/dcaf-labs/drip/pkg/consumer"
	"github.com/dcaf-labs/drip/pkg/service/alert"
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
			queuerepository.NewAccountUpdateQueue,
			queuerepository.NewTransactionUpdateQueue,
			api.GetDefaultClientProvider,
			solana.NewSolanaClient,
			tokenregistry.NewTokenRegistry,
			orcawhirlpool.NewOrcaWhirlpoolClient,
			coingecko.NewCoinGeckoClient,
			processor.NewProcessor,
			alert.NewAlertService,
			ixparser.NewIxParser,
		),
		fx.Invoke(
			database.RunMigrations,
			token.NewTokenJob,
			tokenaccount.NewTokenAccountJob,
			consumer.Server,
		),
		fx.NopLogger,
	}
}