package redisdb

import (
	"context"
	"tiered-service-backend/config"

	"github.com/redis/go-redis/v9"
)

func NewClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Config.REDIS_ADDR,
		DB:   0,
	})

	// Test connection
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic("failed to connect to Redis: " + err.Error())
	}

	return rdb
}
