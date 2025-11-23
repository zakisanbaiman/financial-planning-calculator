#!/bin/bash

# ローカルでCIを実行するスクリプト
# 使用方法: ./scripts/run-ci-local.sh [workflow-name]
# 例: ./scripts/run-ci-local.sh lint

set -e

WORKFLOW_NAME="${1:-all}"

echo "🚀 Running CI workflows locally..."
echo ""

case "$WORKFLOW_NAME" in
  lint)
    echo "📋 Running Lint workflow..."
    make ci-lint
    ;;
  test)
    echo "📋 Running Test workflow..."
    make ci-test
    ;;
  pr-check)
    echo "📋 Running PR Check workflow..."
    make ci-pr-check
    ;;
  e2e)
    echo "📋 Running E2E Tests workflow..."
    echo "⚠️  Note: This requires database and servers to be running"
    make ci-e2e
    ;;
  all)
    echo "📋 Running all CI workflows (except E2E)..."
    make ci-all
    ;;
  quick)
    echo "📋 Running quick CI checks..."
    make ci-quick
    ;;
  *)
    echo "❌ Unknown workflow: $WORKFLOW_NAME"
    echo ""
    echo "Available workflows:"
    echo "  lint      - Run lint checks (backend + frontend)"
    echo "  test      - Run tests (backend + frontend build)"
    echo "  pr-check  - Run quick PR checks"
    echo "  e2e       - Run E2E tests (requires DB and servers)"
    echo "  all       - Run all workflows (except E2E)"
    echo "  quick     - Run quick checks (lint + pr-check)"
    exit 1
    ;;
esac

echo ""
echo "✅ CI workflow completed!"

