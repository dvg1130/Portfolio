package redisdb

import (
	"context"
	"fmt"

	"github.com/dvg1130/Portfolio/secure-backend/config"
	"github.com/redis/go-redis/v9"
)

func RedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.AuthConfig.REDIS_ADDR,
		DB:   0,
	})

	//test connectin
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic("failed to connect to Redis: " + err.Error())
	}
	fmt.Println("connected to redis")
	return rdb

}
