package Initialization

import (
	config2 "admin-user/rpc/internal/config"
	mysql2 "admin-user/rpc/internal/pkg/mysql"
	redis2 "admin-user/rpc/internal/pkg/redis"
	goredis "github.com/go-redis/redis"
	"gorm.io/gorm"
)

// DefaultDataSource 全局的DataSource，会用来初始化各个模块的repository（不是直接调用变量，而是调用GetDataSource）
var DefaultDataSource *DataSources

type DataSources struct {
	DB          *gorm.DB
	RedisClient *goredis.Client
}

// NewDataSources 根据给定配置初始化 Mysql,Redis
func NewDataSources(config config2.Config) (*DataSources, error) {
	// 初始化mysql
	db, err := mysql2.New(config.Mysql)
	if err != nil {
		return nil, err
	}

	// 初始化redis
	redisClient, err := redis2.New(config.Redis)
	if err != nil {
		return nil, err
	}
	ds := &DataSources{
		DB:          db,
		RedisClient: redisClient,
	}
	DefaultDataSource = ds
	return ds, nil
}

// GetDataSource 会被其他模块调用，返回全局的DataSource
func GetDataSource() *DataSources {
	return DefaultDataSource
}
