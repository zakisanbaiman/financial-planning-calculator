package web

import (
	"net/http"
	"os"
	"time"

	"github.com/financial-planning-calculator/backend/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupMiddleware configures all middleware for the Echo server
func SetupMiddleware(e *echo.Echo, cfg *config.ServerConfig) {
	// ログミドルウェア - 詳細なリクエスト/レスポンスログ
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: cfg.LogFormat,
		Output: os.Stdout,
	}))

	// リカバリーミドルウェア - パニック時の復旧
	e.Use(middleware.Recover())

	// CORS設定 - フロントエンドからのアクセス許可
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.AllowedOrigins,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Requested-With",
		},
		AllowCredentials: true,
		MaxAge:           86400, // 24時間
	}))

	// セキュリティヘッダー
	if cfg.EnableSecureHeaders {
		e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
			XSSProtection:         "1; mode=block",
			ContentTypeNosniff:    "nosniff",
			XFrameOptions:         "DENY",
			HSTSMaxAge:            31536000,
			ContentSecurityPolicy: "default-src 'self'",
		}))
	}

	// リクエストサイズ制限
	e.Use(middleware.BodyLimit(cfg.MaxRequestSize))

	// レート制限 - API呼び出し頻度制限
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStore(100), // 100 requests per second
	}))

	// タイムアウト設定
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: cfg.RequestTimeout,
	}))

	// リクエストID生成
	e.Use(middleware.RequestID())

	// Gzip圧縮
	if cfg.EnableGzip {
		e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: cfg.GzipLevel,
		}))
	}
}

// CustomHTTPErrorHandler provides consistent error responses
func CustomHTTPErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	} else {
		msg = err.Error()
	}

	// ログ出力
	c.Logger().Error(err)

	// エラーレスポンス
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, map[string]interface{}{
				"error":      msg,
				"status":     code,
				"timestamp":  time.Now().Format(time.RFC3339),
				"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
			})
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
