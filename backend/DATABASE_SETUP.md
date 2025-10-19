# データベースセットアップガイド

## 概要

財務計画計算機のPostgreSQLデータベースセットアップ手順です。

## 前提条件

### PostgreSQLのインストール

#### macOS (Homebrew)
```bash
brew install postgresql
brew services start postgresql
```

#### Ubuntu/Debian
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

#### CentOS/RHEL
```bash
sudo yum install postgresql-server postgresql-contrib
sudo postgresql-setup initdb
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

### Go環境
- Go 1.21以上がインストールされていること

## セットアップ手順

### 1. 自動セットアップ（推奨）

```bash
cd backend
make db-setup
```

このコマンドは以下を自動実行します：
- PostgreSQLの起動確認
- データベースの作成
- マイグレーションの実行
- サンプルデータの投入（オプション）

### 2. 手動セットアップ

#### Step 1: データベース作成

```bash
# PostgreSQLに接続
psql -U postgres

# データベース作成
CREATE DATABASE financial_planning;

# ユーザー作成（オプション）
CREATE USER financial_app WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE financial_planning TO financial_app;

# 接続終了
\q
```

#### Step 2: 環境変数設定

`.env`ファイルを作成：

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=financial_planning
DB_SSLMODE=disable
```

#### Step 3: マイグレーション実行

```bash
# 依存関係のインストール
make deps

# マイグレーション実行
make migrate-up

# マイグレーション状況確認
make migrate-status
```

#### Step 4: サンプルデータ投入（オプション）

```bash
make seed
```

## 利用可能なコマンド

### Makefileコマンド

```bash
# データベース関連
make migrate-up      # マイグレーション実行
make migrate-down    # マイグレーションロールバック
make migrate-status  # マイグレーション状況確認
make seed           # サンプルデータ投入
make db-reset       # データベースリセット
make db-setup       # 対話式セットアップ

# 開発関連
make build          # アプリケーションビルド
make run            # アプリケーション実行
make test           # テスト実行
make dev-setup      # 開発環境セットアップ
```

### 直接コマンド

```bash
# マイグレーション
go run ./cmd/migrate/main.go -command=up
go run ./cmd/migrate/main.go -command=down
go run ./cmd/migrate/main.go -command=status

# シード
go run ./cmd/seed/main.go
```

## データベース構造

### テーブル一覧

1. **users** - ユーザー情報
2. **financial_data** - 財務データ（収入・投資設定）
3. **expense_items** - 支出項目
4. **savings_items** - 貯蓄項目
5. **retirement_data** - 退職・年金情報
6. **goals** - 財務目標

### サンプルデータ

シードデータには以下が含まれます：

- 2人のサンプルユーザー
- 各ユーザーの財務データ（収入・支出・貯蓄）
- 退職計画データ
- 複数の財務目標（緊急資金、住宅購入、老後資金など）

## トラブルシューティング

### よくある問題

#### 1. PostgreSQLに接続できない

```bash
# PostgreSQLの状態確認
pg_isready -h localhost -p 5432

# サービス起動（macOS）
brew services start postgresql

# サービス起動（Linux）
sudo systemctl start postgresql
```

#### 2. データベースが存在しない

```bash
# データベース一覧確認
psql -U postgres -l

# データベース作成
createdb -U postgres financial_planning
```

#### 3. 権限エラー

```bash
# PostgreSQLユーザーでログイン
sudo -u postgres psql

# 権限付与
GRANT ALL PRIVILEGES ON DATABASE financial_planning TO your_user;
```

#### 4. マイグレーションエラー

```bash
# マイグレーション状況確認
make migrate-status

# 問題のあるマイグレーションをロールバック
make migrate-down

# 再実行
make migrate-up
```

### ログ確認

```bash
# PostgreSQLログ確認（macOS）
tail -f /usr/local/var/log/postgres.log

# PostgreSQLログ確認（Linux）
sudo tail -f /var/log/postgresql/postgresql-*.log
```

## 本番環境での注意事項

### セキュリティ

1. **強力なパスワード**: 本番環境では強力なパスワードを使用
2. **SSL接続**: `DB_SSLMODE=require`に設定
3. **ファイアウォール**: 必要なポートのみ開放
4. **定期バックアップ**: 自動バックアップの設定

### パフォーマンス

1. **接続プール**: 適切な接続プール設定
2. **インデックス**: クエリパフォーマンスの監視
3. **統計情報**: 定期的な`ANALYZE`実行

### バックアップ

```bash
# バックアップ作成
pg_dump -h localhost -U postgres financial_planning > backup_$(date +%Y%m%d_%H%M%S).sql

# リストア
psql -h localhost -U postgres financial_planning < backup_file.sql
```

## 開発者向け情報

### 新しいマイグレーションの作成

1. `infrastructure/database/migrations/`に新しいファイルを作成
2. ファイル名: `{version}_{description}.sql`
3. ロールバック用: `{version}_{description}_down.sql`

例：
```
002_add_user_preferences.sql
002_add_user_preferences_down.sql
```

### テストデータベース

テスト用に別のデータベースを使用：

```env
# テスト環境用 .env.test
DB_NAME=financial_planning_test
```

```bash
# テスト用データベース作成
createdb -U postgres financial_planning_test

# テスト実行
DB_NAME=financial_planning_test make migrate-up
DB_NAME=financial_planning_test go test ./...
```

## サポート

問題が発生した場合は、以下の情報を含めてIssueを作成してください：

1. OS・PostgreSQLバージョン
2. エラーメッセージ
3. 実行したコマンド
4. 環境変数設定（パスワードは除く）