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

// CustomHTTPErrorHandler provides consistent error responses using our unified error format
func CustomHTTPErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  any
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message

		// Check if it's our custom validation error
		if validationErr, ok := he.Message.(ValidationErrorResponse); ok {
			// ログ出力
			c.Logger().Warnf("Validation error: %+v", validationErr)

			if !c.Response().Committed {
				err = c.JSON(code, validationErr)
				if err != nil {
					c.Logger().Error(err)
				}
			}
			return
		}
	} else {
		msg = err.Error()
	}

	// ログ出力
	c.Logger().Error(err)

	// 統一されたエラーレスポンス形式を使用
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			errorResponse := map[string]any{
				"error":      getErrorMessageFromStatus(code),
				"details":    msg,
				"timestamp":  time.Now().UTC().Format(time.RFC3339),
				"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
				"code":       getErrorCodeFromStatus(code),
			}
			err = c.JSON(code, errorResponse)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}

// getErrorCodeFromStatus returns appropriate error code based on HTTP status
func getErrorCodeFromStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusTooManyRequests:
		return "TOO_MANY_REQUESTS"
	case http.StatusInternalServerError:
		return "INTERNAL_SERVER_ERROR"
	case http.StatusServiceUnavailable:
		return "SERVICE_UNAVAILABLE"
	case http.StatusRequestTimeout:
		return "TIMEOUT"
	case http.StatusUnprocessableEntity:
		return "VALIDATION_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}

// getErrorMessageFromStatus returns appropriate error message based on HTTP status
func getErrorMessageFromStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "リクエストが無効です"
	case http.StatusUnauthorized:
		return "認証が必要です"
	case http.StatusForbidden:
		return "アクセスが拒否されました"
	case http.StatusNotFound:
		return "リソースが見つかりません"
	case http.StatusConflict:
		return "リソースが競合しています"
	case http.StatusTooManyRequests:
		return "リクエスト数が上限を超えています"
	case http.StatusInternalServerError:
		return "内部サーバーエラーが発生しました"
	case http.StatusServiceUnavailable:
		return "サービスが利用できません"
	case http.StatusRequestTimeout:
		return "リクエストがタイムアウトしました"
	case http.StatusUnprocessableEntity:
		return "入力データを処理できません"
	default:
		return "エラーが発生しました"
	}
}
