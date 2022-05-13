package main

import (
	"context"
	"os/exec"

	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/dcaf-protocol/drip/internal/database/psql"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func main() {
	fxApp := fx.New(getDependencies()...)
	if err := fxApp.Start(context.Background()); err != nil {
		log.WithError(err).Fatalf("failed to start codegen")
	}
	if err := fxApp.Start(context.Background()); err != nil {
		log.WithField("err", err.Error()).Fatalf("starting fx app for codegen")
	}
}

func getDependencies() []fx.Option {
	return []fx.Option{
		fx.Provide(
			configs.NewPSQLConfig,
			psql.NewDatabase,
		),
		fx.Invoke(
			psql.RunMigrations,
			psql.GenerateModels,
		),
		fx.NopLogger,
	}
}

func codegen() {
	configs.LoadEnv()
	cmd := exec.Command("sqlboiler", "psql")
	stdout, err := cmd.Output()
	if err != nil {
		logrus.WithError(err).Fatalf("failed to run sqlboiler")
		return
	}
	logrus.Infof(string(stdout))
}
