package config

import (
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest"
	"time"
)

type Config struct {
	rest.RestConf
	Mysql MySQLConfig
	Redis RedisConfig
	Auth  struct {
		AccessSecret string
		AccessExpire int64
	}
	LogConf logx.LogConf
}

type MySQLConfig struct {
	Name            string
	Addr            string
	UserName        string
	Password        string
	ShowLog         bool
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifeTime int
}

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	DialTimeOut  time.Duration
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
	PoolSize     int
	PoolTimeOut  time.Duration
}
