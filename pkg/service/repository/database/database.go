package database

import (
	"context"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/configs"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(
	config *configs.PSQLConfig,
) (*sqlx.DB, error) {
	if config.IsTestDB {
		db, err := sqlx.Connect("postgres", getConnectionString(config))
		if err != nil {
			return nil, err
		}
		config.DBName = "test_" + uuid.New().String()[0:4]
		_, err = db.Exec("create database " + config.DBName)
		if err != nil {
			return nil, err
		}
		logrus.WithField("database", config.DBName).Info("created new DB")
	}
	return sqlx.Connect("postgres", getConnectionString(config))
}

type gormLogger struct{}

func (g gormLogger) LogMode(level logger.LogLevel) logger.Interface { return gormLogger{} }

func (g gormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	logrus.WithContext(ctx).Info(s, i)
}

func (g gormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	logrus.WithContext(ctx).Warn(s, i)
}

func (g gormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	logrus.WithContext(ctx).Error(s, i)
}

func (g gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
}

// Dummy import of  _ *sqlx.DB to ensure that NewDatabase is called first
func NewGORMDatabase(
	config *configs.PSQLConfig, _ *sqlx.DB,
) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(getConnectionString(config)), &gorm.Config{Logger: gormLogger{}})
}
