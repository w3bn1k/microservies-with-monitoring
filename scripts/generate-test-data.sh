#!/bin/bash
# Скрипт для генерации тестовых данных

echo "🚀 Генерация тестовых данных для всех систем..."

# Генерировать события для заполнения всех систем
echo "📤 Отправка событий..."
for i in {1..10}; do
  echo "  Отправка события $i/10..."
  curl -s -X POST http://localhost:8080/api/v1/events \
    -H "Content-Type: application/json" \
    -d "{\"type\":\"user_action\",\"user_id\":\"user$i\",\"data\":{\"action\":\"test_$i\",\"timestamp\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}}" > /dev/null
  sleep 0.5
done

echo "✅ События отправлены!"
echo ""

# Проверить, что данные появились
echo "🔍 Проверка данных в системах..."

echo "📊 Prometheus метрики:"
prometheus_metrics=$(curl -s http://localhost:9090/api/v1/query?query=events_processed_total | jq -r '.data.result | length' 2>/dev/null || echo "0")
echo "  Метрики: $prometheus_metrics"

echo "🔴 Redis данные:"
redis_keys=$(curl -s http://localhost:8084/api/keys 2>/dev/null | jq -r '.[] | .key' | grep -c "event:" || echo "0")
echo "  Ключи событий: $redis_keys"

echo "🗄️ PostgreSQL данные:"
postgres_count=$(docker-compose exec -T postgres psql -U postgres -d microservices -c "SELECT COUNT(*) FROM transactions;" 2>/dev/null | grep -o '[0-9]\+' | tail -1 || echo "0")
echo "  Транзакций: $postgres_count"

echo "🛠️ Kafka сообщения:"
kafka_messages=$(docker-compose exec kafka kafka-console-consumer --bootstrap-server localhost:29092 --topic user-events --from-beginning --max-messages 10 --timeout-ms 5000 2>/dev/null | wc -l || echo "0")
echo "  Сообщений: $kafka_messages"

echo ""
echo "🎯 Теперь откройте системы мониторинга:"
echo "  📊 Prometheus:  http://localhost:9090"
echo "  📈 Grafana:     http://localhost:3000 (admin/admin)"
echo "  🔍 Jaeger:      http://localhost:16686"
echo "  🗄️ pgAdmin:     http://localhost:8085 (admin@admin.com/admin)"
echo "  🔴 Redis:       http://localhost:8084"
echo "  🛠️ Kafka UI:    http://localhost:8083"
