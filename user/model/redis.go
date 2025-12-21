package model

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"time"
)

var RedisClient *RedisCache

type RedisCache struct {
	redisClient *redis.Client
}

type RedisInfo struct {
	Host     string
	Password string
	Port     int
	DB       int
}

func InitRedisCache(redisInfo *RedisInfo) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisInfo.Host + ":" + strconv.Itoa(redisInfo.Port),
		Password: redisInfo.Password,
		DB:       redisInfo.DB,
		//RESP3 (默认) 这是 Redis 6.0 及以上版本推出的新协议，能显式返回更多数据类型（如 Map, Set 等），减少客户端的数据解析开销，提升性能。
		//RESP2 (兼容) 用于兼容旧版 Redis，或者需要严格向后兼容的场景。
		Protocol: 2, // Connection protocol
	})
	RedisClient = &RedisCache{client}
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	err := RedisClient.Ping(ctx)
	if err != nil {
		zap.L().Error("redis连接异常", zap.Error(err))
	}
}

func (rc *RedisCache) Put(ctx context.Context, key string, value string, expire time.Duration) error {
	err := rc.redisClient.Set(ctx, key, value, expire).Err()
	return err
}

func (rc *RedisCache) Get(ctx context.Context, key string) (value string, err error) {
	value, err = rc.redisClient.Get(ctx, key).Result()
	return value, err
}

func (rc *RedisCache) Ping(ctx context.Context) error {
	err := rc.redisClient.Ping(ctx).Err()
	return err
}
