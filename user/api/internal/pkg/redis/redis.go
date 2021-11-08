package redis

import (
	"context"
	"time"

	"admin/user/api/internal/config"
	"github.com/go-redis/redis"
)

// New 返回初始化后的 redis实例
func New(redisConf config.RedisConfig) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         redisConf.Addr,
		Password:     redisConf.Password,
		DB:           redisConf.DB,
		DialTimeout:  redisConf.DialTimeOut,
		ReadTimeout:  redisConf.ReadTimeOut,
		WriteTimeout: redisConf.WriteTimeOut,
		PoolSize:     redisConf.PoolSize,
		PoolTimeout:  redisConf.PoolTimeOut,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := redisClient.WithContext(ctx).Ping().Result()
	if err != nil {
		return nil, err
	}
	return redisClient, nil
}
