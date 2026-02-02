# Dev Container セットアップガイド

このプロジェクトは [Dev Containers](https://containers.dev/) をサポートしています。Dev Containersを使用することで、一貫性のある開発環境をすばやくセットアップできます。

## 必要なもの

- [Visual Studio Code](https://code.visualstudio.com/)
- [Dev Containers 拡張機能](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
- [Docker Desktop](https://www.docker.com/products/docker-desktop)

## セットアップ手順

### 1. リポジトリをクローン

```bash
git clone https://github.com/zakisanbaiman/financial-planning-calculator.git
cd financial-planning-calculator
```

### 2. Dev Containerで開く

VS Codeでプロジェクトを開き、以下のいずれかの方法でDev Containerを起動します：

**方法1: コマンドパレットから**
1. `Ctrl+Shift+P` (Mac: `Cmd+Shift+P`) でコマンドパレットを開く
2. `Dev Containers: Reopen in Container` を選択

**方法2: 通知から**
1. VS Codeでフォルダーを開く
2. 右下に表示される通知の「Reopen in Container」をクリック

### 3. 自動セットアップ

Dev Containerが起動すると、以下が自動的に実行されます：

- Go 1.24とNode.js 18のインストール
- 必要なVS Code拡張機能のインストール
- Go toolsのインストール (air, golangci-lint, swag)
- 依存関係のインストール (Go modules, npm packages)
- Git hooksのセットアップ
- データベースのマイグレーションとシード

初回起動時は5〜10分程度かかる場合があります。

## 使い方

### 開発サーバーの起動

```bash
# バックエンド + データベース（ホットリロード有効）
make up

# すべてのサービス（フロントエンド含む）
make up-full

# ログを確認
make logs
```

### サービスへのアクセス

- **Backend API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Frontend**: http://localhost:3000
- **pprof**: http://localhost:6060/debug/pprof/

### 便利なコマンド

```bash
# テスト実行
make test

# Lint実行
make lint

# データベース操作
make migrate        # マイグレーション実行
make seed          # シードデータ投入
make shell-db      # PostgreSQLシェルを開く

# コンテナ操作
make down          # サービス停止
make logs          # すべてのログ表示
make logs-api      # バックエンドログのみ表示

# ヘルプ
make help          # すべてのコマンドを表示
make docker-help   # Docker関連コマンドを表示
```

## インストールされる拡張機能

Dev Containerには以下のVS Code拡張機能が自動的にインストールされます：

### 開発ツール
- **Go** - Go言語サポート
- **ESLint** - JavaScriptのLint
- **Prettier** - コードフォーマッター
- **Tailwind CSS IntelliSense** - Tailwind CSSのサポート

### 生産性向上
- **GitLens** - Git機能拡張
- **GitHub Copilot** - AIコード補完
- **Path Intellisense** - パス補完
- **Resource Monitor** - リソース監視

### ドキュメント・API
- **Swagger** - OpenAPI/Swaggerサポート
- **Markdown All in One** - Markdown編集

## 設定のカスタマイズ

Dev Containerの設定をカスタマイズするには、`.devcontainer/devcontainer.json`を編集します。

### 拡張機能の追加

```json
{
  "customizations": {
    "vscode": {
      "extensions": [
        "your.extension.id"
      ]
    }
  }
}
```

### ポート転送の追加

```json
{
  "forwardPorts": [9000],
  "portsAttributes": {
    "9000": {
      "label": "My Service",
      "onAutoForward": "notify"
    }
  }
}
```

## トラブルシューティング

### Dev Containerが起動しない

1. Dockerが起動していることを確認
2. VS Codeのターミナルでエラーメッセージを確認
3. `.devcontainer/devcontainer.json`の構文エラーをチェック

### データベースに接続できない

```bash
# PostgreSQLの状態を確認
docker ps | grep postgres

# データベースコンテナを再起動
make down
make up
```

### 依存関係のインストールに失敗

```bash
# Dev Containerを再ビルド
# コマンドパレット (Ctrl+Shift+P) から
# "Dev Containers: Rebuild Container" を実行
```

### ホットリロードが動作しない

```bash
# コンテナを再起動
make restart

# ログを確認
make logs-api
```

## 既存のDocker環境との違い

Dev Containerは既存の `docker-compose.yml` を利用しますが、以下の点が異なります：

| 項目 | 通常のDocker環境 | Dev Container |
|------|-----------------|--------------|
| 用途 | 開発サーバーの実行 | VS Code統合開発環境 |
| エディタ | ローカルのエディタ | VS Code（コンテナ内） |
| 拡張機能 | ローカルにインストール | コンテナ内にインストール |
| ファイル編集 | ローカルファイルシステム | コンテナ内でマウント |
| セットアップ | 手動 | 自動（post-createスクリプト） |

両方の環境を並行して使用できます。

## 参考リンク

- [Dev Containers公式ドキュメント](https://containers.dev/)
- [VS Code Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers)
- [Dev Container仕様](https://containers.dev/implementors/spec/)
