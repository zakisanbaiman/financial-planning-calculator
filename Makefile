.PHONY: help install setup lint format test clean dev build

# デフォルトターゲット
help:
	@echo "Financial Planning Calculator - Make Commands"
	@echo ""
	@echo "Setup:"
	@echo "  make install    - Install all dependencies"
	@echo "  make setup      - Setup Git hooks and tools"
	@echo ""
	@echo "Development:"
	@echo "  make dev        - Start development servers"
	@echo "  make lint       - Run linters"
	@echo "  make format     - Format code"
	@echo "  make test       - Run all tests"
	@echo ""
	@echo "Build:"
	@echo "  make build      - Build all projects"
	@echo "  make clean      - Clean build artifacts"

# 依存関係のインストール
install:
	@echo "Installing root dependencies..."
	npm install
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "Installing e2e dependencies..."
	cd e2e && npm install
	@echo "Installing backend dependencies..."
	cd backend && go mod download

# Git hooksのセットアップ
setup:
	@echo "Setting up Git hooks..."
	npm run prepare
	@echo "Git hooks installed!"

# Lintの実行
lint:
	@echo "Running linters..."
	npm run lint

# コードフォーマット
format:
	@echo "Formatting code..."
	npm run format
	@echo "Formatting YAML files..."
	npx prettier --write "**/*.{yml,yaml,json,md}"

# テストの実行
test:
	@echo "Running tests..."
	npm run test

# E2Eテストの実行
test-e2e:
	@echo "Running E2E tests..."
	npm run test:e2e

# 統合テストの実行
test-integration:
	@echo "Running integration tests..."
	./scripts/test-integration.sh

# 開発サーバーの起動
dev:
	@echo "Starting development servers..."
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@echo ""
	@echo "Press Ctrl+C to stop"
	@make -j2 dev-backend dev-frontend

dev-backend:
	cd backend && go run main.go

dev-frontend:
	cd frontend && npm run dev

# ビルド
build:
	@echo "Building projects..."
	npm run build:backend
	npm run build:frontend

# クリーンアップ
clean:
	@echo "Cleaning build artifacts..."
	rm -rf frontend/.next
	rm -rf frontend/out
	rm -rf backend/server
	rm -rf e2e/test-results
	rm -rf e2e/playwright-report
	@echo "Clean complete!"

# バックエンドのみ起動
backend:
	cd backend && go run main.go

# フロントエンドのみ起動
frontend:
	cd frontend && npm run dev

# データベースのセットアップ（将来用）
db-setup:
	@echo "Setting up database..."
	# TODO: Add database setup commands

# 依存関係の更新
update:
	@echo "Updating dependencies..."
	cd frontend && npm update
	cd e2e && npm update
	cd backend && go get -u ./...
	cd backend && go mod tidy

# セキュリティチェック
security:
	@echo "Running security checks..."
	cd frontend && npm audit
	cd e2e && npm audit
	cd backend && go list -json -m all | nancy sleuth

# すべてのチェックを実行（CI相当）
ci: lint test
	@echo "All CI checks passed!"
