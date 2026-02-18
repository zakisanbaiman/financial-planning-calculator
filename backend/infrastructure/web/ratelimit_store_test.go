package web

import (
	"math"
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

func TestCustomRateLimiterStore_Allow(t *testing.T) {
	// burst=3, rps=100 → 3 requests should be allowed immediately
	store := NewCustomRateLimiterStore(100, 3, time.Minute)

	for i := 0; i < 3; i++ {
		allowed, err := store.Allow("127.0.0.1")
		require.NoError(t, err)
		assert.True(t, allowed, "request %d should be allowed", i+1)
	}

	// 4th request should be denied (burst exhausted)
	allowed, err := store.Allow("127.0.0.1")
	require.NoError(t, err)
	assert.False(t, allowed, "4th request should be denied")
}

func TestCustomRateLimiterStore_Allow_DifferentIdentifiers(t *testing.T) {
	store := NewCustomRateLimiterStore(100, 1, time.Minute)

	// First request for each IP should be allowed
	ok1, _ := store.Allow("192.168.0.1")
	ok2, _ := store.Allow("192.168.0.2")
	assert.True(t, ok1)
	assert.True(t, ok2)

	// Second request for each IP should be denied
	ok1again, _ := store.Allow("192.168.0.1")
	ok2again, _ := store.Allow("192.168.0.2")
	assert.False(t, ok1again)
	assert.False(t, ok2again)
}

func TestCustomRateLimiterStore_GetInfo_Full(t *testing.T) {
	burst := 10
	store := NewCustomRateLimiterStore(1, burst, time.Minute)

	// No requests consumed yet → remaining should equal burst
	info := store.GetInfo("10.0.0.1")

	assert.Equal(t, burst, info.Limit)
	assert.Equal(t, burst, info.Remaining)
	// When full, reset is approximately now
	assert.WithinDuration(t, time.Now(), info.ResetAt, 2*time.Second)
}

func TestCustomRateLimiterStore_GetInfo_AfterConsumption(t *testing.T) {
	burst := 5
	store := NewCustomRateLimiterStore(1, burst, time.Minute)
	ip := "10.0.0.2"

	// Consume 3 tokens
	for i := 0; i < 3; i++ {
		store.Allow(ip) //nolint:errcheck
	}

	info := store.GetInfo(ip)
	assert.Equal(t, burst, info.Limit)
	// Remaining should be around burst-3 (may be slightly different due to time-based refill)
	assert.LessOrEqual(t, info.Remaining, burst-3)
	// Reset time should be in the future
	assert.True(t, info.ResetAt.After(time.Now()) || math.Abs(float64(info.Remaining-burst)) < 1,
		"reset should be in the future when tokens are missing")
}

func TestCustomRateLimiterStore_GetInfo_DoesNotConsumeToken(t *testing.T) {
	burst := 5
	store := NewCustomRateLimiterStore(100, burst, time.Minute)
	ip := "10.0.0.3"

	// Call GetInfo multiple times without Allow
	for i := 0; i < 10; i++ {
		info := store.GetInfo(ip)
		assert.Equal(t, burst, info.Remaining, "GetInfo should not consume tokens")
	}
}

func TestCustomRateLimiterStore_GetInfo_ResetFields(t *testing.T) {
	store := NewCustomRateLimiterStore(1, 3, time.Minute)
	ip := "10.0.0.4"

	// Consume all tokens at rps=1 so reset will be in the future
	store.Allow(ip) //nolint:errcheck
	store.Allow(ip) //nolint:errcheck
	store.Allow(ip) //nolint:errcheck

	info := store.GetInfo(ip)
	// Reset Unix timestamp should be non-zero and in the future
	assert.Greater(t, info.Reset, time.Now().Unix()-1)
	// ResetAt should not be zero value
	assert.False(t, info.ResetAt.IsZero())
}
