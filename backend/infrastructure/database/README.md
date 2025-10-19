# データベース設計

## 概要

財務計画計算機のPostgreSQLデータベース設計とマイグレーション管理システムです。

## データベース構造

### テーブル一覧

1. **users** - ユーザー情報
2. **financial_data** - 財務データ（収入・投資設定）
3. **expense_items** - 支出項目
4. **savings_items** - 貯蓄項目
5. **retirement_data** - 退職・年金情報
6. **goals** - 財務目標
7. **schema_migrations** - マイグレーション管理（自動作成）

### ER図

```
users (1) ----< financial_data (1) ----< expense_items (*)
  |                                 \----< savings_items (*)
  |
  +-----< retirement_data (1)
  |
  +-----< goals (*)
```

### 主要な制約

- **外部キー制約**: データ整合性を保証
- **CHECK制約**: 
  - 金額は非負値
  - 年齢は1-120の範囲
  - 投資利回りは0-100%の範囲
  - インフレ率は0-50%の範囲
- **UNIQUE制約**: 
  - ユーザーごとに財務データは1件のみ
  - ユーザーごとに退職データは1件のみ
- **NOT NULL制約**: 必須フィールドの保証

### インデックス

パフォーマンス最適化のため以下のインデックスを設定：

- ユーザーID関連の検索用インデックス
- カテゴリ・タイプ別検索用インデックス
- アクティブな目標の検索用複合インデックス
- 目標期日での検索用インデックス

## マイグレーション

### 使用方法

```bash
# 全てのマイグレーションを適用
make migrate-up

# 最新のマイグレーションをロールバック
make migrate-down

# マイグレーション状況を確認
make migrate-status
```

### マイグレーションファイル

- `migrations/001_create_initial_schema.sql` - 初期スキーマ作成
- `migrations/001_create_initial_schema_down.sql` - 初期スキーマロールバック

### 新しいマイグレーションの追加

1. `migrations/` ディレクトリに新しいファイルを作成
2. ファイル名は `{version}_{description}.sql` の形式
3. 対応するロールバックファイル `{version}_{description}_down.sql` も作成
4. `make migrate-up` で適用

例：
```
002_add_user_preferences.sql
002_add_user_preferences_down.sql
```

## シードデータ

### 使用方法

```bash
# サンプルデータを投入
make seed

# データベースをリセット（マイグレーション + シード）
make db-reset
```

### サンプルデータ内容

- 2人のサンプルユーザー
- 各ユーザーの財務データ（収入・支出・貯蓄）
- 退職計画データ
- 複数の財務目標（緊急資金、住宅購入、老後資金など）

## 環境設定

### 環境変数

`.env` ファイルまたは環境変数で以下を設定：

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=financial_planning
DB_SSLMODE=disable
```

### PostgreSQL設定

```sql
-- データベース作成
CREATE DATABASE financial_planning;

-- ユーザー作成（オプション）
CREATE USER financial_app WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE financial_planning TO financial_app;
```

## 開発・テスト

### 開発環境セットアップ

```bash
# 依存関係インストール + マイグレーション + シード
make dev-setup
```

### テストデータベース

テスト用に別のデータベースを使用することを推奨：

```env
# テスト環境用
DB_NAME=financial_planning_test
```

## パフォーマンス考慮事項

### インデックス戦略

- **単一カラムインデックス**: 頻繁に検索されるカラム
- **複合インデックス**: 複数条件での検索用
- **部分インデックス**: 条件付きインデックス（例：アクティブな目標のみ）

### クエリ最適化

- 適切なWHERE句の使用
- JOINの最適化
- 必要なカラムのみのSELECT

### 監視項目

- スロークエリの監視
- インデックス使用状況
- テーブルサイズの監視
- 接続数の監視

## セキュリティ

### データ保護

- 機密データの暗号化（アプリケーション層で実装）
- 適切なアクセス権限設定
- SQLインジェクション対策（パラメータ化クエリ）

### バックアップ

```bash
# データベースバックアップ
pg_dump financial_planning > backup.sql

# リストア
psql financial_planning < backup.sql
```

## トラブルシューティング

### よくある問題

1. **マイグレーション失敗**
   - ロールバックして再実行
   - データベース接続確認

2. **制約違反**
   - データの整合性確認
   - 制約条件の確認

3. **パフォーマンス問題**
   - EXPLAIN ANALYZEでクエリ分析
   - インデックス追加検討

### ログ確認

```bash
# PostgreSQLログ確認
tail -f /var/log/postgresql/postgresql-*.log
```