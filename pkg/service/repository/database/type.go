package database

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

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
