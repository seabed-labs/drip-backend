package database

import (
	"fmt"

	"github.com/dcaf-labs/drip/pkg/service/configs"
)

func getConnectionString(config *configs.PSQLConfig) string {
	if config.URL != "" {
		return config.URL
	}
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName)
}
