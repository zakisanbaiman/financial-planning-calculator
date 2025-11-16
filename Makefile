# Docker環境用Makefile

.PHONY: help build up down logs clean migrate seed reset test

# デフォルトターゲット
help:
	@echo "財務計画計算機 - Docker開発環境"
	@echo "================================"
	@echo "利用可能なコマンド:"
	@echo "  build     - Dockerイメージをビルド"
	@echo "  up        - 開発環境を起動（バックエンド + DB）"
	@echo "  up-full   - 全サービスを起動（フロントエンド含む）"
	@echo "  down      - 環境を停止"
	@echo "  logs      - ログを表示"
	@echo "  clean     - 全てのコンテナとボリュームを削除"
	@echo "  migrate   - データベースマイグレーションを実行"
	@echo "  seed      - サンプルデータを投入"
	@echo "  reset     - データベースをリセット"
	@echo "  test      - テストを実行"
	@echo "  lint      - Lintを実行"
	@echo "  lint-verbose - Lintを実行（詳細ログ）"
	@echo "  go-build  - Goアプリケーションをビルド"
	@echo "  go-fmt    - コードをフォーマット"
	@echo "  go-run    - Goコマンドを実行（例: make go-run CMD='go version'）"
	@echo "  shell-db  - データベースに接続"
	@echo "  shell-api - バックエンドコンテナに接続"

# Dockerイメージをビルド
build:
	@echo "Dockerイメージをビルド中..."
	docker-compose build

# 開発環境を起動（バックエンド + DB）
up:
	@echo "開発環境を起動中..."
	docker-compose up -d postgres backend
	@echo "起動完了！"
	@echo "API: http://localhost:8080"
	@echo "Swagger: http://localhost:8080/swagger/index.html"

# 全サービスを起動（フロントエンド含む）
up-full:
	@echo "全サービスを起動中..."
	docker-compose --profile frontend up -d
	@echo "起動完了！"
	@echo "フロントエンド: http://localhost:3000"
	@echo "API: http://localhost:8080"

# 環境を停止
down:
	@echo "環境を停止中..."
	docker-compose down

# ログを表示
logs:
	docker-compose logs -f

# 特定サービスのログを表示
logs-api:
	docker-compose logs -f backend

logs-db:
	docker-compose logs -f postgres

# 全てのコンテナとボリュームを削除
clean:
	@echo "全てのコンテナとボリュームを削除中..."
	docker-compose down -v --remove-orphans
	docker system prune -f

# データベースマイグレーションを実行
migrate:
	@echo "マイグレーションを実行中..."
	docker-compose run --rm db-tools go run ./cmd/migrate/main.go -command=up

# マイグレーション状況を確認
migrate-status:
	@echo "マイグレーション状況を確認中..."
	docker-compose run --rm db-tools go run ./cmd/migrate/main.go -command=status

# マイグレーションをロールバック
migrate-down:
	@echo "マイグレーションをロールバック中..."
	docker-compose run --rm db-tools go run ./cmd/migrate/main.go -command=down

# サンプルデータを投入
seed:
	@echo "サンプルデータを投入中..."
	docker-compose run --rm db-tools go run ./cmd/seed/main.go

# データベースをリセット（マイグレーション + シード）
reset: migrate seed
	@echo "データベースのリセットが完了しました"

# テストを実行
test:
	@echo "テストを実行中..."
	docker-compose run --rm backend go test ./... -v

# テスト（カバレッジ付き）
test-coverage:
	@echo "カバレッジ付きでテストを実行中..."
	docker-compose run --rm backend go test ./... -v -coverprofile=coverage.out
	docker-compose run --rm backend go tool cover -html=coverage.out -o coverage.html

# データベースに接続
shell-db:
	@echo "PostgreSQLに接続中..."
	docker-compose exec postgres psql -U postgres -d financial_planning

# バックエンドコンテナに接続
shell-api:
	@echo "バックエンドコンテナに接続中..."
	docker-compose exec backend sh

# 開発環境のセットアップ
dev-setup: build up wait-for-services migrate seed
	@echo "開発環境のセットアップが完了しました！"
	@echo ""
	@echo "🎉 セットアップ完了！"
	@echo "API: http://localhost:8080"
	@echo "Swagger: http://localhost:8080/swagger/index.html"
	@echo "データベース: localhost:5432"

# サービスの起動を待機
wait-for-services:
	@echo "サービスの起動を待機中..."
	@timeout=60; \
	while [ $$timeout -gt 0 ]; do \
		if docker-compose exec -T postgres pg_isready -U postgres -d financial_planning >/dev/null 2>&1; then \
			echo "✅ PostgreSQLが起動しました"; \
			break; \
		fi; \
		echo "PostgreSQLの起動を待機中... (残り$${timeout}秒)"; \
		sleep 2; \
		timeout=$$((timeout-2)); \
	done; \
	if [ $$timeout -le 0 ]; then \
		echo "❌ PostgreSQLの起動がタイムアウトしました"; \
		exit 1; \
	fi

# 本番用ビルド
build-prod:
	@echo "本番用イメージをビルド中..."
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml build

# 依存関係の更新
update-deps:
	@echo "依存関係を更新中..."
	docker-compose run --rm backend go mod tidy
	docker-compose run --rm backend go mod download

# Lintを実行
lint:
	@echo "Lintを実行中..."
	docker-compose run --rm backend golangci-lint run -v

# Lintを実行（詳細ログ付き）
lint-verbose:
	@echo "Lintを実行中（詳細ログ）..."
	docker-compose run --rm backend golangci-lint run -v --print-issued-lines --print-linter-name

# Goコマンドを実行（例: make go-run CMD="go version"）
go-run:
	docker-compose run --rm backend $(CMD)

# ビルドを実行
go-build:
	@echo "ビルドを実行中..."
	docker-compose run --rm backend go build -o bin/server ./main.go

# フォーマットを実行
go-fmt:
	@echo "コードをフォーマット中..."
	docker-compose run --rm backend go fmt ./...