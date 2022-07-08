package psql

import (
	"context"
	"time"

	"github.com/dcaf-protocol/drip/pkg/configs"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(
	config *configs.PSQLConfig,
) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", getConnectionString(config))
}

type gormLogger struct{}

func (g gormLogger) LogMode(level logger.LogLevel) logger.Interface { return gormLogger{} }

func (g gormLogger) Info(ctx context.Context, s string, i ...interface{}) { logrus.Info(ctx, s, i) }

func (g gormLogger) Warn(ctx context.Context, s string, i ...interface{}) { logrus.Warn(ctx, s, i) }

func (g gormLogger) Error(ctx context.Context, s string, i ...interface{}) { logrus.Error(ctx, s, i) }

func (g gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
}

func NewGORMDatabase(
	config *configs.PSQLConfig,
) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(getConnectionString(config)), &gorm.Config{Logger: gormLogger{}})
}
