package ormx

import "gorm.io/gorm/logger"

type MysqlConfig struct {
	DataSourceWrite string // mysql链接地址，满足 $user:$password@tcp($ip:$port)/$db?$queries 格式即可
	DataSourceRead  string
	MaxConnIdle     int    `json:",default=20"`
	MaxIdleTime     int64  `json:",default=3600"`
	MaxLeftTime     int64  `json:",default=3600"`
	MaxOpen         int    `json:",default=30"`
	LogLevel        string `json:",default=info"`
}

func getLogLevelFromString(level string) logger.LogLevel {
	var l logger.LogLevel
	switch level { //debug,info,error,severe
	case "debug", "info":
		l = logger.Info
	case "error":
		l = logger.Error
	case "severe":
		l = logger.Silent
	default:
		l = logger.Info
	}
	return l
}
