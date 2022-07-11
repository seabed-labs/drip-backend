package main

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/processor"

	"github.com/dcaf-labs/drip/pkg/clients/solana"
	"github.com/dcaf-labs/drip/pkg/configs"
	"github.com/dcaf-labs/drip/pkg/database/psql"
	"github.com/dcaf-labs/drip/pkg/event"
	"github.com/dcaf-labs/drip/pkg/repository"
	"github.com/dcaf-labs/drip/pkg/repository/query"

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
