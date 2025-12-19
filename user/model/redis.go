package model

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var RedisClient *RedisCache

func init() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
		//RESP3 (默认) 这是 Redis 6.0 及以上版本推出的新协议，能显式返回更多数据类型（如 Map, Set 等），减少客户端的数据解析开销，提升性能。
		//RESP2 (兼容) 用于兼容旧版 Redis，或者需要严格向后兼容的场景。
		Protocol: 2, // Connection protocol
	})
	RedisClient = &RedisCache{client}
}

type RedisCache struct {
	redisClient *redis.Client
}

func (rc *RedisCache) Put(ctx context.Context, key string, value string, expire time.Duration) error {
	err := rc.redisClient.Set(ctx, key, value, expire).Err()
	return err
}

func (rc *RedisCache) Get(ctx context.Context, key string) (value string, err error) {
	value, err = rc.redisClient.Get(ctx, key).Result()
	return value, err
}
