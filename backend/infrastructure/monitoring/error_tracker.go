package monitoring

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/financial-planning-calculator/backend/infrastructure/log"
)

// ErrorTracker はエラー追跡のインターフェース
type ErrorTracker interface {
	// CaptureError はエラーを記録します
	CaptureError(ctx context.Context, err error, tags map[string]string)
	// CaptureMessage はメッセージを記録します
	CaptureMessage(ctx context.Context, message string, level string, tags map[string]string)
	// SetUser はユーザー情報を設定します
	SetUser(ctx context.Context, userID string, email string)
	// Close はエラートラッカーをクローズします
	Close()
}

// ErrorContext はエラーのコンテキスト情報を保持します
type ErrorContext struct {
	RequestID   string
	UserID      string
	Timestamp   time.Time
	Environment string
	Tags        map[string]string
	Extra       map[string]interface{}
	StackTrace  string
}

// DefaultErrorTracker はデフォルトのエラートラッカー（ログベース）
type DefaultErrorTracker struct {
	environment string
}

// NewDefaultErrorTracker は新しいデフォルトエラートラッカーを作成します
func NewDefaultErrorTracker(environment string) *DefaultErrorTracker {
	return &DefaultErrorTracker{
		environment: environment,
	}
}

// CaptureError はエラーを記録します
func (t *DefaultErrorTracker) CaptureError(ctx context.Context, err error, tags map[string]string) {
	errCtx := t.buildErrorContext(ctx, tags)
	
	attrs := []slog.Attr{
		slog.String("environment", errCtx.Environment),
		slog.Time("captured_at", errCtx.Timestamp),
		slog.String("stack_trace", errCtx.StackTrace),
	}
	
	// タグを追加
	for key, value := range errCtx.Tags {
		attrs = append(attrs, slog.String("tag_"+key, value))
	}
	
	// Extraフィールドを追加
	for key, value := range errCtx.Extra {
		attrs = append(attrs, slog.Any("extra_"+key, value))
	}
	
	log.Error(ctx, "エラーが発生しました", err, attrs...)
	
	// Prometheusメトリクスにも記録
	RecordError("application_error", "error")
}

// CaptureMessage はメッセージを記録します
func (t *DefaultErrorTracker) CaptureMessage(ctx context.Context, message string, level string, tags map[string]string) {
	errCtx := t.buildErrorContext(ctx, tags)
	
	attrs := []slog.Attr{
		slog.String("environment", errCtx.Environment),
		slog.Time("captured_at", errCtx.Timestamp),
		slog.String("level", level),
	}
	
	// タグを追加
	for key, value := range errCtx.Tags {
		attrs = append(attrs, slog.String("tag_"+key, value))
	}
	
	switch level {
	case "error":
		log.Error(ctx, message, fmt.Errorf("%s", message), attrs...)
	case "warning":
		log.Warn(ctx, message, attrs...)
	default:
		log.Info(ctx, message, attrs...)
	}
}

// SetUser はユーザー情報を設定します
func (t *DefaultErrorTracker) SetUser(ctx context.Context, userID string, email string) {
	// コンテキストにユーザー情報を設定
	ctx = log.WithUserID(ctx, userID)
}

// Close はエラートラッカーをクローズします
func (t *DefaultErrorTracker) Close() {
	// デフォルト実装では何もしない
}

// buildErrorContext はエラーコンテキストを構築します
func (t *DefaultErrorTracker) buildErrorContext(ctx context.Context, tags map[string]string) *ErrorContext {
	errCtx := &ErrorContext{
		Timestamp:   time.Now().UTC(),
		Environment: t.environment,
		Tags:        tags,
		Extra:       make(map[string]interface{}),
		StackTrace:  getStackTrace(3), // 3つのフレームをスキップ
	}
	
	// コンテキストからリクエストIDを取得
	if requestID, ok := ctx.Value(log.RequestIDKey).(string); ok {
		errCtx.RequestID = requestID
		errCtx.Extra["request_id"] = requestID
	}
	
	// コンテキストからユーザーIDを取得
	if userID, ok := ctx.Value(log.UserIDKey).(string); ok {
		errCtx.UserID = userID
		errCtx.Extra["user_id"] = userID
	}
	
	// コンテキストから操作名を取得
	if operation, ok := ctx.Value(log.OperationKey).(string); ok {
		errCtx.Extra["operation"] = operation
	}
	
	return errCtx
}

// getStackTrace はスタックトレースを取得します
func getStackTrace(skip int) string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	
	frames := runtime.CallersFrames(pcs[:n])
	stackTrace := ""
	
	for {
		frame, more := frames.Next()
		stackTrace += fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		
		if !more {
			break
		}
	}
	
	return stackTrace
}

// グローバルエラートラッカー
var globalErrorTracker ErrorTracker

// InitErrorTracker はエラートラッカーを初期化します
func InitErrorTracker(environment string) {
	// 将来的にはSentry等の外部サービスを使用可能
	// 現在はデフォルトのログベース実装を使用
	globalErrorTracker = NewDefaultErrorTracker(environment)
}

// GetErrorTracker はグローバルエラートラッカーを取得します
func GetErrorTracker() ErrorTracker {
	if globalErrorTracker == nil {
		// デフォルトで開発環境として初期化
		InitErrorTracker("development")
	}
	return globalErrorTracker
}

// CaptureError はグローバルエラートラッカーを使用してエラーを記録します
func CaptureError(ctx context.Context, err error, tags map[string]string) {
	GetErrorTracker().CaptureError(ctx, err, tags)
}

// CaptureMessage はグローバルエラートラッカーを使用してメッセージを記録します
func CaptureMessage(ctx context.Context, message string, level string, tags map[string]string) {
	GetErrorTracker().CaptureMessage(ctx, message, level, tags)
}
