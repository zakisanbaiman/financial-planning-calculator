# ADR 006: Redisキャッシュ戦略の採用（Cache-Aside + Repositoryデコレータ）

## ステータス

採択済み (2026-04-02)

## 背景

Issue #214「Redisキャッシュ戦略の設計・実装」では、現在レートリミット用途にしか使われていない Redis をキャッシュ層として活用し、財務計算のパフォーマンスを向上させることが目標とされている。

本プロジェクトは Go/Echo バックエンドと DDD アーキテクチャを採用しており、以下の状況がある。

- `backend/infrastructure/redis/client.go` に Redis クライアントが実装済み
- `FinancialPlanRepository` と `GoalRepository` が財務計算ユースケースで頻繁に呼ばれる読み取り操作を提供している
- PostgreSQL リポジトリはインターフェース経由で利用されており、デコレータパターンが適用しやすい構造になっている

## 問題

財務計算ユースケース（`CalculateComprehensiveProjection` など）は、同一ユーザーの財務計画と目標データを繰り返し読み取る。これらのデータは更新頻度が低く、読み取り頻度が高いため、キャッシュによる効果が期待できる。

キャッシュ戦略の選択肢として以下を検討した。

| パターン | 概要 |
|---|---|
| Cache-Aside（Lazy Loading） | 読み取り時にキャッシュを確認し、なければDBから取得してキャッシュに保存 |
| Write-Through | 書き込み時に同時にキャッシュへも書き込む |
| Read-Through | キャッシュレイヤー自体がDBからのロードを担う（ライブラリ依存） |
| インメモリキャッシュ | アプリケーション内のマップで保持 |

## 決定

**Cache-Aside パターン + Repository デコレータ**を採用する。

### キャッシュ対象と TTL

| リポジトリメソッド | キャッシュキー | TTL |
|---|---|---|
| `FinancialPlanRepository.FindByID` | `fp:plan:id:{planID}` | 5分 |
| `FinancialPlanRepository.FindByUserID` | `fp:plan:uid:{userID}` | 5分 |
| `GoalRepository.FindByUserID` | `fp:goals:uid:{userID}` | 3分 |
| `GoalRepository.FindActiveGoalsByUserID` | `fp:goals:active:uid:{userID}` | 3分 |

### 無効化戦略（Write-Invalidate）

- `FinancialPlan` の Save/Update 後に `fp:plan:id:{planID}` と `fp:plan:uid:{userID}` を削除
- `Goal` の Save/Update 後に `fp:goals:uid:{userID}` と `fp:goals:active:uid:{userID}` を削除
- 削除失敗はログ出力のみ（TTL が切れれば一貫性は自然に回復する）

### 実装方法

1. `CacheClient` インターフェース（`SetJSON/GetJSON/Delete/DeleteByPattern`）を定義し、具体的な Redis 実装と分離する（テスタビリティ確保）
2. `CachedFinancialPlanRepository` と `CachedGoalRepository` がデコレータとして既存の PostgreSQL 実装をラップする
3. ドメインオブジェクトは非公開フィールドを持つため、キャッシュ専用 DTO（`financialPlanCacheDTO`, `goalCacheDTO`）を使って JSON シリアライズする
4. `main.go` で Redis 接続確認後、接続成功時のみデコレータで置き換える

### fail-open（Redis障害時の動作）

- `redis.Nil`（キャッシュミス）と接続エラーを区別する
- 接続エラー時は `slog.Warn` でログを出力し、PostgreSQL にフォールバックする
- Redis 起動失敗時はキャッシュなしで起動継続する

### 監視

`cache_hits_total` と `cache_misses_total` の Prometheus CounterVec で記録。Grafana の recording rule で `cache_hits_total / (cache_hits_total + cache_misses_total)` としてヒット率を算出する。

## 代替案と却下理由

**Write-Through**: 書き込みパスが複雑になり、トランザクション境界の扱いが難しくなる。読み取り頻度に比べて書き込みが少ないため、Cache-Aside の方が実装コストと効果のバランスが良い。

**Read-Through（ライブラリ利用）**: 外部依存を増やさない方針のため却下。DDD アーキテクチャとの統合も複雑になる。

**インメモリキャッシュ**: 複数インスタンス間でキャッシュが共有されず、データ不整合が発生しやすい。Railway でのスケールアウト時に問題になる。

## 影響

- 新規ファイル: `infrastructure/redis/cache_client.go`, `infrastructure/redis/cache.go`, `infrastructure/repositories/cache_keys.go`, `infrastructure/repositories/cache_dto.go`, `infrastructure/repositories/cached_financial_plan_repository.go`, `infrastructure/repositories/cached_goal_repository.go`
- ドメイン層への追加: `entities.NewFinancialProfileWithID`, `entities.NewRetirementDataWithID`, `aggregates.NewFinancialPlanWithID`（リポジトリ復元用ファクトリ関数）
- `main.go` に Redis 初期化とデコレータ注入を追加
- `infrastructure/monitoring/prometheus.go` に `CacheHitsTotal`, `CacheMissesTotal` を追加
