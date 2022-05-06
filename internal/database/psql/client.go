package psql

import (
	"fmt"

	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/jmoiron/sqlx"
)

func NewDatabase(
	config *configs.PSQLConfig,
) (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName)
	return sqlx.Connect("postgres", psqlInfo)
}
