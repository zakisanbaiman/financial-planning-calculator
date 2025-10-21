#!/bin/bash

# run-integration-tests.sh
# Script to run database integration tests for the financial planning calculator

set -e

echo "🚀 Starting Database Integration Tests"
echo "======================================"

# Check if PostgreSQL is running
if ! command -v psql &> /dev/null; then
    echo "❌ PostgreSQL client not found. Please install PostgreSQL."
    exit 1
fi

# Set test database environment variables
export DB_HOST=${DB_HOST:-localhost}
export DB_PORT=${DB_PORT:-5432}
export DB_USER=${DB_USER:-postgres}
export DB_PASSWORD=${DB_PASSWORD:-password}
export DB_NAME=${DB_NAME:-financial_planning_test}
export DB_SSLMODE=${DB_SSLMODE:-disable}

echo "📊 Database Configuration:"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  User: $DB_USER"
echo "  Database: $DB_NAME"
echo ""

# Check database connection
echo "🔍 Checking database connection..."
if ! PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "SELECT 1;" &> /dev/null; then
    echo "❌ Cannot connect to PostgreSQL. Please ensure PostgreSQL is running and credentials are correct."
    echo "   You can start PostgreSQL with Docker: docker-compose up -d postgres"
    exit 1
fi

echo "✅ Database connection successful"

# Create test database if it doesn't exist
echo "🗄️  Setting up test database..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;" 2>/dev/null || true
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME;" 2>/dev/null || true

echo "✅ Test database ready"

# Change to backend directory
cd "$(dirname "$0")/.."

echo ""
echo "🧪 Running Integration Tests"
echo "============================"

# Run repository integration tests
echo "📦 Testing Repository Layer..."
go test -v ./infrastructure/repositories -run "Test.*Repository.*" -timeout 30s

echo ""
echo "🔄 Testing Cross-Repository Operations..."
go test -v ./infrastructure/repositories -run "TestCrossRepository.*|TestTransactionConsistency|TestRepositoryFactory" -timeout 30s

echo ""
echo "🚀 Testing Concurrent Access..."
go test -v ./infrastructure/repositories -run "TestConcurrent.*" -timeout 60s

echo ""
echo "📈 Running Performance Tests..."
go test -v ./infrastructure/repositories -run "TestDatabaseStressTest|TestDatabaseConnectionLeaks|TestDatabaseDeadlockDetection" -timeout 120s

echo ""
echo "⚡ Running Benchmarks..."
go test -v ./infrastructure/repositories -bench "Benchmark.*" -benchtime=5s -timeout 60s

echo ""
echo "🧹 Cleaning up test database..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;" 2>/dev/null || true

echo ""
echo "✅ All Database Integration Tests Completed Successfully!"
echo "======================================================="

# Generate test coverage report if requested
if [ "$1" = "--coverage" ]; then
    echo ""
    echo "📊 Generating Coverage Report..."
    go test ./infrastructure/repositories -coverprofile=coverage.out -covermode=atomic
    go tool cover -html=coverage.out -o coverage.html
    echo "✅ Coverage report generated: coverage.html"
fi

echo ""
echo "📋 Test Summary:"
echo "  ✅ Repository CRUD operations"
echo "  ✅ Data integrity constraints"
echo "  ✅ Transaction consistency"
echo "  ✅ Concurrent access handling"
echo "  ✅ Performance under load"
echo "  ✅ Connection pool management"
echo "  ✅ Deadlock detection"
echo ""
echo "🎉 Integration testing complete!"