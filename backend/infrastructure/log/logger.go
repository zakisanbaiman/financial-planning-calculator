// backend/infrastructure/log/logger.go
package log

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	// TODO: configからlog levelやformatを読み込めるようにする
	// 今回は簡単のため、JSON形式で標準出力にログを出すように固定
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// Logger はグローバルな構造化ロガーを返します。
func Logger() *slog.Logger {
	return logger
}
