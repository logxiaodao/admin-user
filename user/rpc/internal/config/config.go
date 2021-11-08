package config

import (
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/zrpc"
	"time"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql    MySQLConfig
	Redis    RedisConfig
	LogConf  logx.LogConf
	AuthConf AuthConfig
}

type AuthConfig struct {
	AccessSecret string
	AccessExpire int64
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
	Host         string
	Password     string
	DB           int
	DialTimeOut  time.Duration
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
	PoolSize     int
	PoolTimeOut  time.Duration
}
