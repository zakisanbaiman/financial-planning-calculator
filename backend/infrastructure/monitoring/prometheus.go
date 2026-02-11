package monitoring

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTPRequestsTotal カウンター: 総リクエスト数
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTPRequestDuration ヒストグラム: リクエスト処理時間
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// HTTPRequestSize ヒストグラム: リクエストサイズ
	HTTPRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7), // 100B ~ 100MB
		},
		[]string{"method", "endpoint"},
	)

	// HTTPResponseSize ヒストグラム: レスポンスサイズ
	HTTPResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7), // 100B ~ 100MB
		},
		[]string{"method", "endpoint"},
	)

	// ActiveConnections ゲージ: アクティブな接続数
	ActiveConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_connections",
			Help: "Number of active HTTP connections",
		},
	)

	// DatabaseQueryDuration ヒストグラム: データベースクエリ処理時間
	DatabaseQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query_type"},
	)

	// DatabaseConnectionsActive ゲージ: アクティブなDB接続数
	DatabaseConnectionsActive = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
		},
	)

	// CacheHitRatio ゲージ: キャッシュヒット率
	CacheHitRatio = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cache_hit_ratio",
			Help: "Cache hit ratio (0-1)",
		},
		[]string{"cache_type"},
	)

	// ErrorsTotal カウンター: エラー総数
	ErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
		[]string{"error_type", "severity"},
	)
)

// InitPrometheus はPrometheusメトリクスを初期化します
func InitPrometheus() {
	// メトリクスを登録
	prometheus.MustRegister(HTTPRequestsTotal)
	prometheus.MustRegister(HTTPRequestDuration)
	prometheus.MustRegister(HTTPRequestSize)
	prometheus.MustRegister(HTTPResponseSize)
	prometheus.MustRegister(ActiveConnections)
	prometheus.MustRegister(DatabaseQueryDuration)
	prometheus.MustRegister(DatabaseConnectionsActive)
	prometheus.MustRegister(CacheHitRatio)
	prometheus.MustRegister(ErrorsTotal)
}

// PrometheusMiddleware はPrometheusメトリクスを収集するミドルウェアです
func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// アクティブ接続数を増やす
			ActiveConnections.Inc()
			defer ActiveConnections.Dec()

			// リクエスト開始時刻
			start := time.Now()

			// リクエストサイズを記録
			if c.Request().ContentLength > 0 {
				HTTPRequestSize.WithLabelValues(
					c.Request().Method,
					c.Path(),
				).Observe(float64(c.Request().ContentLength))
			}

			// 次のハンドラを実行
			err := next(c)

			// 処理時間を計算
			duration := time.Since(start).Seconds()

			// ステータスコードを取得
			status := c.Response().Status
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					status = he.Code
				}
			}

			// メトリクスを記録
			HTTPRequestsTotal.WithLabelValues(
				c.Request().Method,
				c.Path(),
				strconv.Itoa(status),
			).Inc()

			HTTPRequestDuration.WithLabelValues(
				c.Request().Method,
				c.Path(),
			).Observe(duration)

			// レスポンスサイズを記録
			HTTPResponseSize.WithLabelValues(
				c.Request().Method,
				c.Path(),
			).Observe(float64(c.Response().Size))

			// エラーの場合はエラーメトリクスを記録
			if status >= 400 {
				severity := "warning"
				if status >= 500 {
					severity = "error"
				}

				errorType := "client_error"
				if status >= 500 {
					errorType = "server_error"
				}

				ErrorsTotal.WithLabelValues(errorType, severity).Inc()
			}

			return err
		}
	}
}

// PrometheusHandler はPrometheusのメトリクスエンドポイントを提供します
func PrometheusHandler() echo.HandlerFunc {
	h := promhttp.Handler()
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

// RecordDatabaseQuery はデータベースクエリのメトリクスを記録します
func RecordDatabaseQuery(queryType string, duration time.Duration) {
	DatabaseQueryDuration.WithLabelValues(queryType).Observe(duration.Seconds())
}

// UpdateDatabaseConnections はアクティブなデータベース接続数を更新します
func UpdateDatabaseConnections(count int) {
	DatabaseConnectionsActive.Set(float64(count))
}

// UpdateCacheHitRatio はキャッシュヒット率を更新します
func UpdateCacheHitRatio(cacheType string, ratio float64) {
	CacheHitRatio.WithLabelValues(cacheType).Set(ratio)
}

// RecordError はエラーを記録します
func RecordError(errorType, severity string) {
	ErrorsTotal.WithLabelValues(errorType, severity).Inc()
}
