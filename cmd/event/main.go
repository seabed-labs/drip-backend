package main

import (
	"context"

	"github.com/dcaf-protocol/drip/internal/pkg/repository/query"

	"github.com/dcaf-protocol/drip/internal/pkg/repository"

	"github.com/dcaf-protocol/drip/internal/pkg/processor"

	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/dcaf-protocol/drip/internal/database/psql"
	"github.com/dcaf-protocol/drip/internal/event"
	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func main() {
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
			configs.NewAppConfig,
			configs.NewPSQLConfig,
			psql.NewDatabase,
			psql.NewGORMDatabase,
			query.Use,
			repository.NewRepository,
			solana.NewSolanaClient,
			processor.NewProcessor,
		),
		fx.Invoke(
			event.Server,
		),
		fx.NopLogger,
	}
}
