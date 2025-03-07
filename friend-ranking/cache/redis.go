package cache

import (
	"context"
	"friend-ranking/config"
	"github.com/go-redis/redis"
)

var (
	Rdb  *redis.Client
	Rctx context.Context
)

// 初始化Redis配置
func init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: "",
		DB:       0,
	})
	Rctx = context.Background()
}

func Zscore(id int, score int) redis.Z {
	return redis.Z{Score: float64(score), Member: id}
}
