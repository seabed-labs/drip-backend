package psql

import (
	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/jmoiron/sqlx"
)

func NewDatabase(
	config *configs.PSQLConfig,
) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", getConnectionString(config))
}
