.PHONY: help install setup lint format test clean dev build docker-help dev-setup up up-full down logs logs-api logs-db migrate migrate-status migrate-down seed reset shell-api shell-db test-docker

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
	@echo "Docker Development:"
	@echo "  make docker-help     - Show Docker-specific commands"
	@echo "  make dev-setup       - First-time Docker setup (build, migrate, seed)"
	@echo "  make up              - Start Docker development environment"
	@echo "  make down            - Stop Docker environment"
	@echo "  make logs            - View all Docker logs"
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

# =============================================================================
# Docker Development Commands
# =============================================================================

docker-help:
	@echo "Docker Development Environment - Commands"
	@echo "=========================================="
	@echo ""
	@echo "Setup & Start:"
	@echo "  make dev-setup       - First-time setup (build, start, migrate, seed)"
	@echo "  make up              - Start backend + database (hot reload enabled)"
	@echo "  make up-full         - Start all services including frontend"
	@echo "  make down            - Stop all containers"
	@echo "  make restart         - Restart all containers"
	@echo ""
	@echo "Database Operations:"
	@echo "  make migrate         - Run database migrations"
	@echo "  make migrate-status  - Check migration status"
	@echo "  make migrate-down    - Rollback last migration"
	@echo "  make seed            - Seed database with sample data"
	@echo "  make reset           - Reset database (down + migrate + seed)"
	@echo ""
	@echo "Development:"
	@echo "  make logs            - View all logs"
	@echo "  make logs-api        - View backend API logs"
	@echo "  make logs-db         - View database logs"
	@echo "  make shell-api       - Access backend container shell"
	@echo "  make shell-db        - Access PostgreSQL shell"
	@echo "  make test-docker     - Run tests in Docker"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean-docker    - Remove containers and volumes"
	@echo "  make rebuild         - Rebuild Docker images"

# First-time setup
dev-setup:
	@echo "🚀 Setting up Docker development environment..."
	docker compose build
	docker compose up -d postgres
	@echo "⏳ Waiting for database to be ready..."
	@sleep 5
	@$(MAKE) migrate
	@$(MAKE) seed
	docker compose up -d backend
	@echo "✅ Setup complete!"
	@echo ""
	@echo "Backend API: http://localhost:8080"
	@echo "Swagger UI:  http://localhost:8080/swagger/index.html"
	@echo "Database:    localhost:5432"
	@echo ""
	@echo "Use 'make logs' to view logs"
	@echo "Use 'make down' to stop"

# Start development environment
up:
	@echo "🚀 Starting Docker development environment..."
	docker compose up -d postgres backend
	@echo "✅ Started! Backend with hot reload at http://localhost:8080"
	@echo "Use 'make logs' to view logs"

# Start all services including frontend
up-full:
	@echo "🚀 Starting all services..."
	docker compose --profile frontend up -d
	@echo "✅ All services started!"
	@echo "Backend:  http://localhost:8080"
	@echo "Frontend: http://localhost:3000"

# Stop all services
down:
	@echo "🛑 Stopping Docker environment..."
	docker compose down
	@echo "✅ Stopped!"

# Restart services
restart:
	@echo "🔄 Restarting services..."
	docker compose restart
	@echo "✅ Restarted!"

# View all logs
logs:
	docker compose logs -f

# View backend logs
logs-api:
	docker compose logs -f backend

# View database logs
logs-db:
	docker compose logs -f postgres

# Run migrations
migrate:
	@echo "📦 Running database migrations..."
	docker compose run --rm db-tools go run ./cmd/migrate/main.go -command=up
	@echo "✅ Migrations complete!"

# Check migration status
migrate-status:
	@echo "📊 Checking migration status..."
	docker compose run --rm db-tools go run ./cmd/migrate/main.go -command=status

# Rollback migration
migrate-down:
	@echo "⏪ Rolling back last migration..."
	docker compose run --rm db-tools go run ./cmd/migrate/main.go -command=down
	@echo "✅ Rollback complete!"

# Seed database
seed:
	@echo "🌱 Seeding database..."
	docker compose run --rm db-tools go run ./cmd/seed/main.go
	@echo "✅ Seeding complete!"

# Reset database
reset:
	@echo "🔄 Resetting database..."
	@$(MAKE) migrate-down
	@$(MAKE) migrate
	@$(MAKE) seed
	@echo "✅ Database reset complete!"

# Access backend container shell
shell-api:
	@echo "🐚 Accessing backend container..."
	docker compose exec backend sh

# Access database shell
shell-db:
	@echo "🐚 Accessing PostgreSQL..."
	docker compose exec postgres psql -U postgres -d financial_planning

# Run tests in Docker
test-docker:
	@echo "🧪 Running tests in Docker..."
	docker compose run --rm backend go test -v ./...

# Clean up Docker resources
clean-docker:
	@echo "🧹 Cleaning up Docker resources..."
	docker compose down -v
	@echo "✅ Cleanup complete!"

# Rebuild Docker images
rebuild:
	@echo "🔨 Rebuilding Docker images..."
	docker compose build --no-cache
	@echo "✅ Rebuild complete!"
