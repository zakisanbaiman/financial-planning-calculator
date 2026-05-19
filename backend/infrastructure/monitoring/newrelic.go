package monitoring

import (
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"
)

var (
	// New Relic アプリケーションインスタンス（nilの場合は監視無効）
	nrApp *newrelic.Application

	// アクティブなHTTP接続数（アトミックカウンター）
	activeConnections int64

	// アクティブなDB接続数
	dbConnectionsActive int64
)

// InitNewRelic は New Relic エージェントを初期化します
func InitNewRelic(licenseKey, appName string) error {
	if licenseKey == "" {
		return fmt.Errorf("NEW_RELIC_LICENSE_KEY が設定されていません")
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(licenseKey),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		return fmt.Errorf("New Relic エージェントの初期化に失敗しました: %w", err)
	}

	nrApp = app
	return nil
}

// GetApplication は New Relic アプリケーションインスタンスを返します
// New Relic が無効の場合は nil を返します
func GetApplication() *newrelic.Application {
	return nrApp
}

// NewRelicMiddleware は New Relic でHTTPリクエストを計測するミドルウェアを返します
func NewRelicMiddleware() echo.MiddlewareFunc {
	// nrecho ミドルウェア（New Relic が無効な場合は nrApp が nil のまま動く）
	nrechoMiddleware := nrecho.Middleware(nrApp)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// まず nrecho でトランザクションを開始し、その後カスタムメトリクスを記録する
		nrechoWrapped := nrechoMiddleware(next)

		return func(c echo.Context) error {
			atomic.AddInt64(&activeConnections, 1)
			defer atomic.AddInt64(&activeConnections, -1)

			start := time.Now()
			reqSize := c.Request().ContentLength

			err := nrechoWrapped(c)

			duration := time.Since(start).Seconds()
			status := c.Response().Status
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					status = he.Code
				}
			}

			if nrApp != nil {
				method := c.Request().Method
				path := c.Path()
				statusStr := strconv.Itoa(status)

				// リクエスト数（メソッド・エンドポイント・ステータス別）
				nrApp.RecordCustomMetric(
					fmt.Sprintf("Custom/HTTP/Requests/%s/%s/%s", method, sanitizePath(path), statusStr),
					1,
				)

				// レスポンスタイム
				nrApp.RecordCustomMetric(
					fmt.Sprintf("Custom/HTTP/Duration/%s/%s", method, sanitizePath(path)),
					duration,
				)

				// リクエストサイズ
				if reqSize > 0 {
					nrApp.RecordCustomMetric("Custom/HTTP/RequestSize", float64(reqSize))
				}

				// レスポンスサイズ
				nrApp.RecordCustomMetric("Custom/HTTP/ResponseSize", float64(c.Response().Size))

				// アクティブ接続数
				nrApp.RecordCustomMetric("Custom/HTTP/ActiveConnections", float64(atomic.LoadInt64(&activeConnections)))

				// エラーメトリクス
				if status >= 400 {
					severity := "warning"
					errorType := "client_error"
					if status >= 500 {
						severity = "error"
						errorType = "server_error"
					}
					nrApp.RecordCustomMetric(
						fmt.Sprintf("Custom/Errors/%s/%s", errorType, severity),
						1,
					)
				}
			}

			return err
		}
	}
}

// RecordDatabaseQuery はデータベースクエリのメトリクスを記録します
func RecordDatabaseQuery(queryType string, duration time.Duration) {
	if nrApp == nil {
		return
	}
	nrApp.RecordCustomMetric(
		fmt.Sprintf("Custom/Database/QueryDuration/%s", queryType),
		duration.Seconds(),
	)
}

// UpdateDatabaseConnections はアクティブなデータベース接続数を更新します
func UpdateDatabaseConnections(count int) {
	atomic.StoreInt64(&dbConnectionsActive, int64(count))
	if nrApp == nil {
		return
	}
	nrApp.RecordCustomMetric("Custom/Database/ActiveConnections", float64(count))
}

// UpdateCacheHitRatio はキャッシュヒット率を更新します
func UpdateCacheHitRatio(cacheType string, ratio float64) {
	if nrApp == nil {
		return
	}
	nrApp.RecordCustomMetric(fmt.Sprintf("Custom/Cache/HitRatio/%s", cacheType), ratio)
}

// RecordError はエラーメトリクスを記録します
func RecordError(errorType, severity string) {
	if nrApp == nil {
		return
	}
	nrApp.RecordCustomMetric(fmt.Sprintf("Custom/Errors/%s/%s", errorType, severity), 1)
}

// RecordCacheHit はキャッシュヒットを記録します
func RecordCacheHit(cacheType string) {
	if nrApp == nil {
		return
	}
	nrApp.RecordCustomMetric(fmt.Sprintf("Custom/Cache/Hits/%s", cacheType), 1)
}

// RecordCacheMiss はキャッシュミスを記録します
func RecordCacheMiss(cacheType string) {
	if nrApp == nil {
		return
	}
	nrApp.RecordCustomMetric(fmt.Sprintf("Custom/Cache/Misses/%s", cacheType), 1)
}

// WrapHTTPHandler は標準の http.Handler を New Relic でラップします
func WrapHTTPHandler(pattern string, handler http.Handler) (string, http.Handler) {
	if nrApp == nil {
		return pattern, handler
	}
	wrappedPattern, wrappedFunc := newrelic.WrapHandleFunc(nrApp, pattern, handler.ServeHTTP)
	return wrappedPattern, http.HandlerFunc(wrappedFunc)
}

// sanitizePath はパスの `/` を `_` に置換してメトリクス名に使えるようにします
func sanitizePath(path string) string {
	result := make([]byte, len(path))
	for i, c := range []byte(path) {
		if c == '/' {
			result[i] = '_'
		} else {
			result[i] = c
		}
	}
	return string(result)
}
