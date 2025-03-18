package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"web-scaffold/settings"
)

var rdb *redis.Client

func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		zap.L().Error("connect to redis failed", zap.Error(err))
	}
	return
}

func Close() (err error) {
	err = rdb.Close()
	if err != nil {
		zap.L().Error("close redis failed", zap.Error(err))
	}
	return
}
