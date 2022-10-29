package database

import (
	"fmt"

	"github.com/dcaf-labs/drip/pkg/service/config"
)

func getConnectionString(dbConfig config.PSQLConfig) string {
	if dbConfig.GetURL() != "" {
		return dbConfig.GetURL()
	}
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConfig.GetHost(),
		dbConfig.GetPort(),
		dbConfig.GetUser(),
		dbConfig.GetPassword(),
		dbConfig.GetDBName())
}
