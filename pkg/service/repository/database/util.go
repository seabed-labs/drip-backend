package database

import (
	"fmt"

	"github.com/dcaf-labs/drip/pkg/service/configs"
)

func getConnectionString(config configs.PSQLConfig) string {
	if config.GetURL() != "" {
		return config.GetURL()
	}
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.GetHost(),
		config.GetPort(),
		config.GetUser(),
		config.GetPassword(),
		config.GetDBName())
}
