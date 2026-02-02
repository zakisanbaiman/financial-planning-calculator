# 財務計画計算機

将来の資産形成と老後の財務計画を可視化するWebアプリケーション

## 概要

このアプリケーションは、ユーザーが現在の収入、支出、貯蓄状況を入力することで、将来の資産推移、老後資金、緊急時の備えを計算し、安心できる財務計画を立てられるようにします。

## 🚀 Dev環境（プレビュー環境）

PRを作成すると、自動的にクラウド上で確認できるプレビュー環境がデプロイされます。

- **プラットフォーム**: Render.com
- **自動デプロイ**: PRの作成・更新時
- **有効期間**: PRクローズ後7日間
- **詳細**: [Dev環境セットアップガイド](./DEV_ENVIRONMENT_SETUP.md)

PRを作成すると、GitHub Actions が自動的にプレビューURLをコメントします。

### 🤖 AI連携によるエラー検知・自動修正

Render.comへのデプロイ時のエラーを自動検知し、AIアシスタント（Claude、Copilot、Cursor）が修正を支援します。

- **MCP (Model Context Protocol)連携**: AIアシスタントがRender.comのデプロイ状態を直接監視
- **自動エラー検知**: デプロイログから一般的なエラーパターンを自動検出
- **PRへの自動通知**: エラー発生時にPRへコメントを追加
- **詳細**: [MCP セットアップガイド](./docs/MCP_SETUP.md) | [クイックリファレンス](./docs/MCP_QUICK_REFERENCE.md)

## 技術スタック

### フロントエンド
- Next.js 14
- TypeScript
- Tailwind CSS
- Chart.js
- React Hook Form + Zod

### バックエンド
- Go
- Echo Framework
- PostgreSQL
- OpenAPI/Swagger

## プロジェクト構造

```
financial-planning-calculator/
├── frontend/          # Next.jsフロントエンド
│   ├── src/
│   │   ├── app/       # App Router
│   │   ├── components/
│   │   ├── lib/
│   │   └── types/
│   ├── package.json
│   ├── next.config.js
│   └── tailwind.config.js
├── backend/           # Goバックエンド
│   ├── config/        # 設定
│   ├── docs/          # OpenAPI仕様
│   ├── go.mod
│   └── main.go
└── README.md
```

## セットアップ

### 🚀 Dev Container（最も簡単な方法）

