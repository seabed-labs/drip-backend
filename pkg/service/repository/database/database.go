package database

import (
	"fmt"

	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(
	dbConfig config.PSQLConfig,
) (*sqlx.DB, error) {
	if dbConfig.GetIsTestDB() {
		db, err := sqlx.Connect("postgres", getConnectionString(dbConfig))
		if err != nil {
			return nil, err
		}
		dbConfig.SetDBName("test_" + uuid.New().String()[0:4])
		_, err = db.Exec(fmt.Sprintf("create database %s", dbConfig.GetDBName()))
		if err != nil {
			return nil, err
		}
		logrus.WithField("database", dbConfig.GetDBName()).Info("created new DB")
	}
	return sqlx.Connect("postgres", getConnectionString(dbConfig))
}

// NewGORMDatabase has a dummy import of  _ *sqlx.DB to ensure that NewDatabase is called first
func NewGORMDatabase(
	dbConfig config.PSQLConfig, _ *sqlx.DB,
) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(getConnectionString(dbConfig)), &gorm.Config{Logger: gormLogger{}})
}
