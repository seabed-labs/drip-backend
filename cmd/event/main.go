package main

import (
	"context"

	"github.com/dcaf-protocol/drip/internal/pkg/processor"

	"github.com/dcaf-protocol/drip/internal/event"
	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/internal/pkg/repository"

	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/dcaf-protocol/drip/internal/database/psql"
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
			psql.NewGORMDatabase,
			repository.Use,
			solana.NewSolanaClient,
			processor.NewProcessor,
		),
		fx.Invoke(
			event.NewDripProgramProcessor,
		),
		fx.NopLogger,
	}
}
