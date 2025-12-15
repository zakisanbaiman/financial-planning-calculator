# バックエンド ホットリロード 実装検証

## 実装内容

バックエンドのホットリロード機能が実装されました。

### 変更ファイル

1. **docker-compose.yml**
   - `command: ["go", "run", "main.go"]` → `command: ["air", "-c", ".air.toml"]`
   - これによりAirがバックエンドを起動し、ファイル変更を監視します

2. **Makefile**
   - Docker開発用のコマンドを追加（Docker Compose V2対応）
   - 合計150行以上の新しいコマンドを追加

3. **README.md**
   - ホットリロード機能の説明を追加
   - Docker関連コマンドを拡充

4. **DOCKER_SETUP.md**
   - ホットリロードの詳細な説明を追加
   - 開発ワークフローを更新

5. **HOT_RELOAD_GUIDE.md（新規作成）**
   - 包括的なホットリロードガイド
   - セットアップ、動作確認、トラブルシューティング

## 検証方法

### 前提条件

- Docker & Docker Compose V2がインストールされていること
- ポート8080、5432が使用可能であること

### 手動検証手順

#### 1. 環境セットアップ

```bash
# プロジェクトルートに移動
cd /path/to/financial-planning-calculator

# 初回セットアップ
make dev-setup
```

期待される出力：
```
🚀 Setting up Docker development environment...
⏳ Waiting for database to be ready...
📦 Running database migrations...
✅ Migrations complete!
🌱 Seeding database...
✅ Seeding complete!
✅ Setup complete!

Backend API: http://localhost:8080
Swagger UI:  http://localhost:8080/swagger/index.html
Database:    localhost:5432
```

#### 2. ホットリロードの確認

**ターミナル1: ログ監視**
```bash
make logs-api
```

**ターミナル2: コード編集**
```bash
# backend/main.go を編集
vim backend/main.go

# 例: 66行目のログメッセージを変更
# 変更前:
log.Printf("サーバーを開始します: http://localhost:%s", cfg.Port)

# 変更後:
log.Printf("🚀 ホットリロード対応サーバーを開始します: http://localhost:%s", cfg.Port)

# 保存
```

**期待される動作:**
1. ファイル保存後、約1秒以内にターミナル1のログに以下が表示される：
```
main.go has changed
building...
running...
🚀 ホットリロード対応サーバーを開始します: http://localhost:8080
```

2. サーバーが再起動され、新しいログメッセージが表示される

#### 3. API動作確認

```bash
# ヘルスチェック
curl http://localhost:8080/health

# Swagger UI確認
open http://localhost:8080/swagger/index.html
```

期待される結果：
- ヘルスチェックが正常に応答
- Swagger UIが表示される
- コード変更後も正常に動作

#### 4. 追加の動作確認

**複数ファイルの変更:**
```bash
# config/server.go を編集
vim backend/config/server.go

# infrastructure/web/routes.go を編集
vim backend/infrastructure/web/routes.go
```

期待される動作：
- 各ファイル保存後、自動的に再ビルドされる
- エラーがなければサーバーが正常に再起動

**構文エラーのテスト:**
```bash
# 意図的に構文エラーを作成
echo "invalid go code" >> backend/main.go

# ログを確認
make logs-api
```

期待される動作：
- ビルドエラーがログに表示される
- エラー修正後、自動的に再ビルドされる

### 自動テスト（将来的に追加可能）

```bash
# E2Eテストでホットリロードをテスト
./scripts/test-hot-reload.sh
```

## 確認項目チェックリスト

- [ ] `make dev-setup` が正常に完了
- [ ] `make up` でサーバーが起動
- [ ] `.go` ファイル編集後、自動的に再ビルド
- [ ] 再ビルド時間が1-3秒以内
- [ ] 構文エラー時に適切なエラーメッセージ表示
- [ ] エラー修正後、自動的に回復
- [ ] `make logs-api` でログが確認可能
- [ ] API（/health）が正常に応答
- [ ] Swagger UIが表示される
- [ ] `make down` で正常に停止

## トラブルシューティング

### Docker ビルドエラー

現在、Alpine Linuxのリポジトリへのネットワーク接続問題により、
Docker imageのビルドが失敗する可能性があります。

**回避策:**
1. ネットワーク接続を確認
2. Docker Desktopを再起動
3. DNS設定を確認

または、既にビルド済みのイメージがある場合はそれを使用してください。

### ポート競合

ポート8080または5432が既に使用されている場合：

```bash
# ポート使用状況確認
lsof -i :8080
lsof -i :5432

# プロセスを停止してから再実行
```

## 実装の技術詳細

### Air 設定

`backend/.air.toml`:
- 監視対象: `.go`, `.tpl`, `.tmpl`, `.html`
- 除外: テストファイル、マイグレーション、シード
- ビルド遅延: 1000ms
- 出力先: `./tmp/main`

### Docker Compose 設定

```yaml
services:
  backend:
    command: ["air", "-c", ".air.toml"]
    volumes:
      - ./backend:/app
      - go_mod_cache:/go/pkg/mod
```

- Airがコンテナ内で実行される
- ホストの`backend/`ディレクトリがマウントされる
- ファイル変更が即座に検知される

## 参考リンク

- [Air GitHub Repository](https://github.com/air-verse/air)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Go Hot Reload Best Practices](https://threedots.tech/post/go-hot-reload/)

## 次のステップ

1. 本番環境用のDockerfileは`go run`ではなくビルド済みバイナリを使用（既に実装済み）
2. E2Eテストにホットリロードのテストケースを追加
3. CI/CDパイプラインでホットリロードが有効になっていないことを確認
