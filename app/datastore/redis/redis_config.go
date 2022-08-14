package redis

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	redisOnce sync.Once
	redisInst *redis.Client
)

func redisInstance() *redis.Client {
	redisOnce.Do(func() {
		redisInst = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Username: "",
			Password: "",
			DB:       0,
		})

		_, err := redisInst.Ping(context.Background()).Result()
		if err != nil {
			panic(err)
		}
	})
	return redisInst
}
