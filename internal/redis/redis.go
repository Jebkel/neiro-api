package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"neiro-api/config"
	"sync"
)

type ManagerRedis struct {
	ClientRedis *redis.Client
	Ctx         context.Context
}

var (
	instance *ManagerRedis
	once     sync.Once
)

func GetRedis() *ManagerRedis {
	once.Do(func() {
		cfg := config.GetConfig().RedisConfig
		instance = &ManagerRedis{
			ClientRedis: redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
				Password: cfg.Password,
				DB:       cfg.DB,
			}),
			Ctx: context.Background(),
		}
	})
	return instance
}
