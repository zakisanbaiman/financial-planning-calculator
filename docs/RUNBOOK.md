# アラート Runbook

financial-planning-calculator のアラート発生時の対処手順書です。

## アラート一覧

| アラート名 | 重篤度 | 条件 |
|-----------|--------|------|
| [HighErrorRate](#highErrorRate) | Critical | エラーレート > 5% が5分以上 |
| [SlowResponseTime](#slowResponseTime) | Warning | P99レスポンス > 2秒 |
| [DatabaseConnectionHigh](#databaseConnectionHigh) | Warning | DB接続数 > 80% |
| [ServiceDown](#serviceDown) | Critical | ヘルスチェック失敗 |
| [HighMemoryUsage](#highMemoryUsage) | Warning | メモリ使用率 > 85% |

---

## New Relic でのアラート設定手順

1. [New Relic](https://one.newrelic.com) にログイン
2. **Alerts** → **Alert Policies** → **New alert policy** を作成
3. ポリシー名: `financial-planning-calculator-alerts`
4. 各アラートを **Add condition** で追加（下記の各セクションを参照）
5. **Notification channels** に Slack または メールを設定

---

## HighErrorRate

**重篤度:** Critical  
**条件:** HTTPエラーレートが5%を超えた状態が5分以上継続

### New Relic Alert 設定

- **Condition type**: NRQL
- **Query**:
  ```sql
  SELECT filter(count(*), WHERE numeric(httpResponseCode) >= 500) / count(*) * 100 
  FROM Transaction 
  WHERE appName = 'financial-planning-calculator'
  ```
- **Threshold**: `> 5` for `5 minutes`

### 原因調査

```bash
# 1. New Relic APM でエラーの詳細を確認
#    APM → financial-planning-calculator → Errors → Error analytics

# 2. ローカル/本番ログを確認
docker logs financial_planning_backend --since 10m | grep '"level":"ERROR"'

# 3. 特定のエンドポイントでエラーが集中していないか確認
#    New Relic → APM → Transactions → Sort by Error rate
```

### 対処手順

1. **New Relic APM** でエラーの発生しているエンドポイントを特定
2. エラーの種類を確認:
   - `500 Internal Server Error`: バックエンドのバグまたはデータベースの問題
   - `502/503/504`: インフラの問題（メモリ不足、接続タイムアウト）
3. データベース接続エラーの場合は [DatabaseConnectionHigh](#databaseConnectionHigh) も確認
4. コードのバグであれば、原因となるコミットを特定して Revert を検討
5. 解決しない場合はサービスを再起動:
   ```bash
   docker compose restart backend
   ```

---

## SlowResponseTime

**重篤度:** Warning  
**条件:** P99レスポンスタイムが2秒を超えた状態が2分以上継続

### New Relic Alert 設定

- **Condition type**: NRQL
- **Query**:
  ```sql
  SELECT percentile(duration, 99) 
  FROM Transaction 
  WHERE appName = 'financial-planning-calculator'
  ```
- **Threshold**: `> 2` for `2 minutes`

### 原因調査

```bash
# 1. New Relic APM でスロートランザクションを確認
#    APM → financial-planning-calculator → Transactions → Sort by Response time (P99)

# 2. データベースのスロークエリを確認
#    New Relic → APM → Databases → Sort by Average response time

# 3. 外部サービスの応答時間を確認
#    New Relic → APM → External services
```

### 対処手順

1. 遅いエンドポイントを New Relic のトランザクション一覧で特定
2. **データベースが原因の場合**:
   - スロークエリの特定と INDEX 追加を検討
   - DB 接続プールの設定を確認
3. **外部サービスが原因の場合**:
   - タイムアウト設定を確認
   - サーキットブレーカーの導入を検討
4. **N+1 クエリが疑われる場合**:
   - New Relic のトランザクショントレースでクエリ数を確認

---

## DatabaseConnectionHigh

**重篤度:** Warning  
**条件:** アクティブなDB接続数が80を超えた状態が2分以上継続（最大100想定）

### New Relic Alert 設定

- **Condition type**: NRQL
- **Query**:
  ```sql
  SELECT latest(numeric(newrelic.timeslice.value)) 
  FROM Metric 
  WHERE metricTimesliceName = 'Custom/Database/ActiveConnections' 
  AND appName = 'financial-planning-calculator'
  ```
- **Threshold**: `> 80` for `2 minutes`

### 原因調査

```bash
# 1. PostgreSQL で現在の接続数を確認
docker exec financial_planning_db psql -U postgres -c "SELECT count(*) FROM pg_stat_activity;"

# 2. 長時間実行中のクエリを確認
docker exec financial_planning_db psql -U postgres -c "
SELECT pid, now() - pg_stat_activity.query_start AS duration, query, state
FROM pg_stat_activity
WHERE (now() - pg_stat_activity.query_start) > interval '1 minutes'
ORDER BY duration DESC;"

# 3. アイドル接続を確認
docker exec financial_planning_db psql -U postgres -c "
SELECT count(*), state FROM pg_stat_activity GROUP BY state;"
```

### 対処手順

1. 長時間の idle 接続がある場合、接続プールの設定（`max_open_conns`, `max_idle_conns`）を見直す
2. スロークエリが原因で接続が詰まっている場合は [SlowResponseTime](#slowResponseTime) も確認
3. 緊急の場合は idle 接続を強制切断:
   ```sql
   SELECT pg_terminate_backend(pid) 
   FROM pg_stat_activity 
   WHERE state = 'idle' AND query_start < now() - interval '5 minutes';
   ```
4. 根本原因が解決しない場合はバックエンドを再起動してコネクションプールをリセット

---

## ServiceDown

**重篤度:** Critical  
**条件:** ヘルスチェックが1分以上失敗

### New Relic Alert 設定

- **Condition type**: NRQL
- **Query**:
  ```sql
  SELECT count(*) 
  FROM Transaction 
  WHERE appName = 'financial-planning-calculator'
  ```
- **Threshold**: `< 1` for `1 minute` (データが届かない = サービスダウン)

または **Synthetics Monitor** でヘルスチェックを設定:
- URL: `https://your-domain/health`
- 間隔: 1分
- アラート: 連続2回失敗

### 原因調査

```bash
# 1. コンテナの状態確認
docker compose ps

# 2. バックエンドのログを確認
docker logs financial_planning_backend --since 5m

# 3. ヘルスチェックを手動で実行
curl -v http://localhost:8080/health

# 4. DB への接続確認
curl -v http://localhost:8080/health/detailed
```

### 対処手順

1. **コンテナが停止している場合**:
   ```bash
   docker compose start backend
   # または完全に再起動
   docker compose restart backend
   ```
2. **コンテナは起動しているがヘルスチェック失敗の場合**:
   - OOM Kill の確認: `docker inspect financial_planning_backend | grep OOMKilled`
   - メモリ不足の場合は [HighMemoryUsage](#highMemoryUsage) を参照
3. **データベース接続失敗の場合**:
   ```bash
   docker compose restart postgres
   docker compose restart backend
   ```
4. 本番環境では、ロードバランサーが自動でトラフィックを切り替えるか確認

---

## HighMemoryUsage

**重篤度:** Warning  
**条件:** プロセスのメモリ使用量が500MiB（約85%相当）を超えた状態が5分以上継続

### New Relic Alert 設定

- **Condition type**: NRQL
- **Query**:
  ```sql
  SELECT latest(numeric(newrelic.timeslice.value)) 
  FROM Metric 
  WHERE metricTimesliceName LIKE 'Memory/Physical%' 
  AND appName = 'financial-planning-calculator'
  ```
- **Threshold**: `> 500` (MB) for `5 minutes`

### 原因調査

```bash
# 1. コンテナのメモリ使用量を確認
docker stats financial_planning_backend --no-stream

# 2. Go の pprof でメモリプロファイルを取得（ENABLE_PPROF=true の場合）
curl -s http://localhost:6060/debug/pprof/heap > heap.prof
go tool pprof -http=:8888 heap.prof

# 3. New Relic APM でメモリ使用量の推移を確認
#    APM → financial-planning-calculator → JVM/Runtime (Go Runtime)
```

### 対処手順

1. **一時的なメモリスパイクの場合** (トラフィック増加など):
   - しばらく監視して自然に下がるか確認
   - GC が実行されれば回復するはず
2. **メモリリークが疑われる場合**:
   - pprof でヒープを分析してリークの原因を特定
   - 原因となるコミットを特定して修正
3. **緊急の場合**:
   ```bash
   # バックエンドを再起動（メモリをリセット）
   docker compose restart backend
   ```
4. コンテナのメモリ制限が小さすぎる場合は `docker-compose.yml` の `mem_limit` を調整
