package database

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/config"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(
	lifecycle fx.Lifecycle,
	dbConfig config.PSQLConfig,
) (*sqlx.DB, error) {
	onStop, err := maybeSetupEmbeddedDB(dbConfig)
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

func maybeSetupEmbeddedDB(dbConfig config.PSQLConfig) (func() error, error) {
	if !dbConfig.GetShouldUseEmbeddedDB() {
		return nil, nil
	}
	embeddedDB := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Username(dbConfig.GetUser()).
		Password(dbConfig.GetPassword()).
		Database(dbConfig.GetDBName()).
		Port(uint32(dbConfig.GetPort())),
	)
	if err := embeddedDB.Start(); err != nil {
		return nil, err
	}
	return embeddedDB.Stop, nil
}
