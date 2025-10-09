#!/bin/bash
# Скрипт для проверки всех систем мониторинга

echo "🔍 Проверка статуса всех систем..."

# Проверка Prometheus
echo "📊 Prometheus:"
if curl -s http://localhost:9090/api/v1/query?query=up > /dev/null; then
  echo "  ✅ Prometheus работает"
  metrics=$(curl -s http://localhost:9090/api/v1/query?query=events_processed_total | jq -r '.data.result | length' 2>/dev/null || echo "0")
  echo "  📈 Метрики: $metrics"
else
  echo "  ❌ Prometheus недоступен"
fi

# Проверка Grafana
echo "📈 Grafana:"
if curl -s http://localhost:3000/api/health > /dev/null; then
  echo "  ✅ Grafana работает"
else
  echo "  ❌ Grafana недоступен"
fi

# Проверка Jaeger
echo "🔍 Jaeger:"
if curl -s http://localhost:16686/api/services > /dev/null; then
  echo "  ✅ Jaeger работает"
  services=$(curl -s http://localhost:16686/api/services | jq -r '.[] | .name' | wc -l 2>/dev/null || echo "0")
  echo "  📊 Сервисов в трейсинге: $services"
else
  echo "  ❌ Jaeger недоступен"
fi

# Проверка pgAdmin
echo "🗄️ pgAdmin:"
if curl -s http://localhost:8085 > /dev/null; then
  echo "  ✅ pgAdmin работает"
else
  echo "  ❌ pgAdmin недоступен"
fi

# Проверка Redis Commander
echo "🔴 Redis Commander:"
if curl -s http://localhost:8084 > /dev/null; then
  echo "  ✅ Redis Commander работает"
else
  echo "  ❌ Redis Commander недоступен"
fi

# Проверка Kafka UI
echo "🛠️ Kafka UI:"
if curl -s http://localhost:8083 > /dev/null; then
  echo "  ✅ Kafka UI работает"
else
  echo "  ❌ Kafka UI недоступен"
fi

echo ""
echo "🎯 Для генерации данных выполните:"
echo "  curl -X POST http://localhost:8080/api/v1/events -H 'Content-Type: application/json' -d '{\"type\":\"user_action\",\"user_id\":\"test\",\"data\":{\"action\":\"test\"}}'"
echo ""
echo "📋 Доступные ссылки:"
echo "  📊 Prometheus:  http://localhost:9090"
echo "  📈 Grafana:     http://localhost:3000 (admin/admin)"
echo "  🔍 Jaeger:      http://localhost:16686"
echo "  🗄️ pgAdmin:     http://localhost:8085 (admin@admin.com/admin)"
echo "  🔴 Redis:       http://localhost:8084"
echo "  🛠️ Kafka UI:    http://localhost:8083"
