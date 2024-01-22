package ormx

import (
	"context"
	"database/sql"
	driver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func NewDbSplit(cfg MysqlConfig) *DbSplit {
	db := newConnectWithMysql(cfg.DataSourceWrite, cfg)
	sp := &DbSplit{
		readObj:  db,
		writeObj: db,
		logger:   NewLog(cfg.LogLevel),
	}
	if cfg.DataSourceRead != "" {
		sp.readObj = newConnectWithMysql(cfg.DataSourceRead, cfg)
	}
	return sp
}

func newConnectWithMysql(dsn string, cfg MysqlConfig) *gorm.DB {
	db, err := gorm.Open(driver.Open(dsn), &gorm.Config{
		Logger: NewLog(LogLevel),
	})
	if err != nil {
		panic("can't connect to mysql" + err.Error())
	}
	dbGo, err := db.DB()
	if err != nil {
		panic("can't connect to mysql" + err.Error())
	}
	if cfg.MaxLeftTime == 0 {
		cfg.MaxLeftTime = 3600
	}
	if cfg.MaxIdleTime == 0 {
		cfg.MaxIdleTime = 3600
	}
	if cfg.MaxConnIdle == 0 {
		cfg.MaxConnIdle = 10
	}
	if cfg.MaxConnIdle > 50 {
		cfg.MaxConnIdle = 20
	}
	if cfg.MaxOpen > 20 {
		cfg.MaxOpen = 20
	}

	dbGo.SetConnMaxLifetime(time.Duration(cfg.MaxLeftTime * int64(time.Second)))
	dbGo.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime * int64(time.Second)))
	dbGo.SetMaxIdleConns(cfg.MaxConnIdle)
	dbGo.SetMaxOpenConns(cfg.MaxOpen)
	return db
}

type DbSplit struct {
	readObj  *gorm.DB
	writeObj *gorm.DB
	logger   *GormLogger
}

func (sp *DbSplit) Read() *gorm.DB {
	if sp.readObj == nil {
		return sp.Write()
	}
	return sp.readObj
}

func (sp *DbSplit) ReadWithContext(ctx context.Context) *gorm.DB {
	return sp.Read().WithContext(ctx)
}

func (sp *DbSplit) Write() *gorm.DB {
	return sp.writeObj
}

func (sp *DbSplit) WriteWithContext(ctx context.Context) *gorm.DB {
	return sp.writeObj.WithContext(ctx)
}

func (sp *DbSplit) Save(v interface{}) *gorm.DB {
	return sp.Write().Save(v)
}

func (sp *DbSplit) SaveWithContext(ctx context.Context, v interface{}) *gorm.DB {
	return sp.WriteWithContext(ctx).Save(v)
}

func (sp *DbSplit) Create(v interface{}) *gorm.DB {
	return sp.Write().Create(v)
}

func (sp *DbSplit) CreateWithContext(ctx context.Context, v interface{}) *gorm.DB {
	return sp.WriteWithContext(ctx).Create(v)
}

func (sp *DbSplit) Table(v string) *gorm.DB {
	return sp.Read().Table(v)
}

func (sp *DbSplit) Model(v interface{}) *gorm.DB {
	return sp.Read().Model(v)
}

func (sp *DbSplit) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return sp.Write().Transaction(fc, opts...)
}

func (sp *DbSplit) TransactionWithContext(ctx context.Context, fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return sp.WriteWithContext(ctx).Transaction(fc, opts...)
}
func (sp *DbSplit) GetLogger() *GormLogger {
	return sp.logger
}

func (sp *DbSplit) Close() error {
	r, _ := sp.readObj.DB()
	if r != nil {
		_ = r.Close()
	}
	w, _ := sp.writeObj.DB()
	if w != nil {
		_ = w.Close()
	}
	return nil
}
