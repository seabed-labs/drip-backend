package database

import (
	"fmt"

	"github.com/dcaf-labs/drip/pkg/service/configs"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(
	config configs.PSQLConfig,
) (*sqlx.DB, error) {
	if config.GetIsTestDB() {
		db, err := sqlx.Connect("postgres", getConnectionString(config))
		if err != nil {
			return nil, err
		}
		config.SetDBName("test_" + uuid.New().String()[0:4])
		_, err = db.Exec(fmt.Sprintf("create database %s", config.GetDBName()))
		if err != nil {
			return nil, err
		}
		logrus.WithField("database", config.GetDBName()).Info("created new DB")
	}
	return sqlx.Connect("postgres", getConnectionString(config))
}

// NewGORMDatabase has a dummy import of  _ *sqlx.DB to ensure that NewDatabase is called first
func NewGORMDatabase(
	config configs.PSQLConfig, _ *sqlx.DB,
) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(getConnectionString(config)), &gorm.Config{Logger: gormLogger{}})
}
