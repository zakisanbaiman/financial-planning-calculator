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
	@echo "CI (Local):"
	@echo "  make ci         - Run all CI checks (lint + test + pr-check)"
	@echo "  make ci-lint    - Run lint workflow (backend + frontend)"
	@echo "  make ci-test    - Run test workflow (backend + frontend)"
	@echo "  make ci-pr-check - Run PR check workflow (quick tests)"
	@echo "  make ci-e2e     - Run E2E tests (requires DB and servers)"
	@echo "  make ci-all     - Run all CI workflows (except E2E)"
	@echo "  make ci-quick   - Run quick CI checks (lint + pr-check)"
	@echo "  ./scripts/run-ci-local.sh [workflow] - Run specific workflow"
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

# CIワークフローをローカルで実行
ci: ci-lint ci-test ci-pr-check
	@echo "✅ All CI checks passed!"

# Lintワークフロー（.github/workflows/lint.yml相当）
ci-lint: ci-lint-backend ci-lint-frontend
	@echo "✅ Lint checks passed!"

ci-lint-backend:
	@echo "🔍 Running Go lint checks..."
	@cd backend && \
		go mod download && \
		go mod tidy && \
		go mod verify && \
		(which golangci-lint > /dev/null && golangci-lint run --timeout=5m --verbose || echo "⚠️  golangci-lint not installed, skipping...") && \
		go fmt ./... && \
		go vet ./...

ci-lint-frontend:
	@echo "🔍 Running Frontend lint checks..."
	@cd frontend && \
		(npm ci || npm install) && \
		([ -f .eslintrc.json ] || echo '{"extends": ["next/core-web-vitals"]}' > .eslintrc.json) && \
		npm run type-check && \
		npm run lint -- --max-warnings 0

# Testワークフロー（.github/workflows/test.yml相当）
ci-test: ci-test-backend ci-test-frontend
	@echo "✅ Test checks passed!"

ci-test-backend:
	@echo "🧪 Running Backend tests..."
	@cd backend && \
		go mod download && \
		go mod tidy && \
		go mod verify && \
		go build -v ./... && \
		go test -v -race -timeout 30s ./... && \
		go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

ci-test-frontend:
	@echo "🧪 Running Frontend build..."
	@cd frontend && \
		(npm ci || npm install) && \
		npm run build

# PR Checkワークフロー（.github/workflows/pr-check.yml相当）
ci-pr-check:
	@echo "🔍 Running PR check (quick tests)..."
	@cd backend && \
		go mod download && \
		go mod tidy && \
		go vet ./... && \
		go test -v -short ./...

# E2Eテストワークフロー（.github/workflows/e2e-tests.yml相当）
# 注意: データベースとサーバーが必要です
ci-e2e:
	@echo "🧪 Running E2E tests..."
	@echo "⚠️  Make sure database and servers are running!"
	@cd e2e && \
		(npm ci || npm install) && \
		npx playwright install --with-deps && \
		npm test

# すべてのCIワークフローを実行（E2E除く）
ci-all: ci-lint ci-test ci-pr-check
	@echo "✅ All CI checks (except E2E) passed!"

# クイックチェック（lint + クイックテスト）
ci-quick: ci-lint-backend ci-pr-check
	@echo "✅ Quick CI checks passed!"
