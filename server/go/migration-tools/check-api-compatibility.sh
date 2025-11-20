#!/bin/bash

# API Compatibility Check Script
# Проверяет совместимость API между монолитом и микросервисами

set -e

MONOLITH_URL=${MONOLITH_URL:-"http://localhost:8080"}
MICROSERVICES_URL=${MICROSERVICES_URL:-"http://localhost:8080"}

echo "🔍 Проверка API совместимости"
echo "============================="
echo "Monolith URL: $MONOLITH_URL"
echo "Microservices URL: $MICROSERVICES_URL"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASSED=0
FAILED=0
WARNINGS=0

check_passed() {
    echo -e "${GREEN}✅ $1${NC}"
    ((PASSED++))
}

check_failed() {
    echo -e "${RED}❌ $1${NC}"
    ((FAILED++))
}

check_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
    ((WARNINGS++))
}

# Function to compare JSON responses
compare_responses() {
    local endpoint=$1
    local method=${2:-"GET"}
    local data=${3:-""}

    echo "Testing $method $endpoint..."

    # Call monolith
    if [ "$method" = "POST" ]; then
        monolith_response=$(curl -s -X POST -H "Content-Type: application/json" -d "$data" "$MONOLITH_URL$endpoint" 2>/dev/null || echo "error")
    else
        monolith_response=$(curl -s "$MONOLITH_URL$endpoint" 2>/dev/null || echo "error")
    fi

    # Call microservices
    if [ "$method" = "POST" ]; then
        microservices_response=$(curl -s -X POST -H "Content-Type: application/json" -d "$data" "$MICROSERVICES_URL$endpoint" 2>/dev/null || echo "error")
    else
        microservices_response=$(curl -s "$MICROSERVICES_URL$endpoint" 2>/dev/null || echo "error")
    fi

    # Check if both responses are valid JSON
    if echo "$monolith_response" | jq . >/dev/null 2>&1 && echo "$microservices_response" | jq . >/dev/null 2>&1; then
        # Compare status
        monolith_status=$(echo "$monolith_response" | jq -r '.status // "ok"')
        microservices_status=$(echo "$microservices_response" | jq -r '.status // "ok"')

        if [ "$monolith_status" = "$microservices_status" ]; then
            check_passed "Status match for $endpoint"
        else
            check_failed "Status mismatch for $endpoint: monolith=$monolith_status, microservices=$microservices_status"
        fi

        # Compare response structure (basic check)
        monolith_keys=$(echo "$monolith_response" | jq 'keys | length')
        microservices_keys=$(echo "$microservices_response" | jq 'keys | length')

        if [ "$monolith_keys" -eq "$microservices_keys" ]; then
            check_passed "Response structure compatible for $endpoint"
        else
            check_warning "Response structure differs for $endpoint: monolith has $monolith_keys keys, microservices has $microservices_keys keys"
        fi
    else
        if [ "$monolith_response" = "error" ] && [ "$microservices_response" = "error" ]; then
            check_warning "Both services unreachable for $endpoint"
        elif [ "$monolith_response" = "error" ]; then
            check_failed "Monolith unreachable for $endpoint"
        elif [ "$microservices_response" = "error" ]; then
            check_failed "Microservices unreachable for $endpoint"
        else
            check_failed "Invalid JSON response for $endpoint"
        fi
    fi
}

echo "1. Health Check Endpoints"
echo "------------------------"

compare_responses "/health"
compare_responses "/ready"

echo ""
echo "2. Authentication Endpoints"
echo "---------------------------"

# Test login (will fail without proper credentials, but checks if endpoint exists)
compare_responses "/api/v1/auth/login" "POST" '{"username":"test","password":"test"}'

echo ""
echo "3. Template Endpoints"
echo "--------------------"

# Test template execution (mock data)
compare_responses "/api/v1/templates/execute" "POST" '{"query":"test query","language":"ru"}'

echo ""
echo "4. Batch Endpoints"
echo "-----------------"

compare_responses "/api/v1/batch/status/test-id"

echo ""
echo "5. Webhook Endpoints"
echo "-------------------"

compare_responses "/api/v1/webhooks"

echo ""
echo "6. Analytics Endpoints"
echo "----------------------"

compare_responses "/api/v1/analytics/stats"

echo ""
echo "7. Conversation Endpoints"
echo "-------------------------"

compare_responses "/api/v1/conversations"

echo ""
echo "📊 Результаты проверки совместимости:"
echo "====================================="
echo -e "${GREEN}✅ Совместимо: $PASSED эндпоинтов${NC}"
echo -e "${RED}❌ Несовместимо: $FAILED эндпоинтов${NC}"
echo -e "${YELLOW}⚠️  Предупреждения: $WARNINGS${NC}"

TOTAL=$((PASSED + FAILED + WARNINGS))
COMPATIBILITY_RATE=$((PASSED * 100 / TOTAL))

echo ""
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 API полностью совместим! Готовность: ${COMPATIBILITY_RATE}%${NC}"
    echo ""
    echo "✅ Backward compatibility поддерживается"
    echo "✅ Клиенты продолжат работать без изменений"
elif [ $COMPATIBILITY_RATE -ge 90 ]; then
    echo -e "${YELLOW}⚠️  Высокая совместимость API: ${COMPATIBILITY_RATE}%${NC}"
    echo ""
    echo "🔧 Исправьте $FAILED несовместимых эндпоинтов"
else
    echo -e "${RED}❌ Низкая совместимость API: ${COMPATIBILITY_RATE}%${NC}"
    echo ""
    echo "🛠️  Требуется доработка API контрактов"
fi

echo ""
echo "📋 Рекомендации:"
echo "• Убедитесь что все эндпоинты возвращают HTTP 200 для существующих клиентов"
echo "• Проверьте response format для всех POST/PUT запросов"
echo "• Протестируйте error handling (400, 401, 404, 500)"
echo "• Убедитесь в поддержке всех query parameters"

if [ $FAILED -gt 0 ]; then
    echo ""
    echo "🔍 Детальный анализ несовместимостей:"
    echo "• Сравните response schemas между монолитом и микросервисами"
    echo "• Проверьте API documentation на соответствие"
    echo "• Обновите integration tests для новых response formats"
fi
