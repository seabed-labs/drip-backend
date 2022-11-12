package database

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"

	"github.com/dcaf-labs/drip/pkg/service/config"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jmoiron/sqlx"
	"github.com/phayes/freeport"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(
	lifecycle fx.Lifecycle,
	dbConfig config.PSQLConfig,
) (*sqlx.DB, error) {
	onStop, err := maybeSetupEmbeddedDBWithRetry(dbConfig, 0, 3)
	if err != nil {
		return nil, err
	}
	lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return onStop()
		},
	})
	return sqlx.Connect("postgres", getConnectionString(dbConfig))
}

// NewGORMDatabase has a dummy import of  _ *sqlx.DB to ensure that NewDatabase is called first
func NewGORMDatabase(
	dbConfig config.PSQLConfig, _ *sqlx.DB,
) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(getConnectionString(dbConfig)), &gorm.Config{Logger: gormLogger{}})
}

func maybeSetupEmbeddedDBWithRetry(dbConfig config.PSQLConfig, try, maxTry int) (func() error, error) {
	if !dbConfig.GetShouldUseEmbeddedDB() {
		return nil, nil
	}
	if try == maxTry {
		return nil, fmt.Errorf("failed to get db conn with retry")
	}
	if onStop, err := maybeSetupEmbeddedDB(dbConfig); err != nil {
		logrus.WithError(err).Errorf("failed to get db conn, retrying...")
		return maybeSetupEmbeddedDBWithRetry(dbConfig, try+1, maxTry)
	} else {
		return onStop, err
	}
}

func maybeSetupEmbeddedDB(dbConfig config.PSQLConfig) (func() error, error) {
	if !dbConfig.GetShouldUseEmbeddedDB() {
		return nil, nil
	}
	root := config.GetProjectRoot()
	port, err := freeport.GetFreePort()
	if err != nil {
		return nil, err
	}
	dbConfig.SetPort(port)
	embeddedDB := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Username(dbConfig.GetUser()).
		Password(dbConfig.GetPassword()).
		Database(dbConfig.GetDBName()).
		Port(uint32(dbConfig.GetPort())).
		Version(embeddedpostgres.V14).
		RuntimePath(fmt.Sprintf("%s/.tmp/%s/extracted", root, dbConfig.GetDBName())).
		DataPath(fmt.Sprintf("%s/.tmp/%s/extracted/data", root, dbConfig.GetDBName())).
		BinariesPath(fmt.Sprintf("%s/.tmp/%s/extracted", root, dbConfig.GetDBName())).
		Logger(nil),
	)
	if err := embeddedDB.Start(); err != nil {
		return nil, err
	}
	return func() error {
		stopErr := embeddedDB.Stop()
		cleanErr := os.RemoveAll(fmt.Sprintf("%s/.tmp/%s/", root, dbConfig.GetDBName()))
		if stopErr != nil || cleanErr != nil {
			return fmt.Errorf("err in db cleanup %w", multierror.Append(stopErr, cleanErr))
		}
		return nil
	}, nil
}
