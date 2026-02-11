// backend/infrastructure/log/logger.go
package log

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

var logger *slog.Logger

// ContextKey はコンテキストキーの型
type ContextKey string

const (
	// RequestIDKey はリクエストIDのコンテキストキー
	RequestIDKey ContextKey = "request_id"
	// UserIDKey はユーザーIDのコンテキストキー
	UserIDKey ContextKey = "user_id"
	// OperationKey は操作名のコンテキストキー
	OperationKey ContextKey = "operation"
)

func init() {
	// 環境変数からログレベルを取得（デフォルト: INFO）
	level := getLogLevel()
	
	// JSON形式で標準出力にログを出力（構造化ロギング）
	opts := &slog.HandlerOptions{
		Level: level,
		// ソースコードの位置情報を追加
		AddSource: true,
	}
	logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

// getLogLevel は環境変数からログレベルを取得します
func getLogLevel() slog.Level {
	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		return slog.LevelInfo
	}

	switch levelStr {
	case "DEBUG", "debug":
		return slog.LevelDebug
	case "INFO", "info":
		return slog.LevelInfo
	case "WARN", "warn", "WARNING", "warning":
		return slog.LevelWarn
	case "ERROR", "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Logger はグローバルな構造化ロガーを返します。
func Logger() *slog.Logger {
	return logger
}

// WithContext はコンテキストから情報を抽出してロガーに追加します
func WithContext(ctx context.Context) *slog.Logger {
	l := logger

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		l = l.With(slog.String("request_id", requestID))
	}

	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		l = l.With(slog.String("user_id", userID))
	}

	if operation, ok := ctx.Value(OperationKey).(string); ok && operation != "" {
		l = l.With(slog.String("operation", operation))
	}

	return l
}

// WithRequestID はリクエストIDをコンテキストに追加します
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithUserID はユーザーIDをコンテキストに追加します
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// WithOperation は操作名をコンテキストに追加します
func WithOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, OperationKey, operation)
}

// Error はエラーログを出力します（コンテキスト情報付き）
func Error(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	l := WithContext(ctx)
	allAttrs := append([]slog.Attr{
		slog.String("error", err.Error()),
		slog.String("error_type", getErrorType(err)),
		slog.String("stack_trace", getStackTrace()),
		slog.Time("timestamp", time.Now().UTC()),
	}, attrs...)
	l.LogAttrs(ctx, slog.LevelError, msg, allAttrs...)
}

// getStackTrace はスタックトレースを取得します
func getStackTrace() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// Warn は警告ログを出力します（コンテキスト情報付き）
func Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	l := WithContext(ctx)
	allAttrs := append([]slog.Attr{
		slog.Time("timestamp", time.Now().UTC()),
	}, attrs...)
	l.LogAttrs(ctx, slog.LevelWarn, msg, allAttrs...)
}

// Info は情報ログを出力します（コンテキスト情報付き）
func Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	l := WithContext(ctx)
	allAttrs := append([]slog.Attr{
		slog.Time("timestamp", time.Now().UTC()),
	}, attrs...)
	l.LogAttrs(ctx, slog.LevelInfo, msg, allAttrs...)
}

// Debug はデバッグログを出力します（コンテキスト情報付き）
func Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	l := WithContext(ctx)
	allAttrs := append([]slog.Attr{
		slog.Time("timestamp", time.Now().UTC()),
	}, attrs...)
	l.LogAttrs(ctx, slog.LevelDebug, msg, allAttrs...)
}

// getErrorType はエラーの型名を取得します
func getErrorType(err error) string {
	if err == nil {
		return ""
	}
	// runtime.FuncForPC を使ってエラータイプを取得
	t := runtime.FuncForPC(0)
	if t != nil {
		return t.Name()
	}
	return "unknown"
}

// UseCaseLogger はユースケース層用のロガー構造体
type UseCaseLogger struct {
	name string
}

// NewUseCaseLogger は新しいUseCaseLoggerを作成します
func NewUseCaseLogger(name string) *UseCaseLogger {
	return &UseCaseLogger{name: name}
}

// StartOperation は操作開始をログに記録し、操作名を付与したコンテキストを返します
func (l *UseCaseLogger) StartOperation(ctx context.Context, operation string, attrs ...slog.Attr) context.Context {
	ctx = WithOperation(ctx, operation)
	allAttrs := append([]slog.Attr{
		slog.String("usecase", l.name),
		slog.String("phase", "start"),
	}, attrs...)
	Info(ctx, "操作開始: "+operation, allAttrs...)
	return ctx
}

// EndOperation は操作完了をログに記録します
func (l *UseCaseLogger) EndOperation(ctx context.Context, operation string, attrs ...slog.Attr) {
	allAttrs := append([]slog.Attr{
		slog.String("usecase", l.name),
		slog.String("phase", "end"),
	}, attrs...)
	Info(ctx, "操作完了: "+operation, allAttrs...)
}

// OperationError は操作エラーをログに記録します
func (l *UseCaseLogger) OperationError(ctx context.Context, operation string, err error, attrs ...slog.Attr) {
	allAttrs := append([]slog.Attr{
		slog.String("usecase", l.name),
		slog.String("phase", "error"),
	}, attrs...)
	Error(ctx, "操作エラー: "+operation, err, allAttrs...)
}
