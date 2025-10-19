# Docker開発環境セットアップガイド

## 概要

財務計画計算機のDocker開発環境です。PostgreSQL、Go API、Next.jsフロントエンドを含む完全なコンテナ化された開発環境を提供します。

## 前提条件

- Docker Desktop（macOS/Windows）または Docker Engine（Linux）
- Docker Compose v2.0以上

### インストール確認

```bash
docker --version
docker-compose --version
```

## クイックスタート

### 1. 開発環境の起動

```bash
# 全体をセットアップ（初回のみ）
make dev-setup

# または段階的に実行
make build    # Dockerイメージをビルド
make up       # バックエンド + データベースを起動
make migrate  # マイグレーション実行
make seed     # サンプルデータ投入
```

### 2. アクセス確認

- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **データベース**: localhost:5432

### 3. 開発開始

```bash
# ログを確認
make logs

# 特定サービスのログ
make logs-api
make logs-db
```

## 利用可能なコマンド

### 基本操作

```bash
make help         # ヘルプ表示
make build        # Dockerイメージをビルド
make up           # 開発環境を起動
make up-full      # フロントエンド含む全サービス起動
make down         # 環境を停止
make clean        # 全削除（データも含む）
```

### データベース操作

```bash
make migrate         # マイグレーション実行
make migrate-status  # マイグレーション状況確認
make migrate-down    # マイグレーションロールバック
make seed           # サンプルデータ投入
make reset          # DB完全リセット
make shell-db       # PostgreSQLに接続
```

### 開発・テスト

```bash
make test           # テスト実行
make test-coverage  # カバレッジ付きテスト
make shell-api      # バックエンドコンテナに接続
make logs          # 全ログ表示
make logs-api      # APIログのみ表示
```

## 開発ワークフロー

### 1. 日常的な開発

```bash
# 朝の作業開始
make up

# コードを編集（ホットリロードで自動反映）
# backend/ 配下のGoファイルを編集

# 夜の作業終了
make down
```

### 2. データベーススキーマ変更

```bash
# 新しいマイグレーションファイルを作成
# backend/infrastructure/database/migrations/002_new_feature.sql

# マイグレーション実行
make migrate

# 確認
make migrate-status
```

### 3. 依存関係の追加

```bash
# コンテナ内でgo mod tidyを実行
make shell-api
go mod tidy
exit

# または
make update-deps
```

## ディレクトリ構造

```
.
├── docker-compose.yml          # メイン設定
├── docker-compose.prod.yml     # 本番用設定
├── Makefile                    # 開発用コマンド
├── .env.docker                 # Docker環境変数
├── backend/
│   ├── Dockerfile              # バックエンドイメージ
│   ├── .air.toml              # ホットリロード設定
│   └── infrastructure/database/
│       ├── migrations/         # マイグレーションファイル
│       ├── seeds/             # シードデータ
│       └── init/              # Docker初期化スクリプト
└── frontend/                   # フロントエンド（将来）
    └── Dockerfile
```

## 環境設定

### 環境変数

Docker環境では`.env.docker`の設定が使用されます：

```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=financial_planning
```

### ポート設定

- **3000**: フロントエンド（Next.js）
- **8080**: バックエンドAPI（Go/Echo）
- **5432**: PostgreSQL

## トラブルシューティング

### よくある問題

#### 1. ポートが既に使用されている

```bash
# 使用中のポートを確認
lsof -i :8080
lsof -i :5432

# 該当プロセスを停止してから再実行
make down
make up
```

#### 2. データベース接続エラー

```bash
# PostgreSQLの状態確認
docker-compose ps postgres

# ログ確認
make logs-db

# ヘルスチェック確認
docker-compose exec postgres pg_isready -U postgres
```

#### 3. マイグレーションエラー

```bash
# マイグレーション状況確認
make migrate-status

# 問題があればロールバック
make migrate-down

# データベースを完全リセット
make clean
make dev-setup
```

#### 4. ビルドエラー

```bash
# キャッシュをクリアして再ビルド
docker-compose build --no-cache

# または完全クリーンアップ
make clean
docker system prune -a
make build
```

### ログ確認

```bash
# 全サービスのログ
make logs

# 特定サービス
docker-compose logs -f backend
docker-compose logs -f postgres

# エラーのみ
docker-compose logs --tail=50 backend | grep -i error
```

## 本番環境

### 本番用ビルド

```bash
# 本番用イメージをビルド
make build-prod

# 本番環境で起動
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### セキュリティ設定

本番環境では以下を変更してください：

1. **パスワード**: `.env.docker`のパスワードを変更
2. **SSL**: HTTPS設定を有効化
3. **ファイアウォール**: 必要なポートのみ開放

## パフォーマンス最適化

### 開発環境

- **ホットリロード**: Air使用でコード変更を自動反映
- **ボリュームマウント**: ローカルファイルを直接マウント
- **Go Modules キャッシュ**: 依存関係のダウンロード時間短縮

### 本番環境

- **マルチステージビルド**: 最小限のイメージサイズ
- **Alpine Linux**: 軽量ベースイメージ
- **PostgreSQL最適化**: 本番用設定を適用

## 開発のベストプラクティス

### 1. データの永続化

```bash
# 開発データを保持したい場合
make down  # コンテナ停止（データは保持）

# 完全にクリーンアップしたい場合
make clean  # データも削除
```

### 2. デバッグ

```bash
# コンテナ内でデバッグ
make shell-api
dlv debug --headless --listen=:2345 --api-version=2

# ログレベル調整
# .env.dockerでGIN_MODE=debugに設定済み
```

### 3. テスト

```bash
# 単体テスト
make test

# 統合テスト（データベース含む）
make test-coverage

# 特定パッケージのテスト
docker-compose run --rm backend go test ./domain/... -v
```

## サポート

問題が発生した場合：

1. `make logs`でログを確認
2. `docker-compose ps`でサービス状態を確認
3. 必要に応じて`make clean && make dev-setup`で環境をリセット