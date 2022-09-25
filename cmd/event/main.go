package main

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/clients/tokenregistry"

	"github.com/dcaf-labs/drip/pkg/service/repository/database"

	"github.com/dcaf-labs/drip/pkg/service/configs"

	"github.com/dcaf-labs/drip/pkg/service/clients/solana"

	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/query"

	"github.com/dcaf-labs/drip/pkg/service/alert"
	"github.com/dcaf-labs/drip/pkg/service/processor"

	"github.com/dcaf-labs/drip/pkg/event"
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
			database.NewDatabase,
			database.NewGORMDatabase,
			query.Use,
			repository.NewRepository,
			solana.NewSolanaClient,
			tokenregistry.NewTokenRegistry,
			processor.NewProcessor,
			alert.NewService,
		),
		fx.Invoke(
			event.Server,
		),
		fx.NopLogger,
	}
}
