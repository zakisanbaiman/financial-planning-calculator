package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/financial-planning-calculator/backend/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSetupMiddleware_RateLimiter(t *testing.T) {
	e := echo.New()
	cfg := &config.ServerConfig{
		AllowedOrigins: []string{"http://localhost:3000"},
		CORSMaxAge:     86400,
		RateLimitRPS:   2, // 低いレート制限でテスト
		RateLimitBurst: 2,
		RequestTimeout: 30 * time.Second,
		MaxRequestSize: "10M",
		EnableGzip:     false,
		LogFormat:      "${method} ${uri} ${status}\n",
	}

	SetupMiddleware(e, cfg)

	// テスト用のエンドポイント
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// 同一IPからの連続リクエストでレート制限をテスト
	t.Run("レート制限内のリクエストは成功する", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Forwarded-For", "192.168.1.100")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		// 最初のリクエストは成功するはず
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異なるIPからのリクエストは別々にカウントされる", func(t *testing.T) {
		// 新しいEchoインスタンスで隔離テスト
		e2 := echo.New()
		SetupMiddleware(e2, cfg)
		e2.GET("/test", func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})

		// IP1からのリクエスト
		req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
		req1.Header.Set("X-Forwarded-For", "10.0.0.1")
		rec1 := httptest.NewRecorder()
		e2.ServeHTTP(rec1, req1)
		assert.Equal(t, http.StatusOK, rec1.Code)

		// IP2からのリクエスト（別のIPなので制限されない）
		req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
		req2.Header.Set("X-Forwarded-For", "10.0.0.2")
		rec2 := httptest.NewRecorder()
		e2.ServeHTTP(rec2, req2)
		assert.Equal(t, http.StatusOK, rec2.Code)
	})

	t.Run("X-Real-IPヘッダーも認識される", func(t *testing.T) {
		e3 := echo.New()
		SetupMiddleware(e3, cfg)
		e3.GET("/test", func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Real-IP", "172.16.0.1")
		rec := httptest.NewRecorder()
		e3.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestSetupMiddleware_RateLimitExceeded(t *testing.T) {
	e := echo.New()
	cfg := &config.ServerConfig{
		AllowedOrigins: []string{"http://localhost:3000"},
		CORSMaxAge:     86400,
		RateLimitRPS:   1,  // 非常に低いレート制限
		RateLimitBurst: 1,  // バーストも1に制限
		RequestTimeout: 30 * time.Second,
		MaxRequestSize: "10M",
		EnableGzip:     false,
		LogFormat:      "${method} ${uri} ${status}\n",
	}

	SetupMiddleware(e, cfg)

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// 同一IPからの連続リクエストでレート制限超過をテスト
	clientIP := "192.168.100.200"
	var rateLimited bool

	// バースト+レートを超えるリクエストを送信
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Forwarded-For", clientIP)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code == http.StatusTooManyRequests {
			rateLimited = true
			// レスポンス内容を検証
			assert.Contains(t, rec.Body.String(), "Too Many Requests")
			break
		}
	}

	assert.True(t, rateLimited, "レート制限が機能していません")
}
