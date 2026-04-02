package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Set は文字列値を TTL 付きで保存します。
func (c *Client) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if err := c.rdb.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("redis: SET %s に失敗しました: %w", key, err)
	}
	return nil
}

// SetJSON は任意の値を JSON シリアライズして TTL 付きで保存します。
func (c *Client) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("redis: JSON シリアライズに失敗しました (key=%s): %w", key, err)
	}
	return c.Set(ctx, key, string(b), ttl)
}

// GetJSON はキャッシュから JSON をデシリアライズして dest に格納します。
// キャッシュミスの場合は redis.Nil エラー（IsNil で判定可能）を透過的に返します。
func (c *Client) GetJSON(ctx context.Context, key string, dest any) error {
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		// redis.Nil（キャッシュミス）をそのまま伝播させ、呼び出し側で IsNil() 判定を可能にする
		return err
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("redis: JSON デシリアライズに失敗しました (key=%s): %w", key, err)
	}
	return nil
}

// Delete は指定したキーを削除します。
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("redis: DEL に失敗しました: %w", err)
	}
	return nil
}

// DeleteByPattern は SCAN + DEL でパターンに一致するキーを削除します。
// KEYS コマンドは本番環境で全キースキャンが発生するため使用しない。
func (c *Client) DeleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := c.rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("redis: SCAN %s に失敗しました: %w", pattern, err)
		}
		if len(keys) > 0 {
			if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("redis: パターン削除の DEL に失敗しました: %w", err)
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}
