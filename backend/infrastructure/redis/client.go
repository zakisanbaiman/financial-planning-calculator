package redis

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client は go-redis クライアントのラッパーです。
type Client struct {
	rdb *redis.Client
}

// NewClient は環境変数 REDIS_HOST / REDIS_PORT / REDIS_PASSWORD を読み取り、Redis クライアントを生成します。
// REDIS_HOST のデフォルトは "localhost"、REDIS_PORT のデフォルトは "6379" です。
// REDIS_PASSWORD が空文字列の場合は認証なしで接続します。
func NewClient() *Client {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}
	password := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Password:     password,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})

	return &Client{rdb: rdb}
}

// Ping は Redis への疎通確認を行います。
func (c *Client) Ping(ctx context.Context) error {
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis: ping に失敗しました: %w", err)
	}
	return nil
}

// Incr はキーの値を 1 インクリメントし、インクリメント後の値を返します。
func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	val, err := c.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis: INCR %s に失敗しました: %w", key, err)
	}
	return val, nil
}

// Expire はキーに TTL を設定します。キーが存在しない場合は false を返します。
func (c *Client) Expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	ok, err := c.rdb.Expire(ctx, key, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("redis: EXPIRE %s に失敗しました: %w", key, err)
	}
	return ok, nil
}

// TTL はキーの残り TTL を返します。キーが存在しない場合は -2、TTL が設定されていない場合は -1 を返します。
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	d, err := c.rdb.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis: TTL %s に失敗しました: %w", key, err)
	}
	return d, nil
}

// Get はキーの文字列値を返します。
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("redis: GET %s に失敗しました: %w", key, err)
	}
	return val, nil
}

// IsNil は go-redis の redis.Nil エラーかどうかを判定します。
// キーが存在しない場合に返される正常な「not found」状態です。
func IsNil(err error) bool {
	return errors.Is(err, redis.Nil)
}

// Underlying は内部の *redis.Client を返します（テストやパイプライン用）。
func (c *Client) Underlying() *redis.Client {
	return c.rdb
}
