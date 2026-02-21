package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/financial-planning-calculator/backend/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthRateLimiterMiddleware_AllowsWithinLimit(t *testing.T) {
	cfg := &config.ServerConfig{
		AuthRateLimitRPS:   10,
		AuthRateLimitBurst: 5,
	}

	e := echo.New()
	authLimiter := AuthRateLimiterMiddleware(cfg)

	e.POST("/api/auth/login", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}, authLimiter)

	// バースト内のリクエストは全て成功するはず
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		req.Header.Set("X-Forwarded-For", "10.0.0.1")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "リクエスト %d は成功するはず", i+1)
	}
}

func TestAuthRateLimiterMiddleware_BlocksExcessiveRequests(t *testing.T) {
	cfg := &config.ServerConfig{
		AuthRateLimitRPS:   1,
		AuthRateLimitBurst: 2, // 最大2リクエストまでバースト許可
	}

	e := echo.New()
	authLimiter := AuthRateLimiterMiddleware(cfg)

	e.POST("/api/auth/login", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}, authLimiter)

	clientIP := "192.168.1.50"
	var rateLimited bool

	// バースト + レートを超えるリクエストを送信
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		req.Header.Set("X-Forwarded-For", clientIP)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code == http.StatusTooManyRequests {
			rateLimited = true

			// レスポンス内容を検証
			var body map[string]any
			err := json.Unmarshal(rec.Body.Bytes(), &body)
			require.NoError(t, err)

			assert.Equal(t, "Too Many Requests", body["error"])
			assert.Equal(t, "AUTH_RATE_LIMIT_EXCEEDED", body["code"])
			assert.Contains(t, body["message"], "authentication attempts")
			assert.NotEmpty(t, body["retry_after"])
			break
		}
	}

	assert.True(t, rateLimited, "認証レートリミットが機能していません")
}

func TestAuthRateLimiterMiddleware_DifferentIPsAreSeparate(t *testing.T) {
	cfg := &config.ServerConfig{
		AuthRateLimitRPS:   1,
		AuthRateLimitBurst: 1, // 各IPに1リクエストのみ許可
	}

	e := echo.New()
	authLimiter := AuthRateLimiterMiddleware(cfg)

	e.POST("/api/auth/login", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}, authLimiter)

	// IP1: 最初のリクエストは成功
	req1 := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req1.Header.Set("X-Forwarded-For", "10.1.1.1")
	rec1 := httptest.NewRecorder()
	e.ServeHTTP(rec1, req1)
	assert.Equal(t, http.StatusOK, rec1.Code)

	// IP2: 別のIPも成功（別カウント）
	req2 := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req2.Header.Set("X-Forwarded-For", "10.1.1.2")
	rec2 := httptest.NewRecorder()
	e.ServeHTTP(rec2, req2)
	assert.Equal(t, http.StatusOK, rec2.Code)

	// IP1: 2回目は制限される
	req3 := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req3.Header.Set("X-Forwarded-For", "10.1.1.1")
	rec3 := httptest.NewRecorder()
	e.ServeHTTP(rec3, req3)
	assert.Equal(t, http.StatusTooManyRequests, rec3.Code)
}

func TestAuthRateLimiterMiddleware_ReturnsAuthRateLimitHeaders(t *testing.T) {
	cfg := &config.ServerConfig{
		AuthRateLimitRPS:   10,
		AuthRateLimitBurst: 5,
	}

	e := echo.New()
	authLimiter := AuthRateLimiterMiddleware(cfg)

	e.POST("/api/auth/login", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}, authLimiter)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req.Header.Set("X-Forwarded-For", "10.2.2.1")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// 認証専用のレートリミットヘッダーが含まれること
	assert.NotEmpty(t, rec.Header().Get("X-Auth-RateLimit-Limit"))
	assert.NotEmpty(t, rec.Header().Get("X-Auth-RateLimit-Remaining"))
	assert.NotEmpty(t, rec.Header().Get("X-Auth-RateLimit-Reset"))

	// Limit は burst 値と一致するはず
	assert.Equal(t, "5", rec.Header().Get("X-Auth-RateLimit-Limit"))
}

func TestAuthRateLimiterMiddleware_NonAuthEndpointsUnaffected(t *testing.T) {
	// 非認証エンドポイントにはauth rate limiterが適用されないことを確認
	cfg := &config.ServerConfig{
		AuthRateLimitRPS:   1,
		AuthRateLimitBurst: 1,
	}

	e := echo.New()
	authLimiter := AuthRateLimiterMiddleware(cfg)

	// 認証エンドポイントのみにミドルウェアを適用
	e.POST("/api/auth/login", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}, authLimiter)

	// 非認証エンドポイント（ミドルウェアなし）
	e.GET("/api/financial-data", func(c echo.Context) error {
		return c.String(http.StatusOK, "Financial Data")
	})

	clientIP := "10.3.3.1"

	// 認証エンドポイントのバーストを使い切る
	req1 := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req1.Header.Set("X-Forwarded-For", clientIP)
	rec1 := httptest.NewRecorder()
	e.ServeHTTP(rec1, req1)
	assert.Equal(t, http.StatusOK, rec1.Code)

	// 認証エンドポイントの2回目は制限される
	req2 := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req2.Header.Set("X-Forwarded-For", clientIP)
	rec2 := httptest.NewRecorder()
	e.ServeHTTP(rec2, req2)
	assert.Equal(t, http.StatusTooManyRequests, rec2.Code)

	// 非認証エンドポイントは影響を受けない
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/financial-data", nil)
		req.Header.Set("X-Forwarded-For", clientIP)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "非認証エンドポイントは制限されないはず (リクエスト %d)", i+1)
	}
}
