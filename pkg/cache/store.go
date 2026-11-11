package cache

import (
	"context"
	"errors"
)

var (
	ErrMiss          = errors.New("缓存:不存在或已经过期")
	ErrInvalidKey    = errors.New("缓存:键为空")
	ErrInvalidTTL    = errors.New("缓存:过期时长小于零")
	ErrEntryTooLarge = errors.New("缓存:单条记录超过适配器容量限制")
)

type Store interface {
	Get(ctx context.Context, key string)
	Set(ctx context.Context, key string)
	Del(ctx context.Context, key string)
}
