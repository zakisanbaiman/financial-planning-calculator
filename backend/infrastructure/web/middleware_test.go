package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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
		// テストごとにユニークな IP を使用（Redis カウンター干渉防止）
		ip := fmt.Sprintf("192.168.1.%d", time.Now().UnixNano()%200+1)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Forwarded-For", ip)
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

		now := time.Now().UnixNano()
		ip1 := fmt.Sprintf("test-mw-sep-%d-1", now)
		ip2 := fmt.Sprintf("test-mw-sep-%d-2", now)

		// IP1からのリクエスト
		req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
		req1.Header.Set("X-Forwarded-For", ip1)
		rec1 := httptest.NewRecorder()
		e2.ServeHTTP(rec1, req1)
		assert.Equal(t, http.StatusOK, rec1.Code)

		// IP2からのリクエスト（別のIPなので制限されない）
		req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
		req2.Header.Set("X-Forwarded-For", ip2)
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

		ip := fmt.Sprintf("test-real-ip-%d", time.Now().UnixNano())
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Real-IP", ip)
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
		RateLimitRPS:   1, // 非常に低いレート制限
		RateLimitBurst: 1, // バーストも1に制限
		RequestTimeout: 30 * time.Second,
		MaxRequestSize: "10M",
		EnableGzip:     false,
		LogFormat:      "${method} ${uri} ${status}\n",
	}

	SetupMiddleware(e, cfg)

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// テストごとにユニークな IP を使用（Redis カウンター干渉防止）
	clientIP := fmt.Sprintf("test-exceeded-%d", time.Now().UnixNano())
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

func TestExtractIdentifier_MultipleXForwardedFor(t *testing.T) {
	// X-Forwarded-For に複数IPが含まれる場合、最初のIPのみが返されることを検証する。
	tests := []struct {
		name     string
		header   string
		expected string
	}{
		{"単一IP", "192.168.1.1", "192.168.1.1"},
		{"複数IP（スペースあり）", "192.168.1.1, 10.0.0.1, 172.16.0.1", "192.168.1.1"},
		{"複数IP（スペースなし）", "192.168.1.1,10.0.0.1", "192.168.1.1"},
		{"前後スペース", "  192.168.1.100  , 10.0.0.1", "192.168.1.100"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("X-Forwarded-For", tt.header)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			identifier, err := extractIdentifier(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, identifier)
		})
	}
}

func TestSetupMiddleware_DenyHandler_RetryAfterIsDynamic(t *testing.T) {
	// DenyHandler が "60s" のハードコードではなく動的な値を返すことを検証する。
	e := echo.New()
	cfg := &config.ServerConfig{
		AllowedOrigins: []string{"http://localhost:3000"},
		CORSMaxAge:     86400,
		RateLimitRPS:   1,
		RateLimitBurst: 1,
		RequestTimeout: 30 * time.Second,
		MaxRequestSize: "10M",
	}
	SetupMiddleware(e, cfg)
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	clientIP := fmt.Sprintf("test-retry-dynamic-%d", time.Now().UnixNano())
	var retryAfter string
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Forwarded-For", clientIP)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code == http.StatusTooManyRequests {
			var body map[string]any
			json.Unmarshal(rec.Body.Bytes(), &body) //nolint:errcheck
			retryAfter, _ = body["retry_after"].(string)
			break
		}
	}

	assert.NotEmpty(t, retryAfter, "retry_after フィールドが存在しない")
	// ハードコード "60s" ではないことを確認
	assert.NotEqual(t, "60s", retryAfter, "retry_after がハードコードされた '60s' のままです")
	// "<数値>s" 形式で、0以上180秒以下であること（ウィンドウ=3分）
	secStr := strings.TrimSuffix(retryAfter, "s")
	secs, err := strconv.Atoi(secStr)
	assert.NoError(t, err, "retry_after の秒数をパースできない: %s", retryAfter)
	assert.GreaterOrEqual(t, secs, 0)
	assert.LessOrEqual(t, secs, 180)
}
