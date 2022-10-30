package database

import (
	"fmt"

	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(
	dbConfig config.PSQLConfig,
) (*sqlx.DB, error) {
	if dbConfig.GetIsTestDB() {
		if err := setupTestDB(dbConfig); err != nil {
			return nil, err
		}
	}
	return sqlx.Connect("postgres", getConnectionString(dbConfig))
}

// NewGORMDatabase has a dummy import of  _ *sqlx.DB to ensure that NewDatabase is called first
func NewGORMDatabase(
	dbConfig config.PSQLConfig, _ *sqlx.DB,
) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(getConnectionString(dbConfig)), &gorm.Config{Logger: gormLogger{}})
}

func setupTestDB(dbConfig config.PSQLConfig) error {
	// the db in the config is not guaranteed to exist yet
	// connect to "postgres" to setup the db defined in the config
	dbName := dbConfig.GetDBName()
	defer dbConfig.SetDBName(dbName)
	dbConfig.SetDBName("postgres")

	db, err := sqlx.Connect("postgres", getConnectionString(dbConfig))
	if err != nil {
		return err
	}
	var count int
	if err := db.QueryRow("SELECT count(*) FROM pg_catalog.pg_database where datname=$1;", dbName).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)); err != nil {
			return err
		}
	}

	dbConfig.SetDBName(dbName)
	logrus.WithField("database", dbConfig.GetDBName()).Info("created new DB")
	return nil
}
