package datastore

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisRateLimiter struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisRateLimiter(addr string, logger *zap.Logger) *RedisRateLimiter {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisRateLimiter{client: client, logger: logger}
}

func (r *RedisRateLimiter) Increment(key string) (int, error) {
	ctx := context.Background()
	count, err := r.client.Incr(ctx, key).Result()
	return int(count), err
}

func (r *RedisRateLimiter) Get(key string) (int, error) {
	ctx := context.Background()
	count, err := r.client.Get(ctx, key).Int()
	return count, err
}

func (r *RedisRateLimiter) Set(key string, value int) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RedisRateLimiter) Reset(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}
