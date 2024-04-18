package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"neiro-api/config"
)

var redisClient *redis.Client

func Init() {
	cfg := config.GetConfig().RedisConfig
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

func GetRedis() *redis.Client {
	return redisClient
}
