package web

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	redisclient "github.com/financial-planning-calculator/backend/infrastructure/redis"
)

// RateLimitInfo は特定の識別子に対するレートリミットの現在状態を保持します。
type RateLimitInfo struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Reset     int64     `json:"reset"`    // Unix タイムスタンプ（秒）
	ResetAt   time.Time `json:"reset_at"` // 人間が読める形式のリセット時刻
}

// CustomRateLimiterStore は Redis を使ったレートリミットストアです。
// echo middleware.RateLimiterStore を実装し、per-IP のリミット情報も公開します。
//
// アルゴリズム: INCR + EXPIRE による固定ウィンドウカウンター
//   - 各識別子ごとに "ratelimit:<window>:<identifier>" というキーを使用
//   - Redisが利用不可の場合は fail-open（リクエストを通す）
type CustomRateLimiterStore struct {
	burst     int
	window    time.Duration
	redis     *redisclient.Client
}

// NewCustomRateLimiterStore は新しい CustomRateLimiterStore を生成します。
// rps パラメータは後方互換性のために受け取りますが、Redis 実装では burst と window を使用します。
func NewCustomRateLimiterStore(rps float64, burst int, window time.Duration) *CustomRateLimiterStore {
	return &CustomRateLimiterStore{
		burst:  burst,
		window: window,
		redis:  redisclient.NewClient(),
	}
}

// redisKey は識別子に対応する Redis キーを返します。
// 固定ウィンドウのため、window 単位で切り捨てた Unix 時刻をキーに含めます。
func (s *CustomRateLimiterStore) redisKey(identifier string) string {
	windowStart := time.Now().Truncate(s.window).Unix()
	return fmt.Sprintf("ratelimit:%d:%s", windowStart, identifier)
}

// Allow は echo middleware.RateLimiterStore を実装します。
// リクエストが許可される場合は true、レートリミット超過の場合は false を返します。
// Redis 障害時は fail-open でリクエストを通します。
func (s *CustomRateLimiterStore) Allow(identifier string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	key := s.redisKey(identifier)

	// INCR でカウンターをインクリメント
	count, err := s.redis.Incr(ctx, key)
	if err != nil {
		// Redis 障害時: fail-open（リクエストを通す）してエラーログを記録
		slog.Error("レートリミット: Redis INCR に失敗しました（fail-open で通過）",
			slog.String("identifier", identifier),
			slog.String("error", err.Error()),
		)
		return true, nil
	}

	// 最初のリクエスト時のみ EXPIRE を設定
	if count == 1 {
		_, expireErr := s.redis.Expire(ctx, key, s.window)
		if expireErr != nil {
			slog.Error("レートリミット: Redis EXPIRE に失敗しました",
				slog.String("identifier", identifier),
				slog.String("error", expireErr.Error()),
			)
		}
	}

	return count <= int64(s.burst), nil
}

// GetInfo は識別子の現在のレートリミット状態を返します。
// トークンを消費しません。Redis 障害時はデフォルト値（フル）を返します。
func (s *CustomRateLimiterStore) GetInfo(identifier string) RateLimitInfo {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	key := s.redisKey(identifier)

	// 現在のカウンターを取得
	val, err := s.redis.Get(ctx, key)
	if err != nil {
		if !redisclient.IsNil(err) {
			slog.Error("レートリミット: Redis GET に失敗しました",
				slog.String("identifier", identifier),
				slog.String("error", err.Error()),
			)
		}
		// キーが存在しない、または Redis 障害時はフル状態を返す
		return RateLimitInfo{
			Limit:     s.burst,
			Remaining: s.burst,
			Reset:     time.Now().Add(s.window).Unix(),
			ResetAt:   time.Now().Add(s.window).UTC(),
		}
	}

	// 現在のカウント値をパース
	var count int
	fmt.Sscanf(val, "%d", &count)

	remaining := s.burst - count
	if remaining < 0 {
		remaining = 0
	}

	// TTL からリセット時刻を計算
	ttl, ttlErr := s.redis.TTL(ctx, key)
	var resetAt time.Time
	if ttlErr != nil || ttl <= 0 {
		resetAt = time.Now().Add(s.window)
	} else {
		resetAt = time.Now().Add(ttl)
	}

	return RateLimitInfo{
		Limit:     s.burst,
		Remaining: remaining,
		Reset:     resetAt.Unix(),
		ResetAt:   resetAt.UTC(),
	}
}
