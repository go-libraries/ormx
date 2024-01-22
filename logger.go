package ormx

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm/logger"
	"time"
)

type GormLogger struct {
	logger.LogLevel
	TracerFunc func(ctx context.Context, begin time.Time, sql string, rowsAffected int64, err error)
	CallerSkip int
}

func NewLog(level string) *GormLogger {
	return &GormLogger{
		LogLevel:   getLogLevelFromString(level),
		CallerSkip: 3,
	}
}

func (l *GormLogger) SetTrace(fn func(ctx context.Context, begin time.Time, sql string, rowsAffected int64, err error)) {
	l.TracerFunc = fn
}

func (l *GormLogger) SetTraceFunc(fn func(ctx context.Context, begin time.Time, sql string, rowsAffected int64, err error)) {
	l.SetTrace(fn)
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *GormLogger) Info(ctx context.Context, format string, params ...interface{}) {
	logx.WithContext(ctx).WithCallerSkip(l.CallerSkip).Infof(format, params)
}

func (l *GormLogger) Warn(ctx context.Context, format string, params ...interface{}) {
	logx.WithContext(ctx).WithCallerSkip(l.CallerSkip).Errorf(format, params)
}

func (l *GormLogger) Error(ctx context.Context, format string, params ...interface{}) {
	logx.WithContext(ctx).WithCallerSkip(l.CallerSkip).Errorf(format, params)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sqlStr, rowsAffected := fc()
	if l.TracerFunc != nil {
		l.TracerFunc(ctx, begin, sqlStr, rowsAffected, err)
	}
	if err != nil {
		logx.WithContext(ctx).WithCallerSkip(l.CallerSkip).Errorf("sql: %s, rowsAffected: %d  err: %+v", sqlStr, rowsAffected, err)
	}
	logx.WithContext(ctx).WithCallerSkip(l.CallerSkip).Infof("ql: %s, rowsAffected: %d times: %dms err: %+v", sqlStr, rowsAffected, int64(time.Now().Sub(begin)/time.Millisecond), err)
}
