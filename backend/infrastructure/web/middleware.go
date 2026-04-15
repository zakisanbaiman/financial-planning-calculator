package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/financial-planning-calculator/backend/config"
	"github.com/financial-planning-calculator/backend/infrastructure/log"
	"github.com/financial-planning-calculator/backend/infrastructure/monitoring"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// botMessagesPath はBot SSEエンドポイントのパス
const botMessagesPath = "/api/bot/messages"

// SetupMiddleware configures all middleware for the Echo server.
// Returns the CustomRateLimiterStore so it can be reused for the status endpoint.
func SetupMiddleware(e *echo.Echo, cfg *config.ServerConfig) *CustomRateLimiterStore {
	// パフォーマンス監視ミドルウェア（New Relic APM）
	e.Use(monitoring.NewRelicMiddleware())

	// ログミドルウェア - slog による構造化リクエストログ
	e.Use(SlogRequestLogger())

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

	// セキュリティヘッダー（本番環境: ENABLE_SECURE_HEADERS=true / 開発環境: ENABLE_SECURE_HEADERS=false）
	// 開発環境では Swagger UI が CSP の制約で動作しなくなるため無効化する
	if cfg.EnableSecureHeaders {
		e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
			XSSProtection:         "1; mode=block",
			ContentTypeNosniff:    "nosniff",
			XFrameOptions:         "DENY",
			HSTSMaxAge:            31536000,
			HSTSExcludeSubdomains: false,
			HSTSPreloadEnabled:    true,
			ContentSecurityPolicy: cfg.ContentSecurityPolicy,
		}))
	}

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
		Skipper: func(c echo.Context) bool {
			// ヘルスチェック・メトリクスはレートリミット対象外
			path := c.Path()
			return path == "/health" || path == "/health/detailed" || path == "/ready" || path == "/metrics"
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusTooManyRequests, map[string]any{
				"error":   "Too Many Requests",
				"message": "Rate limit exceeded. Please wait before retrying.",
				"code":    "RATE_LIMIT_EXCEEDED",
			})
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			info := rateLimitStore.GetInfo(identifier)
			retryAfterSec := info.Reset - time.Now().Unix()
			if retryAfterSec < 0 {
				retryAfterSec = 0
			}
			return c.JSON(http.StatusTooManyRequests, map[string]any{
				"error":       "Too Many Requests",
				"message":     "Rate limit exceeded. Please wait before retrying.",
				"code":        "RATE_LIMIT_EXCEEDED",
				"retry_after": fmt.Sprintf("%ds", retryAfterSec),
			})
		},
	}))

	// X-RateLimit-* response headers
	e.Use(RateLimitHeaderMiddleware(rateLimitStore))

	// タイムアウト設定（SSEエンドポイントは除外）
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: cfg.RequestTimeout,
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == botMessagesPath
		},
	}))

	// リクエストID生成
	e.Use(middleware.RequestID())

	// Gzip圧縮（SSEエンドポイントは除外）
	if cfg.EnableGzip {
		e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: cfg.GzipLevel,
			Skipper: func(c echo.Context) bool {
				return c.Request().URL.Path == botMessagesPath
			},
		}))
	}

	return rateLimitStore
}

// extractIdentifier returns the client IP from common proxy headers or the real IP.
func extractIdentifier(c echo.Context) (string, error) {
	// X-Forwarded-For は "client, proxy1, proxy2" 形式で複数IPを含む場合がある。
	// 最左がオリジナルクライアントIPのため、最初のIPのみを使用する。
	// Note: このヘッダーはクライアントが偽装可能。信頼できるリバースプロキシがある
	// 環境では、プロキシが付与する最右IPを使う設計への移行を将来的に検討すること。
	ip := c.Request().Header.Get("X-Forwarded-For")
	if ip != "" {
		parts := strings.Split(ip, ",")
		ip = strings.TrimSpace(parts[0])
	}
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

// AuthRateLimiterMiddleware creates a stricter rate limiter middleware for authentication endpoints.
// This protects against brute-force attacks on login, register, and other auth endpoints.
func AuthRateLimiterMiddleware(cfg *config.ServerConfig) echo.MiddlewareFunc {
	authStore := NewCustomRateLimiterStore(
		float64(cfg.AuthRateLimitRPS),
		cfg.AuthRateLimitBurst,
		5*time.Minute,
	)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			identifier, err := extractIdentifier(c)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]any{
					"error":   "Internal Server Error",
					"message": "Failed to identify client",
					"code":    "INTERNAL_ERROR",
				})
			}

			allowed, _ := authStore.Allow(identifier)
			if !allowed {
				info := authStore.GetInfo(identifier)
				return c.JSON(http.StatusTooManyRequests, map[string]any{
					"error":       "Too Many Requests",
					"message":     "Too many authentication attempts. Please wait before retrying.",
					"code":        "AUTH_RATE_LIMIT_EXCEEDED",
					"retry_after": fmt.Sprintf("%ds", info.Reset-time.Now().Unix()),
				})
			}

			// Add auth-specific rate limit headers
			info := authStore.GetInfo(identifier)
			h := c.Response().Header()
			h.Set("X-Auth-RateLimit-Limit", fmt.Sprintf("%d", info.Limit))
			h.Set("X-Auth-RateLimit-Remaining", fmt.Sprintf("%d", info.Remaining))
			h.Set("X-Auth-RateLimit-Reset", fmt.Sprintf("%d", info.Reset))

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

// SlogRequestLogger は slog を使った構造化リクエストロガーミドルウェアを返します。
// Echo 標準の LoggerMiddleware の代わりに使用し、JSON 構造化ログで統一します。
func SlogRequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			// リクエストIDを取得（RequestID ミドルウェアが後段で付与）
			requestID := res.Header().Get(echo.HeaderXRequestID)
			ctx := log.WithRequestID(req.Context(), requestID)

			latency := time.Since(start)
			status := res.Status

			attrs := []slog.Attr{
				slog.String("method", req.Method),
				slog.String("uri", req.RequestURI),
				slog.Int("status", status),
				slog.String("latency", latency.String()),
				slog.Int64("latency_ms", latency.Milliseconds()),
				slog.String("remote_ip", c.RealIP()),
				slog.String("bytes_in", req.Header.Get(echo.HeaderContentLength)),
				slog.String("bytes_out", strconv.FormatInt(res.Size, 10)),
				slog.String("user_agent", req.UserAgent()),
			}

			if err != nil {
				attrs = append(attrs, slog.String("error", err.Error()))
			}

			switch {
			case status >= 500:
				log.WithContext(ctx).LogAttrs(ctx, slog.LevelError, "request", attrs...)
			case status >= 400:
				log.WithContext(ctx).LogAttrs(ctx, slog.LevelWarn, "request", attrs...)
			default:
				log.WithContext(ctx).LogAttrs(ctx, slog.LevelInfo, "request", attrs...)
			}

			return nil
		}
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
