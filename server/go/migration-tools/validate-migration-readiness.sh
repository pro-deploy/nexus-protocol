#!/bin/bash

# Migration Readiness Validation Script
# Проверяет готовность инфраструктуры к миграции

set -e

echo "🔍 Проверка готовности к миграции на микросервисы"
echo "=================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

PASSED=0
FAILED=0

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
}

echo ""
echo "1. Проверка Kubernetes кластера..."
echo "-----------------------------------"

# Check kubectl access
if kubectl cluster-info >/dev/null 2>&1; then
    check_passed "kubectl доступен"
else
    check_failed "kubectl недоступен"
fi

# Check namespace exists
if kubectl get namespace nexus-prod >/dev/null 2>&1; then
    check_passed "namespace nexus-prod существует"
else
    check_warning "namespace nexus-prod не найден (будет создан)"
fi

# Check nodes
NODE_COUNT=$(kubectl get nodes --no-headers | wc -l)
if [ "$NODE_COUNT" -ge 3 ]; then
    check_passed "Kubernetes кластер имеет $NODE_COUNT nodes (минимум 3)"
else
    check_failed "Kubernetes кластер имеет только $NODE_COUNT nodes (нужно минимум 3)"
fi

echo ""
echo "2. Проверка инфраструктуры..."
echo "------------------------------"

# Check PostgreSQL
if kubectl get pods -n nexus-prod -l app=postgres >/dev/null 2>&1; then
    POSTGRES_READY=$(kubectl get pods -n nexus-prod -l app=postgres -o jsonpath='{.items[0].status.conditions[?(@.type=="Ready")].status}')
    if [ "$POSTGRES_READY" = "True" ]; then
        check_passed "PostgreSQL кластер готов"
    else
        check_failed "PostgreSQL кластер не готов"
    fi
else
    check_warning "PostgreSQL кластер не найден"
fi

# Check Redis
if kubectl get pods -n nexus-prod -l app=redis >/dev/null 2>&1; then
    REDIS_READY=$(kubectl get pods -n nexus-prod -l app=redis -o jsonpath='{.items[0].status.conditions[?(@.type=="Ready")].status}')
    if [ "$REDIS_READY" = "True" ]; then
        check_passed "Redis кластер готов"
    else
        check_failed "Redis кластер не готов"
    fi
else
    check_warning "Redis кластер не найден"
fi

# Check Keycloak
if kubectl get pods -n nexus-prod -l app=keycloak >/dev/null 2>&1; then
    KEYCLOAK_READY=$(kubectl get pods -n nexus-prod -l app=keycloak -o jsonpath='{.items[0].status.conditions[?(@.type=="Ready")].status}')
    if [ "$KEYCLOAK_READY" = "True" ]; then
        check_passed "Keycloak готов"
    else
        check_failed "Keycloak не готов"
    fi
else
    check_warning "Keycloak не найден"
fi

echo ""
echo "3. Проверка секретов и конфигурации..."
echo "--------------------------------------"

# Check secrets
if kubectl get secret nexus-secrets -n nexus-prod >/dev/null 2>&1; then
    check_passed "nexus-secrets существует"
else
    check_failed "nexus-secrets не найден"
fi

# Check configmaps
if kubectl get configmap nexus-config -n nexus-prod >/dev/null 2>&1; then
    check_passed "nexus-config configmap существует"
else
    check_warning "nexus-config configmap не найден"
fi

echo ""
echo "4. Проверка образов Docker..."
echo "-----------------------------"

# Check if images exist in registry (mock check)
# In real scenario, check your container registry
check_warning "Проверка образов Docker - требуется ручная проверка"
check_warning "Убедитесь что все образы nexus-protocol/*:v1.1.0 загружены в registry"

echo ""
echo "5. Проверка CI/CD..."
echo "-------------------"

# Check if GitHub Actions workflows exist
if [ -d ".github/workflows" ]; then
    WORKFLOW_COUNT=$(find .github/workflows -name "*.yml" | wc -l)
    if [ "$WORKFLOW_COUNT" -ge 2 ]; then
        check_passed "Найдено $WORKFLOW_COUNT CI/CD pipeline'ов"
    else
        check_warning "Найдено только $WORKFLOW_COUNT CI/CD pipeline'ов (рекомендуется минимум 2)"
    fi
else
    check_failed "CI/CD pipelines не найдены в .github/workflows/"
fi

echo ""
echo "6. Проверка кода сервисов..."
echo "----------------------------"

# Check if service directories exist
SERVICES=("ai-service" "auth-service" "batch-service" "webhook-service" "analytics-service" "conversation-service" "api-gateway")
for service in "${SERVICES[@]}"; do
    if [ -d "services/$service" ]; then
        check_passed "Сервис $service готов"
    else
        check_warning "Сервис $service не найден (нужно создать директорию services/$service)"
    fi
done

echo ""
echo "7. Проверка мониторинга..."
echo "--------------------------"

# Check Prometheus
if kubectl get pods -n monitoring -l app=prometheus >/dev/null 2>&1; then
    check_passed "Prometheus найден"
else
    check_warning "Prometheus не найден в namespace monitoring"
fi

# Check Grafana
if kubectl get pods -n monitoring -l app=grafana >/dev/null 2>&1; then
    check_passed "Grafana найден"
else
    check_warning "Grafana не найден в namespace monitoring"
fi

echo ""
echo "8. Проверка сети и безопасности..."
echo "----------------------------------"

# Check network policies
NP_COUNT=$(kubectl get networkpolicies -n nexus-prod --no-headers 2>/dev/null | wc -l)
if [ "$NP_COUNT" -ge 1 ]; then
    check_passed "Найдено $NP_COUNT network policies"
else
    check_warning "Network policies не найдены (рекомендуется настроить)"
fi

# Check service mesh (Istio)
if kubectl get pods -n istio-system -l app=istiod >/dev/null 2>&1; then
    check_passed "Service mesh (Istio) найден"
else
    check_warning "Service mesh не найден (рекомендуется Istio)"
fi

echo ""
echo "📊 Результаты проверки:"
echo "======================"
echo -e "${GREEN}✅ Пройдено: $PASSED проверок${NC}"
echo -e "${RED}❌ Провалено: $FAILED проверок${NC}"

TOTAL=$((PASSED + FAILED))
SUCCESS_RATE=$((PASSED * 100 / TOTAL))

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 Все проверки пройдены! Готовность к миграции: 100%${NC}"
    echo ""
    echo "🚀 Можно приступать к Этапу 1: AI Service migration"
elif [ $SUCCESS_RATE -ge 80 ]; then
    echo -e "${YELLOW}⚠️  Высокая готовность к миграции: ${SUCCESS_RATE}%${NC}"
    echo ""
    echo "🔧 Исправьте критические проблемы перед началом миграции"
else
    echo -e "${RED}❌ Низкая готовность к миграции: ${SUCCESS_RATE}%${NC}"
    echo ""
    echo "🛠️  Требуется дополнительная подготовка инфраструктуры"
fi

echo ""
echo "📋 Следующие шаги:"
echo "1. Исправить все FAILED проверки"
echo "2. Настроить недостающие компоненты"
echo "3. Протестировать staging окружение"
echo "4. Создать план rollback'а"
echo ""
echo "📖 Подробная документация: MICROSERVICES_MIGRATION_PLAN.md"
