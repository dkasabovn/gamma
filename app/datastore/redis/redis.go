package redis

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	wrapperOnce sync.Once
	wrapperInst *redisWrapper
)

type redisWrapper struct {
	client *redis.Client
}

func GetRedis() Redis {
	wrapperOnce.Do(func() {
		wrapperInst = &redisWrapper{
			client: redisInstance(),
		}
	})
	return wrapperInst
}

func (w *redisWrapper) Get(ctx context.Context, key string) (string, error) {
	return w.client.Get(ctx, key).Result()
}

func (w *redisWrapper) Set(ctx context.Context, key string, value string, exp time.Duration) error {
	return w.client.Set(ctx, key, value, exp).Err()
}

func (w *redisWrapper) Exists(ctx context.Context, key string) bool {
	if w.client.Get(ctx, key).Err() == redis.Nil {
		return false
	}
	return true
}