[Dev Containers](https://containers.dev/)を使用すると、VS Code上で一貫性のある開発環境をすばやくセットアップできます。

**必要なもの:**
- [Visual Studio Code](https://code.visualstudio.com/)
- [Dev Containers拡張機能](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
- [Docker Desktop](https://www.docker.com/products/docker-desktop)

**使い方:**
1. VS Codeでこのプロジェクトを開く
2. コマンドパレット（`Ctrl+Shift+P` / `Cmd+Shift+P`）を開く
3. `Dev Containers: Reopen in Container` を選択
4. 自動セットアップが完了するまで待つ（初回は5〜10分）

詳細は [Dev Containerガイド](.devcontainer/README.md) をご覧ください。

### 前提条件（ローカル開発の場合）
- Node.js 18.20.0+
- Go 1.24.0+
- PostgreSQL 13+
- golangci-lint 1.64+ (開発時)

### バージョン管理

このプロジェクトでは以下のバージョン管理ツールをサポートしています：

**direnv（推奨）- Docker版Goを自動使用**
```bash
# direnvをインストール
brew install direnv

# シェルに統合（.zshrcまたは.bashrcに追加）
eval "$(direnv hook zsh)"  # zshの場合
eval "$(direnv hook bash)" # bashの場合

# プロジェクトディレクトリで許可
cd financial-planning-calculator
direnv allow

# これで'go'コマンドが自動的にDocker内で実行されます！
go version  # 🐳 Docker内のGo 1.24.10を使用
```

**goenv（Go）**
```bash
# .go-versionファイルが自動的に読み込まれます
cd backend
goenv install 1.24.0
```

**mise（Go + Node.js）- asdfの後継、高速・高機能**
```bash
# .mise.tomlファイルが自動的に読み込まれます
# miseのインストール
brew install mise

# シェルに統合（.zshrcまたは.bashrcに追加）
echo 'eval "$(mise activate zsh)"' >> ~/.zshrc  # zshの場合
echo 'eval "$(mise activate bash)"' >> ~/.bashrc # bashの場合

# プロジェクトディレクトリでツールをインストール
cd financial-planning-calculator
mise install
```

**手動インストール**
```bash
# Homebrewの場合
brew install go@1.24
brew install node@18
```

### Docker開発環境（推奨）

```bash
# 初回セットアップ（ビルド + DB起動 + マイグレーション + シード）
make dev-setup

# 2回目以降の起動
make up

# 停止
make down

# その他のコマンド
make help
make docker-help  # Docker関連のコマンド一覧
```

**Docker環境で使えるコマンド:**
- `make test` - テスト実行
- `make lint` - Lint実行
- `make shell-api` - バックエンドコンテナに接続
- `make shell-db` - データベースに接続
- `make logs` - ログ表示
- `make logs-api` - バックエンドログのみ表示
- `make migrate` - DBマイグレーション実行
- `make seed` - サンプルデータ投入

**ホットリロード機能:**
- バックエンドコードは[Air](https://github.com/air-verse/air)により自動的にホットリロードされます
- `.go`ファイルを編集すると自動的に再ビルド・再起動されます
- フロントエンドも同様にNext.jsの開発サーバーで自動リロードされます

### ローカル開発環境

**フロントエンド**
```bash
cd frontend
npm install
cp .env.example .env.local
npm run dev
```

**バックエンド**
```bash
cd backend
go mod tidy
cp .env.example .env
go run main.go
```

### API仕様

Swagger UIは http://localhost:8080/swagger/index.html で確認できます。

### パフォーマンスプロファイリング（pprof）

開発環境ではpprofが有効化されており、パフォーマンス分析が可能です。

**pprofの使い方:**

当プロジェクトディレクトリ内のGoはDocker内のものを利用するため、pprofを実行する場合は他ディレクトリにで実行する。

```bash
# ブラウザでプロファイル一覧を確認
open http://localhost:6060/debug/pprof/

# CPU プロファイル（30秒間）
go tool pprof 'http://localhost:6060/debug/pprof/profile?seconds=30'

# メモリプロファイル
go tool pprof 'http://localhost:6060/debug/pprof/heap'

# ゴルーチン
go tool pprof 'http://localhost:6060/debug/pprof/goroutine'

# インタラクティブモードで分析
# pprofコマンド内で使えるコマンド:
# - top: 上位の関数を表示
# - list <関数名>: 関数の詳細を表示
# - web: グラフをブラウザで表示（graphviz必要）
```

**可視化ツール:**
```bash
# graphvizをインストール（グラフ表示用）
brew install graphviz

# CPUプロファイルをグラフで表示
go tool pprof -http=:8081 'http://localhost:6060/debug/pprof/profile?seconds=30'
```

**本番環境での注意:**
- pprofは開発環境のみで有効（`ENABLE_PPROF=true`）
- 本番環境では必ず無効化すること（セキュリティリスク）

## 開発

### フロントエンド開発
- `npm run dev` - 開発サーバー起動
- `npm run build` - プロダクションビルド
- `npm run lint` - ESLint実行
- `npm run type-check` - TypeScript型チェック

### バックエンド開発
- `go run main.go` - サーバー起動
- `go test ./...` - テスト実行
- `go mod tidy` - 依存関係整理

## 機能

### 実装予定機能
- [ ] 財務データ入力・管理
- [ ] 資産推移シミュレーション
- [ ] 老後資金計算
- [ ] 緊急資金計算
- [ ] 目標設定・進捗管理
- [ ] データ可視化（グラフ・チャート）
- [ ] PDFレポート生成

## ライセンス

MIT License