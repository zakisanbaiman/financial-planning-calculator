#!/bin/bash
# Post-create script for dev container
# このスクリプトは、dev containerが作成された後に一度だけ実行されます

set -e

echo "🚀 Setting up dev container..."

# ワークスペースディレクトリに移動
cd /workspace

# Git設定の確認
if [ ! -f ~/.gitconfig ]; then
    echo "📝 Configuring Git..."
    git config --global core.autocrlf false
    git config --global core.eol lf
    git config --global pull.rebase false
fi

# Bashヒストリーファイルの設定
echo "📝 Setting up bash history..."
HISTORY_DIR="/commandhistory"
mkdir -p "$HISTORY_DIR"
touch "$HISTORY_DIR/.bash_history"
if ! grep -q "HISTFILE=" ~/.bashrc; then
    echo "export HISTFILE=$HISTORY_DIR/.bash_history" >> ~/.bashrc
    echo "export PROMPT_COMMAND='history -a'" >> ~/.bashrc
fi

# Zshヒストリーファイルの設定
if [ -f ~/.zshrc ]; then
    mkdir -p "$HISTORY_DIR"
    touch "$HISTORY_DIR/.zsh_history"
    if ! grep -q "HISTFILE=" ~/.zshrc; then
        echo "export HISTFILE=$HISTORY_DIR/.zsh_history" >> ~/.zshrc
    fi
fi

# Go toolsのインストール
echo "🔧 Installing Go tools..."
go install github.com/air-verse/air@v1.52.3
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.0
go install github.com/swaggo/swag/cmd/swag@latest

# バックエンドの依存関係をダウンロード
echo "📦 Downloading Go dependencies..."
cd /workspace/backend
go mod download
go mod verify

# フロントエンドの依存関係をインストール
echo "📦 Installing frontend dependencies..."
cd /workspace/frontend
npm ci

# E2Eテストの依存関係をインストール
echo "📦 Installing e2e dependencies..."
cd /workspace/e2e
npm ci

# ルートの依存関係をインストール
echo "📦 Installing root dependencies..."
cd /workspace
npm ci

# Git hooksのセットアップ
echo "🪝 Setting up Git hooks..."
npm run prepare

# データベースのセットアップ（PostgreSQLが起動している場合）
echo "🗄️  Checking database..."
if pg_isready -h postgres -U postgres > /dev/null 2>&1; then
    echo "✅ Database is ready"
    
    # マイグレーションの実行
    echo "📦 Running database migrations..."
    cd /workspace/backend
    if go run ./cmd/migrate/main.go -command=up 2>/dev/null; then
        echo "✅ Migrations completed"
    else
        echo "⚠️  Migrations skipped (may already be up to date)"
    fi
    
    # シードデータの投入
    echo "🌱 Seeding database..."
    if go run ./cmd/seed/main.go 2>/dev/null; then
        echo "✅ Seeding completed"
    else
        echo "⚠️  Seeding skipped (may already have data)"
    fi
else
    echo "⚠️  Database is not ready yet. You can run migrations later with: make migrate"
fi

# 完了メッセージ
echo ""
echo "✅ Dev Container setup complete!"
echo ""
echo "🎉 You can now start developing!"
echo ""
echo "📚 Useful commands:"
echo "  make up          - Start backend + database with hot reload"
echo "  make down        - Stop all services"
echo "  make logs        - View logs"
echo "  make test        - Run tests"
echo "  make lint        - Run linters"
echo "  make help        - Show all available commands"
echo ""
echo "🌐 Services:"
echo "  Backend API:  http://localhost:8080"
echo "  Swagger UI:   http://localhost:8080/swagger/index.html"
echo "  Frontend:     http://localhost:3000 (run 'make up-full')"
echo "  pprof:        http://localhost:6060/debug/pprof/"
echo ""
