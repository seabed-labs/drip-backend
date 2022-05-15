package psql

import (
	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(
	config *configs.PSQLConfig,
) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", getConnectionString(config))
}

func NewGORMDatabase(
	config *configs.PSQLConfig,
) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(getConnectionString(config)), &gorm.Config{})
}
