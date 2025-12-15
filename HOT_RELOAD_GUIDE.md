# バックエンド ホットリロード ガイド

## 概要

このプロジェクトのバックエンドは、[Air](https://github.com/air-verse/air)を使用したホットリロード機能を搭載しています。コードを編集すると、自動的に再ビルド・再起動され、開発効率が大幅に向上します。

## セットアップ

### 前提条件
- Docker & Docker Compose V2
- Make

### 起動方法

```bash
# 初回セットアップ
make dev-setup

# 2回目以降
make up

# ログを表示（ホットリロードの動作確認）
make logs-api
```

## ホットリロードの動作確認

### 1. バックエンドサーバーを起動

```bash
make up
```

### 2. ログを監視

別のターミナルウィンドウで：

```bash
make logs-api
```

### 3. コードを編集

例えば、`backend/main.go`を編集：

```go
// 変更前
log.Printf("サーバーを開始します: http://localhost:%s", cfg.Port)

// 変更後
log.Printf("🚀 サーバーを開始します: http://localhost:%s", cfg.Port)
```

保存すると、ログに以下のような表示が出ます：

```
building...
running...
🚀 サーバーを開始します: http://localhost:8080
```

### 4. APIの動作確認

```bash
curl http://localhost:8080/health
```

## Air の設定

設定ファイル：`backend/.air.toml`

### 主な設定項目

```toml
[build]
  # ビルドコマンド
  cmd = "go build -o ./tmp/main ."
  
  # 出力先
  bin = "./tmp/main"
  
  # 変更検知後の遅延時間（ミリ秒）
  delay = 1000
  
  # 監視対象の拡張子
  include_ext = ["go", "tpl", "tmpl", "html"]
  
  # 除外するディレクトリ
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "docs", 
                 "infrastructure/database/migrations", 
                 "infrastructure/database/seeds"]
  
  # 除外するファイルパターン
  exclude_regex = ["_test.go"]
```

## トラブルシューティング

### 1. ホットリロードが動作しない

**原因：** Airが正しく起動していない

**解決策：**
```bash
# コンテナの状態を確認
docker compose ps

# ログを確認
make logs-api

# コンテナを再起動
make down
make up
```

### 2. ビルドエラーが発生する

**原因：** コードに構文エラーがある

**解決策：**
```bash
# ログでエラー内容を確認
make logs-api

# エラーを修正後、自動的に再ビルドされる
```

### 3. 変更が反映されない

**原因：** ファイルの保存が正しくされていない、または除外対象のファイル

**解決策：**
- ファイルが正しく保存されているか確認
- `.air.toml`の`include_ext`や`exclude_regex`を確認
- テストファイル（`*_test.go`）は自動的に除外される

### 4. メモリ使用量が増える

**原因：** 頻繁な再ビルドによる一時ファイルの蓄積

**解決策：**
```bash
# コンテナを再起動
make down
make up

# または完全にクリーンアップ
make clean-docker
make dev-setup
```

## パフォーマンス

### ビルド時間

- 初回ビルド：約5-10秒
- ホットリロード時：約1-3秒

### リソース使用量

- CPU：中程度（ビルド時のみ）
- メモリ：約200-300MB（通常時）
- ディスク：約100MB（一時ファイル）

## ベストプラクティス

### 1. 効率的な開発

```bash
# ターミナル1: サーバー起動
make up

# ターミナル2: ログ監視
make logs-api

# ターミナル3: コード編集
vim backend/main.go

# ターミナル4: APIテスト
curl http://localhost:8080/api/...
```

### 2. デバッグ

ホットリロード中でも、通常のデバッグ手法が使用可能：

```go
// ログ出力
log.Printf("Debug: %+v", someVariable)

// pprofによるプロファイリング
// http://localhost:6060/debug/pprof/ にアクセス
```

### 3. テスト駆動開発

ホットリロードはテストファイルを除外しているため、テストは別途実行：

```bash
# Docker内でテスト実行
make test-docker

# または、ローカルで
cd backend && go test ./...
```

## まとめ

- ✅ `.go`ファイルの変更を自動検知
- ✅ 約1秒で自動再ビルド・再起動
- ✅ 開発効率が大幅に向上
- ✅ テストファイルは除外されるため、不要な再起動なし
- ✅ `backend/.air.toml`で柔軟にカスタマイズ可能

## 関連リンク

- [Air GitHub](https://github.com/air-verse/air)
- [Air 設定ドキュメント](https://github.com/air-verse/air#configuration)
- [Docker Development Guide](./DOCKER_SETUP.md)
