# 監視・ログ・エラー追跡機能

財務計画計算機アプリケーションの監視、ログ、エラー追跡機能のドキュメントです。

## 概要

本アプリケーションでは、サービス運用性向上のために以下の機能を提供しています：

- **Prometheusメトリクス収集**: HTTPリクエスト、データベース、キャッシュなどのパフォーマンスメトリクス
- **構造化ログ**: JSON形式の構造化ログ（slog使用）
- **エラー追跡**: 詳細なスタックトレースとコンテキスト情報を含むエラー追跡

## 機能

### 1. Prometheusメトリクス

#### 収集されるメトリクス

| メトリクス名 | 種類 | 説明 | ラベル |
|------------|------|------|-------|
| `http_requests_total` | Counter | 総HTTPリクエスト数 | method, endpoint, status |
| `http_request_duration_seconds` | Histogram | HTTPリクエスト処理時間（秒） | method, endpoint |
| `http_request_size_bytes` | Histogram | HTTPリクエストサイズ（バイト） | method, endpoint |
| `http_response_size_bytes` | Histogram | HTTPレスポンスサイズ（バイト） | method, endpoint |
| `http_active_connections` | Gauge | アクティブなHTTP接続数 | - |
| `database_query_duration_seconds` | Histogram | データベースクエリ処理時間（秒） | query_type |
| `database_connections_active` | Gauge | アクティブなDB接続数 | - |
| `cache_hit_ratio` | Gauge | キャッシュヒット率（0-1） | cache_type |
| `errors_total` | Counter | エラー総数 | error_type, severity |

#### メトリクスエンドポイント

```
GET /metrics
```

このエンドポイントはPrometheus形式でメトリクスを公開します。

#### 設定

`.env`ファイルで以下の設定が可能です：

```env
# メトリクス収集を有効化（デフォルト: true）
ENABLE_PROMETHEUS_METRICS=true

# メトリクスエンドポイント（デフォルト: /metrics）
METRICS_ENDPOINT=/metrics
```

#### Prometheusサーバーの設定例

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'financial-planning-calculator'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### 2. 構造化ログ

#### ログレベル

環境変数`LOG_LEVEL`で設定可能：

- `DEBUG`: デバッグ情報を含むすべてのログ
- `INFO`: 通常の動作情報（デフォルト）
- `WARN`: 警告レベルのログ
- `ERROR`: エラーのみ

```env
# .env
LOG_LEVEL=INFO
```

#### ログ形式

すべてのログはJSON形式で出力されます：

```json
{
  "time": "2026-02-11T10:00:00Z",
  "level": "ERROR",
  "source": {
    "function": "main.handler",
    "file": "/app/main.go",
    "line": 42
  },
  "msg": "エラーが発生しました",
  "request_id": "abc123",
  "user_id": "user456",
  "error": "database connection failed",
  "error_type": "DatabaseError",
  "stack_trace": "...",
  "timestamp": "2026-02-11T10:00:00Z"
}
```

#### ログコンテキスト

以下のコンテキスト情報が自動的に付与されます：

- `request_id`: リクエストID（X-Request-IDヘッダー）
- `user_id`: ユーザーID（認証済みの場合）
- `operation`: 操作名（ユースケース層で設定）

### 3. エラー追跡

#### ErrorTrackerインターフェース

エラー追跡は`ErrorTracker`インターフェースを通じて行われます：

```go
type ErrorTracker interface {
    CaptureError(ctx context.Context, err error, tags map[string]string)
    CaptureMessage(ctx context.Context, message string, level string, tags map[string]string)
    SetUser(ctx context.Context, userID string, email string)
    Close()
}
```

#### 使用例

```go
import "github.com/financial-planning-calculator/backend/infrastructure/monitoring"

// エラーをキャプチャ
ctx := log.WithRequestID(r.Context(), requestID)
ctx = log.WithUserID(ctx, userID)

tags := map[string]string{
    "component": "payment",
    "severity": "high",
}
monitoring.CaptureError(ctx, err, tags)

// メッセージをキャプチャ
monitoring.CaptureMessage(ctx, "処理が完了しました", "info", tags)
```

#### 拡張性

現在はログベースの実装（`DefaultErrorTracker`）を使用していますが、将来的にSentry等の外部サービスと統合可能な設計になっています。

```env
# .env
ENABLE_ERROR_TRACKING=true
ERROR_TRACKING_ENVIRONMENT=production
```

## 運用

### ヘルスチェック

アプリケーションのヘルスチェックは以下のエンドポイントで確認できます：

```
GET /health          # シンプルなヘルスチェック
GET /health/detailed # 詳細なヘルスチェック（DB接続など）
GET /ready           # 準備状態チェック
```

### パフォーマンスプロファイリング

開発環境ではpprofサーバーが有効化されます：

```env
ENABLE_PPROF=true
PPROF_PORT=6060
```

アクセス方法：
```bash
# CPUプロファイル
go tool pprof http://localhost:6060/debug/pprof/profile

# メモリプロファイル
go tool pprof http://localhost:6060/debug/pprof/heap

# ゴルーチン
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

### ログ集約

構造化ログ（JSON形式）のため、以下のツールと統合が容易です：

- **Fluentd/Fluent Bit**: ログ収集・転送
- **Elasticsearch**: ログ保存・検索
- **Kibana**: ログ可視化
- **Loki**: Grafanaとの統合

Docker環境での例：

```yaml
# docker-compose.yml
services:
  app:
    logging:
      driver: "fluentd"
      options:
        fluentd-address: localhost:24224
        tag: financial-planning-calculator
```

### Grafanaダッシュボード

Prometheusメトリクスを可視化するためのGrafanaダッシュボード例：

1. **HTTPリクエスト統計**
   - リクエスト数（rate）
   - エラー率
   - レスポンスタイム（p50, p95, p99）

2. **データベース統計**
   - クエリ処理時間
   - アクティブ接続数

3. **キャッシュ統計**
   - ヒット率
   - ミス率

4. **エラー統計**
   - エラー発生率
   - エラー種別

## トラブルシューティング

### メトリクスが表示されない

1. Prometheusが初期化されているか確認：
   ```bash
   curl http://localhost:8080/metrics
   ```

2. アプリケーションログでエラーを確認：
   ```bash
   docker logs financial-planning-calculator-backend
   ```

### ログが出力されない

1. ログレベルを確認：
   ```env
   LOG_LEVEL=DEBUG  # より詳細なログを出力
   ```

2. ログ形式が正しいか確認（JSON形式であるべき）

### エラーが追跡されない

1. エラートラッキングが有効化されているか確認：
   ```env
   ENABLE_ERROR_TRACKING=true
   ```

2. コンテキストが正しく設定されているか確認：
   ```go
   ctx = log.WithRequestID(ctx, requestID)
   ctx = log.WithUserID(ctx, userID)
   ```

## 参考資料

- [Prometheus公式ドキュメント](https://prometheus.io/docs/introduction/overview/)
- [Grafana公式ドキュメント](https://grafana.com/docs/)
- [Go slogパッケージ](https://pkg.go.dev/log/slog)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)
