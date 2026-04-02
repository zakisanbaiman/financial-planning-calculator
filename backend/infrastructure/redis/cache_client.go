package redis

import (
	"context"
	"time"
)

// CacheClient はキャッシュ操作のインターフェース
// テスト時にモックを注入できるよう、具象型ではなくインターフェースで依存する
type CacheClient interface {
	SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error
	GetJSON(ctx context.Context, key string, dest any) error
	Delete(ctx context.Context, keys ...string) error
	DeleteByPattern(ctx context.Context, pattern string) error
}
