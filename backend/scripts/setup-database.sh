#!/bin/bash

# Database setup script for Financial Planning Calculator

set -e

echo "財務計画計算機 - データベースセットアップスクリプト"
echo "=================================================="

# Check if PostgreSQL is running
if ! pg_isready -h localhost -p 5432 > /dev/null 2>&1; then
    echo "❌ PostgreSQLが起動していません。PostgreSQLを起動してから再実行してください。"
    echo "   macOS: brew services start postgresql"
    echo "   Linux: sudo systemctl start postgresql"
    exit 1
fi

echo "✅ PostgreSQLが起動しています"

# Check if database exists
DB_NAME=${DB_NAME:-financial_planning}
if ! psql -h localhost -p 5432 -U postgres -lqt | cut -d \| -f 1 | grep -qw $DB_NAME; then
    echo "📦 データベース '$DB_NAME' を作成中..."
    createdb -h localhost -p 5432 -U postgres $DB_NAME
    echo "✅ データベース '$DB_NAME' を作成しました"
else
    echo "✅ データベース '$DB_NAME' は既に存在します"
fi

# Run migrations
echo "🔄 マイグレーションを実行中..."
go run ./cmd/migrate/main.go -command=up

# Seed data (optional)
read -p "サンプルデータを投入しますか？ (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🌱 サンプルデータを投入中..."
    go run ./cmd/seed/main.go
    echo "✅ サンプルデータの投入が完了しました"
fi

echo ""
echo "🎉 データベースのセットアップが完了しました！"
echo ""
echo "次のコマンドでマイグレーション状況を確認できます:"
echo "  make migrate-status"
echo ""
echo "次のコマンドでアプリケーションを起動できます:"
echo "  make run"