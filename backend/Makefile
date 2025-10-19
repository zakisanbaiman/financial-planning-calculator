# Makefile for Financial Planning Calculator Backend (Local Development)

.PHONY: help build run test clean migrate-up migrate-down migrate-status seed

# Default target
help:
	@echo "財務計画計算機バックエンド - ローカル開発用"
	@echo "=========================================="
	@echo "利用可能なコマンド:"
	@echo "  build         - アプリケーションをビルド"
	@echo "  run           - アプリケーションを実行"
	@echo "  test          - テストを実行"
	@echo "  test-coverage - カバレッジ付きでテストを実行"
	@echo "  clean         - ビルド成果物を削除"
	@echo "  migrate-up    - データベースマイグレーションを実行"
	@echo "  migrate-down  - 最新のマイグレーションをロールバック"
	@echo "  migrate-status- マイグレーション状況を確認"
	@echo "  seed          - サンプルデータを投入"
	@echo "  db-reset      - データベースをリセット（全削除→マイグレーション→シード）"
	@echo ""
	@echo "Docker開発環境を使用する場合は、プロジェクトルートの Makefile を使用してください"

# Build the application
build:
	@echo "アプリケーションをビルド中..."
	go build -o bin/server ./main.go
	go build -o bin/migrate ./cmd/migrate/main.go
	go build -o bin/seed ./cmd/seed/main.go

# Run the application
run:
	@echo "アプリケーションを実行中..."
	go run main.go

# Run tests
test:
	@echo "テストを実行中..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "カバレッジ付きでテストを実行中..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "カバレッジレポートが coverage.html に生成されました"

# Clean build artifacts
clean:
	@echo "ビルド成果物を削除中..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Database migration up
migrate-up:
	@echo "データベースマイグレーションを実行中..."
	go run ./cmd/migrate/main.go -command=up

# Database migration down
migrate-down:
	@echo "マイグレーションをロールバック中..."
	go run ./cmd/migrate/main.go -command=down

# Check migration status
migrate-status:
	@echo "マイグレーション状況を確認中..."
	go run ./cmd/migrate/main.go -command=status

# Seed database with sample data
seed:
	@echo "サンプルデータを投入中..."
	go run ./cmd/seed/main.go

# Reset database (drop all, migrate, seed)
db-reset: migrate-down migrate-up seed
	@echo "データベースのリセットが完了しました"

# Install dependencies
deps:
	@echo "依存関係をインストール中..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "コードをフォーマット中..."
	go fmt ./...

# Lint code
lint:
	@echo "コードをリント中..."
	golangci-lint run

# Generate documentation
docs:
	@echo "ドキュメントを生成中..."
	swag init -g main.go -o ./docs

# Development setup
dev-setup: deps migrate-up seed
	@echo "開発環境のセットアップが完了しました"

# Database setup (interactive)
db-setup:
	@echo "データベースセットアップを開始します..."
	./scripts/setup-database.sh