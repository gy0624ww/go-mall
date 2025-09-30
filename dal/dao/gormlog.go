package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-study-lab/go-mall/common/logger"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	SlowThreshold time.Duration
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		// 500ms 算慢查询，可以放到配置文件中
		SlowThreshold: 500 * time.Millisecond,
	}
}

// var _ gormLogger.Interface = &GormLogger{}
func (l *GormLogger) LogMode(lev gormLogger.LogLevel) gormLogger.Interface {
	return &GormLogger{}
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logger.Info(ctx, msg, "data", data)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logger.Warn(ctx, msg, "data", data)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logger.Error(ctx, msg, "data", data)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// 获取运行时间
	duration := time.Since(begin).Milliseconds()
	// 获取SQL 语句和返回条数
	sql, rows := fc()
	// Gorm 错误时记录错误日志
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error(ctx, "SQL ERROR", "sql", sql, "rows", rows, "dur(ms)", duration)
	}
	// 慢查询日志
	if duration > l.SlowThreshold.Milliseconds() {
		logger.Warn(ctx, "SQL SLOW", "sql", sql, "rows", rows, "dur(ms)", duration)
	} else {
		logger.Debug(ctx, "SQL DEBUG", "sql", sql, "rows", rows, "dur(ms)", duration)
	}
}
