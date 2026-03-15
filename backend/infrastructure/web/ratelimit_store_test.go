package web

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCustomRateLimiterStore(t *testing.T) {
	store := NewCustomRateLimiterStore(10, 5, time.Minute)
	require.NotNil(t, store)
	assert.Equal(t, 5, store.burst)
}

// TestCustomRateLimiterStore_Allow_* は Redis への接続が不要なフォールバック動作をテストします。
// Redis が利用不可の場合は fail-open（全て許可）されることを確認します。

func TestCustomRateLimiterStore_Allow_FailOpen(t *testing.T) {
	// Redis 接続なしでも Allow は true を返すこと（fail-open）
	// 実際の Redis 接続が利用不可な場合のテスト
	store := NewCustomRateLimiterStore(100, 3, time.Minute)
	require.NotNil(t, store)

	// fail-open: Redis 障害時でもエラーは返さない
	allowed, err := store.Allow("127.0.0.1")
	// エラーは返さない（fail-open の場合は nil）
	assert.NoError(t, err)
	// Redis が利用可能ならバースト内なので true、利用不可なら fail-open で true
	assert.True(t, allowed)
}

func TestCustomRateLimiterStore_GetInfo_ReturnsValidStruct(t *testing.T) {
	store := NewCustomRateLimiterStore(1, 10, time.Minute)
	require.NotNil(t, store)

	// GetInfo は Redis 障害時でもデフォルト値を返すこと
	info := store.GetInfo("10.0.0.1")

	assert.Equal(t, 10, info.Limit)
	assert.GreaterOrEqual(t, info.Remaining, 0)
	assert.LessOrEqual(t, info.Remaining, info.Limit)
	assert.False(t, info.ResetAt.IsZero())
	assert.Greater(t, info.Reset, int64(0))
}

func TestCustomRateLimiterStore_GetInfo_RemainingNotNegative(t *testing.T) {
	store := NewCustomRateLimiterStore(1, 5, time.Minute)

	// Remaining は常に 0 以上であること
	info := store.GetInfo("10.0.0.2")
	assert.GreaterOrEqual(t, info.Remaining, 0)
}

func TestCustomRateLimiterStore_GetInfo_ResetInFuture(t *testing.T) {
	store := NewCustomRateLimiterStore(1, 3, time.Minute)

	// ResetAt は現在時刻以降であること
	info := store.GetInfo("10.0.0.3")
	assert.False(t, info.ResetAt.IsZero())
	// Reset フィールド（Unix 秒）は現在より後
	assert.GreaterOrEqual(t, info.Reset, time.Now().Unix()-1)
}

func TestCustomRateLimiterStore_Allow_DoesNotReturnError(t *testing.T) {
	store := NewCustomRateLimiterStore(100, 5, time.Minute)

	// Allow は Redis 障害時でもエラーを返さないこと（fail-open）
	_, err := store.Allow("192.168.0.1")
	assert.NoError(t, err)
}
