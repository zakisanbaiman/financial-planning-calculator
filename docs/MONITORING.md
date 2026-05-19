# 監視・ログ・エラー追跡機能

財務計画計算機アプリケーションの監視、ログ、エラー追跡機能のドキュメントです。

## 概要

本アプリケーションでは **New Relic APM** を使用してサービスを監視しています。

- **APM（アプリケーションパフォーマンス監視）**: HTTPリクエスト、レスポンスタイム、エラー率
- **分散トレーシング**: リクエスト内部の処理フローを可視化
- **カスタムメトリクス**: DB接続数、キャッシュヒット率などのアプリ固有メトリクス
- **アラート**: 閾値超過時の Slack 通知
- **構造化ログ**: JSON形式の構造化ログ（slog使用）

## セットアップ

### 1. New Relic アカウント登録

1. [New Relic サインアップ](https://newrelic.com/signup)（メールのみ、クレカ不要）
2. Settings → API Keys → **INGEST - LICENSE** の Key をコピー

### 2. 環境変数の設定

```env
# backend/.env
NEW_RELIC_LICENSE_KEY=your_license_key_here
NEW_RELIC_APP_NAME=financial-planning-calculator
```

### 3. アプリケーション起動

```bash
docker compose up -d
```

起動後、数分で New Relic にデータが届き始めます。

---

## New Relic でのメトリクス確認

### APM ダッシュボード

**アクセス**: [one.newrelic.com](https://one.newrelic.com) → APM → financial-planning-calculator

| 画面 | 確認できる内容 |
|------|--------------|
| Summary | スループット、レスポンスタイム、エラー率のリアルタイム概要 |
| Transactions | エンドポイント別のパフォーマンス一覧 |
| Databases | クエリ別の処理時間 |
| Errors | エラー一覧とスタックトレース |
| Distributed Tracing | リクエスト内部の処理フロー |

### カスタムメトリクス

以下のカスタムメトリクスが New Relic に送信されます：

| メトリクス名（NRQL） | 説明 |
|---------------------|------|
| `Custom/HTTP/Requests/{method}/{path}/{status}` | HTTPリクエスト数 |
| `Custom/HTTP/Duration/{method}/{path}` | リクエスト処理時間（秒） |
| `Custom/HTTP/ActiveConnections` | アクティブなHTTP接続数 |
| `Custom/Database/QueryDuration/{query_type}` | DBクエリ処理時間（秒） |
| `Custom/Database/ActiveConnections` | アクティブなDB接続数 |
| `Custom/Cache/HitRatio/{cache_type}` | キャッシュヒット率（0-1） |
| `Custom/Cache/Hits/{cache_type}` | キャッシュヒット数 |
| `Custom/Cache/Misses/{cache_type}` | キャッシュミス数 |
| `Custom/Errors/{error_type}/{severity}` | エラー数 |

### NRQL クエリ例

```sql
-- エラー率（直近5分）
SELECT filter(count(*), WHERE numeric(httpResponseCode) >= 500) / count(*) * 100 AS 'Error Rate %'
FROM Transaction WHERE appName = 'financial-planning-calculator' SINCE 5 minutes ago

-- P99 レスポンスタイム
SELECT percentile(duration, 99) AS 'P99 Response Time (s)'
FROM Transaction WHERE appName = 'financial-planning-calculator' SINCE 30 minutes ago TIMESERIES

-- エンドポイント別スループット
SELECT rate(count(*), 1 minute) AS 'Requests/min'
FROM Transaction WHERE appName = 'financial-planning-calculator'
FACET request.uri SINCE 30 minutes ago LIMIT 10
```

---

## アラート設定

アラートは New Relic の **Alert Policies** で設定します。詳細な設定手順と対処方法は [RUNBOOK.md](RUNBOOK.md) を参照してください。

### アラートポリシー作成手順

1. [New Relic](https://one.newrelic.com) にログイン
2. **Alerts** → **Alert Policies** → **New alert policy**
3. ポリシー名: `financial-planning-calculator-alerts`
4. 通知先に Slack を追加:
   - Notification channels → Add notification channel → Slack
   - Webhook URL を設定

### 設定するアラート

| アラート名 | 重篤度 | 条件 |
|-----------|--------|------|
| HighErrorRate | Critical | エラーレート > 5% が5分以上 |
| SlowResponseTime | Warning | P99 > 2秒 |
| ServiceDown | Critical | データが1分以上届かない |
| DatabaseConnectionHigh | Warning | DB接続数 > 80 |
| HighMemoryUsage | Warning | メモリ > 500MB が5分以上 |

---

## 構造化ログ

### ログレベル設定

```env
# .env
LOG_LEVEL=INFO  # DEBUG / INFO / WARN / ERROR
```

### ログ形式

すべてのログは JSON 形式で出力されます：

```json
{
  "time": "2026-02-11T10:00:00Z",
  "level": "ERROR",
  "msg": "エラーが発生しました",
  "request_id": "abc123",
  "user_id": "user456",
  "error": "database connection failed",
  "stack_trace": "..."
}
```

New Relic は Go エージェントのログ転送機能（`ConfigAppLogForwardingEnabled(true)`）を通じてこれらのログを自動収集します。

---

## エラー追跡

### ErrorTracker インターフェース

```go
import "github.com/financial-planning-calculator/backend/infrastructure/monitoring"

// エラーをキャプチャ（New Relic APM に通知）
tags := map[string]string{
    "component": "payment",
    "severity":  "high",
}
monitoring.CaptureError(ctx, err, tags)

// メッセージをキャプチャ
monitoring.CaptureMessage(ctx, "処理完了", "info", tags)
```

エラーは New Relic の **Errors** セクションでスタックトレース付きで確認できます。

---

## ヘルスチェック

```bash
# シンプルなヘルスチェック
GET /health

# 詳細なヘルスチェック（DB接続など）
GET /health/detailed

# 準備状態チェック
GET /ready
```

---

## パフォーマンスプロファイリング（開発環境）

```env
ENABLE_PPROF=true
PPROF_PORT=6060
```

```bash
# メモリプロファイル
go tool pprof http://localhost:6060/debug/pprof/heap

# CPU プロファイル
go tool pprof http://localhost:6060/debug/pprof/profile
```

---

## トラブルシューティング

### New Relic にデータが届かない

1. License Key が正しいか確認:
   ```bash
   docker exec financial_planning_backend env | grep NEW_RELIC
   ```
2. 起動ログで New Relic の初期化状態を確認:
   ```bash
   docker logs financial_planning_backend | grep "New Relic"
   # ✅ New Relic エージェントを初期化しました → 正常
   # ⚠️ New Relic 初期化失敗 → License Key を確認
   ```
3. New Relic にデータが届くまで数分かかる場合があります

### ログが出力されない

```env
LOG_LEVEL=DEBUG  # より詳細なログを出力
```

---

## 参考資料

- [New Relic Go エージェント](https://docs.newrelic.com/docs/apm/agents/go-agent/get-started/introduction-new-relic-go/)
- [New Relic NRQL リファレンス](https://docs.newrelic.com/docs/query-your-data/nrql-new-relic-query-language/get-started/nrql-syntax-clauses-functions/)
- [Go slog パッケージ](https://pkg.go.dev/log/slog)
- [アラート Runbook](RUNBOOK.md)
