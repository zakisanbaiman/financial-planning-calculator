package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/financial-planning-calculator/backend/config"
	"github.com/financial-planning-calculator/backend/infrastructure/log"
	"github.com/financial-planning-calculator/backend/infrastructure/monitoring"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupMiddleware configures all middleware for the Echo server.
// Returns the CustomRateLimiterStore so it can be reused for the status endpoint.
func SetupMiddleware(e *echo.Echo, cfg *config.ServerConfig) *CustomRateLimiterStore {
	// パフォーマンス監視ミドルウェア（Prometheus）
	e.Use(monitoring.PrometheusMiddleware())

	// ログミドルウェア - 詳細なリクエスト/レスポンスログ
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: cfg.LogFormat,
		Output: os.Stdout,
	}))

	// リカバリーミドルウェア - パニック時の復旧とエラー追跡
	e.Use(RecoveryMiddlewareWithErrorTracking())

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
		MaxAge:           cfg.CORSMaxAge,
	}))

	// セキュリティヘッダー（開発環境ではSwagger UI動作のため無効化）
	// TODO: 本番環境では適切なCSPを設定すること
	// if cfg.EnableSecureHeaders {
	// 	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
	// 		XSSProtection:      "1; mode=block",
	// 		ContentTypeNosniff: "nosniff",
	// 		XFrameOptions:      "DENY",
	// 		HSTSMaxAge:         31536000,
	// 	}))
	// }

	// リクエストサイズ制限
	e.Use(middleware.BodyLimit(cfg.MaxRequestSize))

	// Rate limiting - per-IP API request throttling (custom store for /api/rate-limit/status)
	rateLimitStore := NewCustomRateLimiterStore(
		float64(cfg.RateLimitRPS),
		cfg.RateLimitBurst,
		3*time.Minute,
	)
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: rateLimitStore,
		IdentifierExtractor: extractIdentifier,
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusTooManyRequests, map[string]any{
				"error":   "Too Many Requests",
				"message": "Rate limit exceeded. Please wait before retrying.",
				"code":    "RATE_LIMIT_EXCEEDED",
			})
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return c.JSON(http.StatusTooManyRequests, map[string]any{
				"error":       "Too Many Requests",
				"message":     "Rate limit exceeded. Please wait before retrying.",
				"code":        "RATE_LIMIT_EXCEEDED",
				"retry_after": "60s",
			})
		},
	}))

	// X-RateLimit-* response headers
	e.Use(RateLimitHeaderMiddleware(rateLimitStore))

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

	return rateLimitStore
}

// extractIdentifier returns the client IP from common proxy headers or the real IP.
func extractIdentifier(c echo.Context) (string, error) {
	ip := c.Request().Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.Request().Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = c.RealIP()
	}
	return ip, nil
}

// RateLimitHeaderMiddleware attaches X-RateLimit-{Limit,Remaining,Reset} headers to every response.
func RateLimitHeaderMiddleware(store *CustomRateLimiterStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			identifier, _ := extractIdentifier(c)
			info := store.GetInfo(identifier)

			h := c.Response().Header()
			h.Set("X-RateLimit-Limit", fmt.Sprintf("%d", info.Limit))
			h.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", info.Remaining))
			h.Set("X-RateLimit-Reset", fmt.Sprintf("%d", info.Reset))

			return next(c)
		}
	}
}

// CustomHTTPErrorHandler provides consistent error responses using our unified error format
func CustomHTTPErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  any
	)

	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := log.WithRequestID(c.Request().Context(), requestID)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message

		// Check if it's our custom validation error
		if validationErr, ok := he.Message.(ValidationErrorResponse); ok {
			// 構造化ログ出力
			log.Warn(ctx, "バリデーションエラー",
				slog.Int("status_code", code),
				slog.String("path", c.Request().URL.Path),
				slog.String("method", c.Request().Method),
			)

			if !c.Response().Committed {
				err = c.JSON(code, validationErr)
				if err != nil {
					log.Error(ctx, "レスポンス送信エラー", err)
				}
			}
			return
		}
	} else {
		msg = err.Error()
	}

	// 構造化エラーログ出力
	log.Error(ctx, "HTTPエラー", err,
		slog.Int("status_code", code),
		slog.String("path", c.Request().URL.Path),
		slog.String("method", c.Request().Method),
	)

	// 統一されたエラーレスポンス形式を使用
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			errorResponse := map[string]any{
				"error":      getErrorMessageFromStatus(code),
				"details":    msg,
				"timestamp":  time.Now().UTC().Format(time.RFC3339),
				"request_id": requestID,
				"code":       getErrorCodeFromStatus(code),
			}
			err = c.JSON(code, errorResponse)
		}
		if err != nil {
			log.Error(ctx, "レスポンス送信エラー", err)
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

// RecoveryMiddlewareWithErrorTracking はパニック時の復旧とエラー追跡を提供します
func RecoveryMiddlewareWithErrorTracking() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					var err error
					switch x := r.(type) {
					case string:
						err = echo.NewHTTPError(http.StatusInternalServerError, x)
					case error:
						err = echo.NewHTTPError(http.StatusInternalServerError, x.Error())
					default:
						err = echo.NewHTTPError(http.StatusInternalServerError, "パニックが発生しました")
					}

					// リクエストIDを取得
					requestID := c.Response().Header().Get(echo.HeaderXRequestID)
					ctx := log.WithRequestID(c.Request().Context(), requestID)

					// エラー追跡システムに記録
					tags := map[string]string{
						"panic":      "true",
						"method":     c.Request().Method,
						"path":       c.Path(),
						"request_id": requestID,
					}
					monitoring.CaptureError(ctx, err, tags)

					// 構造化ログ出力（スタックトレース付き）
					log.Error(ctx, "パニックが発生しました", err,
						slog.String("path", c.Request().URL.Path),
						slog.String("method", c.Request().Method),
					)

					// Prometheusメトリクスに記録
					monitoring.RecordError("panic", "critical")

					// エラーレスポンスを返す
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
