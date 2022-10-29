package database

import (
	"errors"
	"fmt"

	"github.com/dcaf-labs/drip/pkg/service/configs"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	// Needed for loading drivers
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const migrationDir = "./pkg/service/repository/database/migrations"

// Importing DB mainly so that we can migrate a integrationtest db
func RunMigrations(
	db *sqlx.DB, config *configs.PSQLConfig,
) (int, error) {
	if err := db.Ping(); err != nil {
		logrus.WithField("err", err.Error()).Error("could not ping DB")
		return 0, err
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{DatabaseName: config.DBName})
	if err != nil {
		logrus.WithField("err", err.Error()).Error("could not get DB driver")
		return 0, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		// file://path/to/directory
		fmt.Sprintf("file://%s/%s", configs.GetProjectRoot(), migrationDir),
		config.DBName, driver)
	if err != nil {
		logrus.WithField("err", err.Error()).Error("could not apply migrations")
		return 0, err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logrus.WithField("err", err.Error()).Error("could not sync DB")
		return 0, err
	}
	version, _, err := m.Version()
	logrus.WithField("version", version).Info("database migrated")
	return int(version), err
}
