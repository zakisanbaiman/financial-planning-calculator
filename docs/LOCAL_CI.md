# ローカルでCIを実行する方法

GitHub ActionsのCIワークフローをローカルで実行する方法を説明します。

## 方法1: Makefileコマンド（推奨）

最も簡単な方法です。各CIワークフローを個別に実行できます。

### すべてのCIチェックを実行

```bash
make ci
```

### 個別のワークフローを実行

```bash
# Lintチェック（バックエンド + フロントエンド）
make ci-lint

# テスト（バックエンド + フロントエンドビルド）
make ci-test

# PRチェック（クイックテスト）
make ci-pr-check

# E2Eテスト（データベースとサーバーが必要）
make ci-e2e

# すべてのワークフロー（E2E除く）
make ci-all

# クイックチェック（lint + pr-check）
make ci-quick
```

## 方法2: スクリプトを使用

```bash
# すべてのワークフローを実行（E2E除く）
./scripts/run-ci-local.sh all

# 特定のワークフローを実行
./scripts/run-ci-local.sh lint
./scripts/run-ci-local.sh test
./scripts/run-ci-local.sh pr-check
./scripts/run-ci-local.sh e2e
./scripts/run-ci-local.sh quick
```

## 方法3: actを使用（GitHub Actionsをそのまま実行）

`act`を使用すると、GitHub Actionsのワークフローファイルをそのままローカルで実行できます。

### インストール

```bash
# macOS
brew install act

# または公式サイトから
# https://github.com/nektos/act
```

### 使用方法

```bash
# すべてのワークフローをリスト表示
act -l

# 特定のワークフローを実行
act -W .github/workflows/lint.yml
act -W .github/workflows/test.yml
act -W .github/workflows/pr-check.yml
act -W .github/workflows/e2e-tests.yml

# 特定のジョブを実行
act -j golangci-lint -W .github/workflows/lint.yml
act -j test-backend -W .github/workflows/test.yml

# プルリクエストイベントで実行
act pull_request -W .github/workflows/pr-check.yml
```

### 注意事項

- `act`はDockerが必要です
- E2Eテストはサービス（PostgreSQL）の設定が必要です
- 一部のアクション（例: `actions/cache`）は完全には再現されない場合があります

## 各ワークフローの内容

### Lintワークフロー (`.github/workflows/lint.yml`)

- **バックエンド**: golangci-lint、go fmt、go vet
- **フロントエンド**: TypeScript型チェック、ESLint

### Testワークフロー (`.github/workflows/test.yml`)

- **バックエンド**: ビルド、テスト（race detector付き）、カバレッジ
- **フロントエンド**: ビルド

### PR Checkワークフロー (`.github/workflows/pr-check.yml`)

- **バックエンド**: go vet、クイックテスト（-shortフラグ付き）

### E2E Testsワークフロー (`.github/workflows/e2e-tests.yml`)

- PostgreSQLデータベースのセットアップ
- バックエンドサーバーの起動
- フロントエンドサーバーの起動
- PlaywrightによるE2Eテスト

## トラブルシューティング

### golangci-lintが見つからない

```bash
# macOS
brew install golangci-lint

# またはGoでインストール
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55
```

### E2Eテストが失敗する

E2Eテストを実行するには、データベースとサーバーが起動している必要があります：

```bash
# Docker Composeでデータベースを起動
docker-compose up -d postgres

# バックエンドサーバーを起動（別ターミナル）
cd backend && go run main.go

# フロントエンドサーバーを起動（別ターミナル）
cd frontend && npm run dev

# その後、E2Eテストを実行
make ci-e2e
```

### npm ciが失敗する

`package-lock.json`が`package.json`と同期していない可能性があります：

```bash
# フロントエンド
cd frontend && rm -rf node_modules package-lock.json && npm install

# E2E
cd e2e && rm -rf node_modules package-lock.json && npm install
```

## 推奨ワークフロー

1. **開発中**: `make ci-quick`でクイックチェック
2. **コミット前**: `make ci`で全チェック
3. **プルリクエスト前**: `make ci-all`でE2E以外の全チェック
4. **リリース前**: `make ci-all && make ci-e2e`で全チェック

