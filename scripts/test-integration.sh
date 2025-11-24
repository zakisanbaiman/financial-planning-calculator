#!/bin/bash

# Integration Test Script
# Tests frontend-backend connectivity and API endpoints

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BACKEND_URL="${BACKEND_URL:-http://localhost:8080}"
FRONTEND_URL="${FRONTEND_URL:-http://localhost:3000}"
API_URL="${BACKEND_URL}/api"

echo "========================================="
echo "財務計画計算機 統合テスト"
echo "========================================="
echo ""
echo "Backend URL: ${BACKEND_URL}"
echo "Frontend URL: ${FRONTEND_URL}"
echo "API URL: ${API_URL}"
echo ""

# Function to test endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local description=$3
    local data=$4
    
    echo -n "Testing ${description}... "
    
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X ${method} "${endpoint}" \
            -H "Content-Type: application/json" 2>&1)
    else
        response=$(curl -s -w "\n%{http_code}" -X ${method} "${endpoint}" \
            -H "Content-Type: application/json" \
            -d "${data}" 2>&1)
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP ${http_code})"
        return 0
    elif [ "$http_code" -ge 400 ] && [ "$http_code" -lt 500 ]; then
        echo -e "${YELLOW}⚠ WARN${NC} (HTTP ${http_code})"
        return 1
    else
        echo -e "${RED}✗ FAIL${NC} (HTTP ${http_code})"
        return 1
    fi
}

# Test counter
total_tests=0
passed_tests=0
failed_tests=0

# 1. Backend Health Check
echo "========================================="
echo "1. バックエンドヘルスチェック"
echo "========================================="
total_tests=$((total_tests + 1))
if test_endpoint "GET" "${BACKEND_URL}/health" "Basic health check"; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

total_tests=$((total_tests + 1))
if test_endpoint "GET" "${BACKEND_URL}/health/detailed" "Detailed health check"; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

total_tests=$((total_tests + 1))
if test_endpoint "GET" "${BACKEND_URL}/ready" "Readiness check"; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

echo ""

# 2. API Info
echo "========================================="
echo "2. API情報"
echo "========================================="
total_tests=$((total_tests + 1))
if test_endpoint "GET" "${API_URL}/" "API info endpoint"; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

echo ""

# 3. CORS Check
echo "========================================="
echo "3. CORS設定確認"
echo "========================================="
echo -n "Testing CORS preflight... "
response=$(curl -s -w "\n%{http_code}" -X OPTIONS "${API_URL}/financial-data" \
    -H "Origin: http://localhost:3000" \
    -H "Access-Control-Request-Method: POST" \
    -H "Access-Control-Request-Headers: Content-Type" 2>&1)

http_code=$(echo "$response" | tail -n1)
total_tests=$((total_tests + 1))

if [ "$http_code" -eq 204 ] || [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✓ PASS${NC} (HTTP ${http_code})"
    passed_tests=$((passed_tests + 1))
else
    echo -e "${RED}✗ FAIL${NC} (HTTP ${http_code})"
    failed_tests=$((failed_tests + 1))
fi

echo ""

# 4. Calculation Endpoints
echo "========================================="
echo "4. 計算エンドポイント"
echo "========================================="

# Asset Projection
asset_data='{
  "user_id": "test-user-001",
  "monthly_income": 400000,
  "monthly_expenses": 280000,
  "current_savings": 1500000,
  "investment_return": 5.0,
  "inflation_rate": 2.0,
  "projection_years": 30
}'

total_tests=$((total_tests + 1))
if test_endpoint "POST" "${API_URL}/calculations/asset-projection" "Asset projection calculation" "$asset_data"; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

# Retirement Calculation
retirement_data='{
  "user_id": "test-user-001",
  "current_age": 35,
  "retirement_age": 65,
  "life_expectancy": 90,
  "monthly_retirement_expenses": 250000,
  "pension_amount": 150000,
  "current_savings": 1500000,
  "monthly_savings": 120000,
  "investment_return": 5.0,
  "inflation_rate": 2.0
}'

total_tests=$((total_tests + 1))
if test_endpoint "POST" "${API_URL}/calculations/retirement" "Retirement calculation" "$retirement_data"; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

# Emergency Fund
emergency_data='{
  "user_id": "test-user-001",
  "monthly_expenses": 280000,
  "emergency_months": 6,
  "current_savings": 1500000
}'

total_tests=$((total_tests + 1))
if test_endpoint "POST" "${API_URL}/calculations/emergency-fund" "Emergency fund calculation" "$emergency_data"; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

echo ""

# 5. Frontend Accessibility
echo "========================================="
echo "5. フロントエンドアクセス確認"
echo "========================================="
echo -n "Testing frontend homepage... "
response=$(curl -s -w "\n%{http_code}" "${FRONTEND_URL}" 2>&1)
http_code=$(echo "$response" | tail -n1)
total_tests=$((total_tests + 1))

if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
    echo -e "${GREEN}✓ PASS${NC} (HTTP ${http_code})"
    passed_tests=$((passed_tests + 1))
else
    echo -e "${YELLOW}⚠ WARN${NC} (HTTP ${http_code}) - Frontend may not be running"
    failed_tests=$((failed_tests + 1))
fi

echo ""

# Summary
echo "========================================="
echo "テスト結果サマリー"
echo "========================================="
echo "Total Tests: ${total_tests}"
echo -e "Passed: ${GREEN}${passed_tests}${NC}"
echo -e "Failed: ${RED}${failed_tests}${NC}"
echo ""

if [ $failed_tests -eq 0 ]; then
    echo -e "${GREEN}✓ すべてのテストが成功しました！${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠ 一部のテストが失敗しました${NC}"
    exit 1
fi
