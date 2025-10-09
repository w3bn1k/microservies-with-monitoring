#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Тестирование API микросервисов${NC}"
echo "=================================="

# Test Producer
echo -e "\n${YELLOW}📤 Тестирование Producer Service${NC}"
echo "Health check:"
curl -s http://localhost:8080/health | jq '.' 2>/dev/null || curl -s http://localhost:8080/health

echo -e "\nСтатистика Producer:"
curl -s http://localhost:8080/api/v1/stats | jq '.' 2>/dev/null || curl -s http://localhost:8080/api/v1/stats

echo -e "\nГенерация события:"
curl -s -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{"type":"user_action","user_id":"test_user","data":{"action":"test"}}' | jq '.' 2>/dev/null || curl -s -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{"type":"user_action","user_id":"test_user","data":{"action":"test"}}'

# Test Consumer
echo -e "\n${YELLOW}📥 Тестирование Consumer Service${NC}"
echo "Health check:"
curl -s http://localhost:8081/health | jq '.' 2>/dev/null || curl -s http://localhost:8081/health

echo -e "\nСтатистика Consumer:"
curl -s http://localhost:8081/api/v1/stats | jq '.' 2>/dev/null || curl -s http://localhost:8081/api/v1/stats

# Test Monitor
echo -e "\n${YELLOW}📊 Тестирование Monitor Service${NC}"
echo "Health check:"
curl -s http://localhost:8082/health | jq '.' 2>/dev/null || curl -s http://localhost:8082/health

echo -e "\nСтатистика Monitor:"
curl -s http://localhost:8082/api/v1/stats | jq '.' 2>/dev/null || curl -s http://localhost:8082/api/v1/stats

echo -e "\nТранзакции:"
curl -s http://localhost:8082/api/v1/transactions | jq '.' 2>/dev/null || curl -s http://localhost:8082/api/v1/transactions

# Test Metrics
echo -e "\n${YELLOW}📈 Тестирование метрик${NC}"
echo "Producer metrics (первые 10 строк):"
curl -s http://localhost:8080/metrics | head -10

echo -e "\n${GREEN}✅ Тестирование завершено!${NC}"
echo -e "\n${BLUE}Доступные URL:${NC}"
echo "  🌐 Producer API:     http://localhost:8080"
echo "  🌐 Consumer API:     http://localhost:8081"
echo "  🌐 Monitor API:      http://localhost:8082"
echo "  📈 Kafka UI:         http://localhost:8083"
echo "  🔴 Redis Commander:  http://localhost:8084"
echo "  📊 Prometheus:       http://localhost:9090"
echo "  📈 Grafana:          http://localhost:3000"
echo "  🔍 Jaeger:           http://localhost:16686"
echo "  🗄️  pgAdmin:          http://localhost:8085"
