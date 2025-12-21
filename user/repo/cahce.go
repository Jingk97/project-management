package repo

import (
	"context"
	"time"
)

// 主要目的是封装缓存层。进行抽象
// 后续无论是redis实现还是mysql实现缓存就不必定论，需要什么实现什么接口即实现了此接口类型

type Cache interface {
	Put(ctx context.Context, key, value string, expire time.Duration) error
	Get(ctx context.Context, key string) (value string, err error)
	Ping(ctx context.Context) error
}
