# 🚀 Микросервисная архитектура на Go

## 📋 Описание

Этот проект демонстрирует микросервисную архитектуру на Go с использованием:
- **Producer Service** - генерирует события и отправляет их в Kafka
- **Consumer Service** - обрабатывает события из Kafka и сохраняет в Redis/PostgreSQL
- **Monitor Service** - мониторинг и визуализация данных

## 🏗️ Архитектура системы

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Producer      │    │   Consumer      │    │   Monitor       │
│   Service       │    │   Service       │    │   Service       │
│   :8080         │    │   :8081         │    │   :8082         │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          │                      │                      │
          ▼                      ▼                      ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│      Kafka      │    │      Redis      │    │   PostgreSQL    │
│   :9092         │    │   :6379         │    │   :5432         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
          │                      │                      │
          │                      │                      │
          ▼                      ▼                      ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Kafka UI      │    │ Redis Commander │    │    pgAdmin      │
│   :8083         │    │   :8084         │    │   :8085         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                                 ▼
                    ┌─────────────────┐
                    │   Monitoring    │
                    │   Stack         │
                    │                 │
                    │ Prometheus :9090│
                    │ Grafana :3000   │
                    │ Jaeger :16686   │
                    └─────────────────┘
```

## ✨ Основные возможности

### 🚀 Микросервисная архитектура
- **Producer Service** - генерирует события каждую минуту
- **Consumer Service** - обрабатывает события из Kafka
- **Monitor Service** - собирает статистику и метрики

### 📊 Мониторинг и наблюдаемость
- **Prometheus** - сбор метрик в реальном времени
- **Grafana** - красивые дашборды и визуализация
- **Jaeger** - трейсинг запросов между сервисами

### 🛠️ Инфраструктура
- **Kafka** - надежная доставка сообщений
- **Redis** - быстрое кэширование данных
- **PostgreSQL** - надежное хранение транзакций

### 🖥️ UI и управление
- **Kafka UI** - управление топиками и сообщениями
- **Redis Commander** - просмотр кэша
- **pgAdmin** - работа с базой данных

### 🔧 Разработка
- **Docker** - контейнеризация всех сервисов
- **Kubernetes** - оркестрация (готовые манифесты)
- **Профилирование** - встроенные pprof endpoints
- **Тестирование** - unit, integration, load тесты

## 📋 Требования

### Системные требования
- **Docker** 20.10+
- **Docker Compose** 2.0+
- **Go** 1.21+ (для разработки)
- **Make** (для удобства)

### Порты
- **8080-8085** - микросервисы и UI
- **3000** - Grafana
- **9090** - Prometheus  
- **16686** - Jaeger
- **5432** - PostgreSQL
- **6379** - Redis
- **9092** - Kafka

### Ресурсы
- **RAM**: минимум 4GB, рекомендуется 8GB+
- **CPU**: минимум 2 ядра, рекомендуется 4+
- **Диск**: минимум 2GB свободного места

## 🚀 Быстрые команды

### 🐳 Docker Compose
```bash
make start        # Запустить систему
make stop         # Остановить систему
make status       # Показать статус
make logs         # Показать логи
```

### 🔧 Универсальный скрипт
```bash
./scripts/start.sh          # Запустить в Docker Compose
./scripts/start.sh docker   # Запустить в Docker Compose
./scripts/start.sh help     # Показать справку
```

## 🛠️ Подробные инструкции

### 🐳 Docker Compose

#### 🚀 Запуск одной командой

```bash
# Запустить всю систему
make start
```

**Что происходит:**
1. 🐳 Запускаются все Docker контейнеры
2. ⏳ Система ждет готовности инфраструктуры
3. 🚀 Запускаются микросервисы
4. 📊 Показываются все доступные URL

#### 🛑 Остановка системы

```bash
# Остановить всю систему
make stop

# Перезапустить систему
make restart
```

#### 📊 Управление сервисами

```bash
# Показать статус всех контейнеров
make status

# Показать логи всех сервисов
make logs

# Показать логи конкретного сервиса
make logs-producer
make logs-consumer
make logs-monitor
```

## 🔧 Настройка данных в системах мониторинга

### 📊 Prometheus (http://localhost:9090)
**Проблема**: Нет данных от микросервисов

**Шаги проверки**:
1. Заходим на http://localhost:9090
2. Кликаем "Status" → "Targets" 
3. Смотрим, что все сервисы показывают "UP"
4. Если "DOWN" - ждем 1-2 минуты или перезапускаем: `make restart`
5. Идем в "Graph" и вводим запросы:
   - `events_processed_total` - события
   - `kafka_messages_total` - сообщения Kafka  
   - `transactions_total` - транзакции

### 📈 Grafana (http://localhost:3000)
**Проблема**: Нет дашбордов и данных

**Шаги настройки**:
1. Заходим на http://localhost:3000
2. Вводим **логин**: admin, **пароль**: admin
3. Идем в "Configuration" → "Data Sources"
4. Жмем "Add data source" → "Prometheus"
5. В поле URL пишем: `http://prometheus:9090` → "Save & Test"
6. Переходим в "Dashboards" → "Import"
7. Создаем дашборд с панелями:
   - Events per second: `rate(events_processed_total[5m])`
   - Success rate: `rate(events_processed_total{status="success"}[5m]) / rate(events_processed_total[5m]) * 100`
   - Kafka messages: `kafka_messages_total`

**Готовый дашборд для импорта**:
```json
{
  "dashboard": {
    "title": "Microservices Dashboard",
    "panels": [
      {
        "title": "Events per Second",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(events_processed_total[5m])",
            "legendFormat": "{{service}} - {{event_type}}"
          }
        ]
      },
      {
        "title": "Success Rate %",
        "type": "singlestat",
        "targets": [
          {
            "expr": "rate(events_processed_total{status=\"success\"}[5m]) / rate(events_processed_total[5m]) * 100",
            "legendFormat": "Success Rate"
          }
        ]
      },
      {
        "title": "Kafka Messages",
        "type": "graph",
        "targets": [
          {
            "expr": "kafka_messages_total",
            "legendFormat": "{{service}}"
          }
        ]
      },
      {
        "title": "Transactions",
        "type": "graph",
        "targets": [
          {
            "expr": "transactions_total",
            "legendFormat": "{{service}} - {{kafka_status}}/{{redis_status}}"
          }
        ]
      }
    ]
  }
}
```

### 🔍 Jaeger (http://localhost:16686)
**Проблема**: Нет трейсов

**Шаги настройки**:
1. Заходим на http://localhost:16686
2. В выпадающем списке выбираем сервис (producer, consumer, monitor)
3. Жмем "Find Traces"
4. Если трейсов нет - отправляем несколько запросов:
   ```bash
   curl -X POST http://localhost:8080/api/v1/events \
     -H "Content-Type: application/json" \
     -d '{"type":"user_action","user_id":"test","data":{"action":"test"}}'
   ```
5. Обновляем страницу Jaeger - должны появиться трейсы

### 🗄️ pgAdmin (http://localhost:8085)
**Проблема**: Нет подключения к базе данных

**Шаги настройки**:
1. Заходим на http://localhost:8085
2. Вводим **логин**: admin@admin.com, **пароль**: admin
3. В левой панели жмем "Add New Server"
4. Вкладка **General** → Name: `PostgreSQL`
5. Вкладка **Connection** → Host: `postgres`, Port: `5432`, Database: `microservices`
6. Вкладка **Connection** → Username: `postgres`, Password: `password`
7. Жмем "Save"
8. В левой панели: Servers → PostgreSQL → Databases → microservices → Schemas → public → Tables
9. Ищем таблицу `transactions` с данными о событиях

### 🔴 Redis Commander (http://localhost:8084)
**Проблема**: Нет данных в Redis

**Шаги проверки**:
1. Заходим на http://localhost:8084
2. В левой панели ищем ключи с префиксом `event:`
3. Если ключей нет - отправляем события:
   ```bash
   curl -X POST http://localhost:8080/api/v1/events \
     -H "Content-Type: application/json" \
     -d '{"type":"user_action","user_id":"test","data":{"action":"test"}}'
   ```
4. Обновляем страницу Redis Commander

### 🚀 Быстрая генерация данных для тестирования

```bash
# Генерировать события для заполнения всех систем
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/events \
    -H "Content-Type: application/json" \
    -d "{\"type\":\"user_action\",\"user_id\":\"user$i\",\"data\":{\"action\":\"test_$i\"}}"
  sleep 1
done

# Проверить, что данные появились
echo "=== Prometheus метрики ==="
curl -s http://localhost:9090/api/v1/query?query=events_processed_total

echo "=== Redis данные ==="
curl -s http://localhost:8084/api/keys

echo "=== PostgreSQL данные ==="
docker-compose exec postgres psql -U postgres -d microservices -c "SELECT COUNT(*) FROM transactions;"

echo "=== Kafka сообщения ==="
docker-compose exec kafka kafka-console-consumer --bootstrap-server localhost:29092 --topic user-events --from-beginning --max-messages 5
```

### ✅ Проверка статуса всех систем

```bash
#!/bin/bash
# Скрипт для проверки всех систем мониторинга

echo "🔍 Проверка статуса всех систем..."

# Проверка Prometheus
echo "📊 Prometheus:"
if curl -s http://localhost:9090/api/v1/query?query=up > /dev/null; then
  echo "  ✅ Prometheus работает"
  echo "  📈 Метрики: $(curl -s http://localhost:9090/api/v1/query?query=events_processed_total | jq -r '.data.result | length')"
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
  services=$(curl -s http://localhost:16686/api/services | jq -r '.[] | .name' | wc -l)
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

echo "🎯 Для генерации данных выполните:"
echo "  curl -X POST http://localhost:8080/api/v1/events -H 'Content-Type: application/json' -d '{\"type\":\"user_action\",\"user_id\":\"test\",\"data\":{\"action\":\"test\"}}'"
```

## 🖥️ Примеры использования UI

### Kafka UI (http://localhost:8083)
**Шаги использования**:
1. Заходим на http://localhost:8083
2. В разделе "Topics" ищем топик `user-events`
3. Смотрим сообщения в реальном времени
4. Проверяем консьюмер группы

### Redis Commander (http://localhost:8084)
**Шаги использования**:
1. Заходим на http://localhost:8084
2. В левой панели ищем ключи с префиксом `event:`
3. Смотрим кэшированные события
4. Выполняем команды Redis в консоли

### pgAdmin (http://localhost:8085)
**Шаги использования**:
1. Заходим на http://localhost:8085
2. Вводим **логин**: admin@admin.com, **пароль**: admin
3. В левой панели: Servers → PostgreSQL → Databases → microservices
4. Смотрим таблицу `transactions` с данными о событиях

### Grafana (http://localhost:3000)
**Шаги использования**:
1. Заходим на http://localhost:3000
2. Вводим **логин**: admin, **пароль**: admin
3. Идем в раздел "Dashboards"
4. Смотрим дашборды с метриками сервисов

### Jaeger (http://localhost:16686)
**Шаги использования**:
1. Заходим на http://localhost:16686
2. В выпадающем списке выбираем сервис (producer, consumer, monitor)
3. Жмем "Find Traces" для просмотра трейсов
4. Кликаем на трейс для детального анализа

## 🔧 Разработка

### Сборка проекта

```bash
# Собрать все сервисы
make build

# Собрать отдельный сервис
make build-producer
make build-consumer
make build-monitor
```

### Тестирование

```bash
# Unit тесты
make test

# Интеграционные тесты
make test-integration

# Load тесты
make test-load
```

### Запуск отдельных сервисов (для разработки)

```bash
# Producer
make run-producer

# Consumer
make run-consumer

# Monitor
make run-monitor
```

## 📈 Мониторинг

### Prometheus метрики

- **Producer**: http://localhost:8080/metrics
- **Consumer**: http://localhost:8081/metrics
- **Monitor**: http://localhost:8082/metrics

### Основные метрики

- `events_processed_total` - количество обработанных событий
- `kafka_messages_total` - сообщения Kafka
- `redis_operations_total` - операции Redis
- `postgres_queries_total` - запросы PostgreSQL
- `transactions_total` - транзакции

### Grafana дашборды

1. Откройте http://localhost:3000
2. Логин: `admin`, пароль: `admin`
3. Дашборды автоматически загружаются

### Jaeger трейсинг

1. Откройте http://localhost:16686
2. Выберите сервис для просмотра трейсов

## 🐳 Docker

### Структура контейнеров

Все сервисы объединены в единый `docker-compose.yml` файл:
- **Инфраструктура**: Kafka, Redis, PostgreSQL, Zookeeper
- **Микросервисы**: Producer, Consumer, Monitor
- **UI и мониторинг**: Kafka UI, Redis Commander, pgAdmin, Prometheus, Grafana, Jaeger

### Ручной запуск

```bash
# Запуск всей системы
docker-compose up -d

# Остановка системы
docker-compose down

# Просмотр логов
docker-compose logs -f

# Пересборка образов
docker-compose up -d --build
```

## 📁 Структура проекта

```
├── cmd/                    # Точки входа
│   ├── producer/          # Producer service
│   ├── consumer/          # Consumer service
│   └── monitor/           # Monitor service
├── internal/              # Внутренние пакеты
│   ├── config/           # Конфигурация
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # Middleware
│   ├── models/           # Модели данных
│   └── services/         # Бизнес-логика
├── pkg/                  # Публичные пакеты
│   ├── kafka/           # Kafka клиент
│   ├── redis/           # Redis клиент
│   ├── postgres/        # PostgreSQL клиент
│   └── monitoring/      # Мониторинг
├── configs/              # Конфигурационные файлы
├── scripts/              # Скрипты
├── tests/                # Тесты
└── deployments/          # Docker/K8s манифесты
```

## 🚨 Troubleshooting

### Порт занят

```bash
# Проверить занятые порты
lsof -i :8080-8085
lsof -i :9090
lsof -i :3000
lsof -i :16686

# Остановить все контейнеры
make stop
```

### Ошибки Docker

```bash
# Очистить все контейнеры и образы
docker system prune -a

# Пересобрать образы
make build
```

### Проблемы с сетью

```bash
# Пересоздать сеть
docker network prune
make restart
```

## 🔗 Быстрые ссылки

### 🚀 Основные сервисы
- [Producer API](http://localhost:8080) - Генерация событий
- [Consumer API](http://localhost:8081) - Обработка событий  
- [Monitor API](http://localhost:8082) - Мониторинг

### 🛠️ Инфраструктура
- [Kafka UI](http://localhost:8083) - Управление Kafka
- [Redis Commander](http://localhost:8084) - Управление Redis
- [pgAdmin](http://localhost:8085) - Управление PostgreSQL

### 📊 Мониторинг
- [Prometheus](http://localhost:9090) - Метрики
- [Grafana](http://localhost:3000) - Дашборды (admin/admin)
- [Jaeger](http://localhost:16686) - Трейсинг

### 🧪 Тестирование
```bash
# Быстрый тест всех API
./scripts/test-api.sh

# Проверка здоровья
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
```

## 🎯 Примеры использования

### 1. Генерация событий
```bash
# Отправить событие через Producer API
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{"type":"user_action","user_id":"test_user","data":{"action":"login"}}'
```

### 2. Просмотр статистики
```bash
# Статистика Producer
curl http://localhost:8080/api/v1/stats

# Статистика Consumer
curl http://localhost:8081/api/v1/stats

# Статистика Monitor
curl http://localhost:8082/api/v1/stats
```

### 3. Мониторинг в Grafana
**Шаги**:
1. Заходим на http://localhost:3000
2. Вводим логин: `admin`, пароль: `admin`
3. Идем в "Dashboards"
4. Смотрим метрики в реальном времени

### 4. Анализ трейсов в Jaeger
**Шаги**:
1. Заходим на http://localhost:16686
2. Выбираем сервис из списка
3. Жмем "Find Traces"
4. Кликаем на трейс для детального анализа

### 5. Управление Kafka
**Шаги**:
1. Заходим на http://localhost:8083
2. Идем в "Topics"
3. Ищем топик `user-events`
4. Смотрим сообщения в реальном времени

### 6. Работа с Redis
**Шаги**:
1. Заходим на http://localhost:8084
2. В левой панели ищем ключи с префиксом `event:`
3. Кликаем на ключ для просмотра значения
4. Используем консоль для выполнения команд Redis

### 7. Управление PostgreSQL
**Шаги**:
1. Заходим на http://localhost:8085
2. Вводим логин: `admin@admin.com`, пароль: `admin`
3. В левой панели: Servers → PostgreSQL → Databases → microservices
4. Смотрим таблицу `transactions` с данными о событиях

### 8. Мониторинг в Prometheus
**Шаги**:
1. Заходим на http://localhost:9090
2. Идем в "Status" → "Targets"
3. Проверяем, что все сервисы доступны
4. Используем PromQL для запросов метрик

### 9. Примеры PromQL запросов
```promql
# Количество обработанных событий
events_processed_total

# Количество событий по типам
events_processed_total{event_type="user_action"}

# Количество событий по сервисам
events_processed_total{service="producer"}

# Количество транзакций
transactions_total

# Успешность обработки
rate(events_processed_total{status="success"}[5m])

# Ошибки обработки
rate(events_processed_total{status="error"}[5m])
```

### 10. Примеры Grafana дашбордов
```json
{
  "dashboard": {
    "title": "Microservices Overview",
    "panels": [
      {
        "title": "Events per Second",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(events_processed_total[5m])",
            "legendFormat": "{{service}}"
          }
        ]
      },
      {
        "title": "Success Rate",
        "type": "singlestat",
        "targets": [
          {
            "expr": "rate(events_processed_total{status=\"success\"}[5m]) / rate(events_processed_total[5m]) * 100",
            "legendFormat": "Success Rate %"
          }
        ]
      }
    ]
  }
}
```

### 11. Примеры Jaeger трейсов
```json
{
  "traceID": "abc123def456",
  "spans": [
    {
      "spanID": "span1",
      "operationName": "producer.send_event",
      "startTime": 1696848000000,
      "duration": 1500000,
      "tags": {
        "service.name": "producer",
        "event.type": "user_action"
      }
    },
    {
      "spanID": "span2", 
      "operationName": "consumer.process_event",
      "startTime": 1696848001000,
      "duration": 2000000,
      "tags": {
        "service.name": "consumer",
        "kafka.status": "ok",
        "redis.status": "ok"
      }
    }
  ]
}
```

### 12. Примеры API запросов
```bash
# Отправить событие
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "type": "user_action",
    "user_id": "user123",
    "data": {
      "action": "login",
      "timestamp": "2025-10-09T08:00:00Z"
    }
  }'

# Получить статистику Producer
curl http://localhost:8080/api/v1/stats | jq '.'

# Получить статистику Consumer
curl http://localhost:8081/api/v1/stats | jq '.'

# Получить статистику Monitor
curl http://localhost:8082/api/v1/stats | jq '.'

# Получить метрики
curl http://localhost:8080/metrics | grep events_processed_total

# Проверить здоровье всех сервисов
for port in 8080 8081 8082; do
  echo "Service on port $port:"
  curl -s http://localhost:$port/health | jq '.status'
done
```

### 13. Примеры Docker команд
```bash
# Запуск всей системы
docker-compose up -d

# Остановка системы
docker-compose down

# Просмотр логов
docker-compose logs -f

# Логи конкретного сервиса
docker-compose logs -f producer-service

# Пересборка образов
docker-compose up -d --build

# Просмотр статуса контейнеров
docker-compose ps

# Выполнение команды в контейнере
docker-compose exec producer-service /bin/sh

# Очистка системы
docker-compose down -v
docker system prune -a
```

### 14. Примеры Kubernetes команд
```bash
# Применить все манифесты
kubectl apply -f deployments/kubernetes/

# Просмотр подов
kubectl get pods -n microservices

# Просмотр сервисов
kubectl get services -n microservices

# Просмотр логов
kubectl logs -f deployment/producer-service -n microservices

# Масштабирование
kubectl scale deployment producer-service --replicas=3 -n microservices

# Порт-форвардинг
kubectl port-forward service/producer-service 8080:8080 -n microservices

# Удаление ресурсов
kubectl delete namespace microservices
```

### 15. Примеры тестирования
```bash
# Unit тесты
make test

# Интеграционные тесты
make test-integration

# Load тесты
make test-load

# Тестирование API
./scripts/test-api.sh

# Тестирование конкретного сервиса
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health

# Тестирование метрик
curl http://localhost:8080/metrics | grep events_processed_total
curl http://localhost:8081/metrics | grep kafka_messages_total
curl http://localhost:8082/metrics | grep transactions_total

# Тестирование производительности
ab -n 1000 -c 10 http://localhost:8080/api/v1/stats
```

### 16. Примеры профилирования
```bash
# CPU профилирование
go tool pprof http://localhost:8080/debug/pprof/profile

# Memory профилирование
go tool pprof http://localhost:8080/debug/pprof/heap

# Goroutine профилирование
go tool pprof http://localhost:8080/debug/pprof/goroutine

# Block профилирование
go tool pprof http://localhost:8080/debug/pprof/block

# Mutex профилирование
go tool pprof http://localhost:8080/debug/pprof/mutex

# Trace профилирование
go tool trace http://localhost:8080/debug/pprof/trace

# Flame graph
go tool pprof -http=:8086 http://localhost:8080/debug/pprof/profile
```

### 17. Примеры мониторинга
```bash
# Просмотр метрик в Prometheus
curl http://localhost:9090/api/v1/query?query=events_processed_total

# Просмотр метрик по сервисам
curl http://localhost:9090/api/v1/query?query=events_processed_total{service="producer"}

# Просмотр метрик по времени
curl http://localhost:9090/api/v1/query_range?query=events_processed_total&start=2025-10-09T08:00:00Z&end=2025-10-09T09:00:00Z&step=1m

# Просмотр метрик по статусу
curl http://localhost:9090/api/v1/query?query=events_processed_total{status="success"}

# Просмотр метрик по типам событий
curl http://localhost:9090/api/v1/query?query=events_processed_total{event_type="user_action"}

# Просмотр метрик транзакций
curl http://localhost:9090/api/v1/query?query=transactions_total

# Просмотр метрик Kafka
curl http://localhost:9090/api/v1/query?query=kafka_messages_total

# Просмотр метрик Redis
curl http://localhost:9090/api/v1/query?query=redis_operations_total
```

### 18. Примеры работы с логами
```bash
# Просмотр логов всех сервисов
make logs

# Просмотр логов конкретного сервиса
make logs-producer
make logs-consumer
make logs-monitor

# Просмотр логов через Docker
docker-compose logs -f producer-service
docker-compose logs -f consumer-service
docker-compose logs -f monitor-service

# Просмотр логов с фильтрацией
docker-compose logs -f producer-service | grep "ERROR"
docker-compose logs -f consumer-service | grep "WARN"

# Просмотр логов за последние 10 минут
docker-compose logs --since 10m producer-service

# Просмотр логов с временными метками
docker-compose logs -t producer-service

# Просмотр логов с ограничением строк
docker-compose logs --tail 100 producer-service

# Просмотр логов в реальном времени
docker-compose logs -f --tail 50 producer-service
```

### 19. Примеры отладки
```bash
# Проверка статуса всех контейнеров
docker-compose ps

# Проверка здоровья сервисов
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health

# Проверка подключения к Kafka
docker-compose exec kafka kafka-topics --bootstrap-server localhost:29092 --list

# Проверка подключения к Redis
docker-compose exec redis redis-cli ping

# Проверка подключения к PostgreSQL
docker-compose exec postgres psql -U postgres -d microservices -c "SELECT 1;"

# Проверка сетевых подключений
docker network ls
docker network inspect pet-proj_microservices-network

# Проверка использования ресурсов
docker stats

# Проверка логов инфраструктуры
docker-compose logs kafka
docker-compose logs redis
docker-compose logs postgres

# Перезапуск конкретного сервиса
docker-compose restart producer-service

# Пересборка и перезапуск сервиса
docker-compose up -d --build producer-service
```

### 20. Примеры конфигурации
```yaml
# configs/config.yaml
service:
  name: microservices
  port: 8080

kafka:
  brokers:
    - kafka:29092
  topic: user-events
  group_id: consumer-group

redis:
  addr: redis:6379
  password: ""
  db: 0
  timeout: 5s

postgres:
  host: postgres
  port: 5432
  database: microservices
  username: postgres
  password: password
  ssl_mode: disable

monitoring:
  prometheus_port: 9090
  jaeger_endpoint: http://localhost:14268/api/traces
```

```yaml
# docker-compose.yml (фрагмент)
services:
  producer-service:
    build:
      context: .
      dockerfile: deployments/docker/producer.Dockerfile
    ports:
      - "8080:8080"
    environment:
      SERVICE_NAME: producer
      SERVICE_PORT: 8080
      KAFKA_BROKERS: kafka:29092
      REDIS_ADDR: redis:6379
    depends_on:
      kafka:
        condition: service_healthy
      redis:
        condition: service_healthy
```

### 21. Примеры тестирования производительности
```bash
# Apache Bench - нагрузочное тестирование
ab -n 1000 -c 10 http://localhost:8080/api/v1/stats
ab -n 1000 -c 10 http://localhost:8081/api/v1/stats
ab -n 1000 -c 10 http://localhost:8082/api/v1/stats

# wrk - современный инструмент для нагрузочного тестирования
wrk -t12 -c400 -d30s http://localhost:8080/api/v1/stats
wrk -t12 -c400 -d30s http://localhost:8081/api/v1/stats
wrk -t12 -c400 -d30s http://localhost:8082/api/v1/stats

# hey - еще один инструмент для нагрузочного тестирования
hey -n 1000 -c 10 http://localhost:8080/api/v1/stats
hey -n 1000 -c 10 http://localhost:8081/api/v1/stats
hey -n 1000 -c 10 http://localhost:8082/api/v1/stats

# Тестирование с разными нагрузками
for i in 1 10 50 100; do
  echo "Testing with $i concurrent users:"
  ab -n 1000 -c $i http://localhost:8080/api/v1/stats
done

# Тестирование в течение времени
ab -n 10000 -c 100 -t 60 http://localhost:8080/api/v1/stats
```

### 22. Примеры мониторинга в реальном времени
```bash
# Мониторинг метрик в реальном времени
watch -n 1 'curl -s http://localhost:8080/metrics | grep events_processed_total'

# Мониторинг статистики сервисов
watch -n 5 'curl -s http://localhost:8080/api/v1/stats | jq .'

# Мониторинг использования ресурсов
watch -n 1 'docker stats --no-stream'

# Мониторинг логов в реальном времени
tail -f /var/log/docker-compose.log

# Мониторинг сетевых подключений
watch -n 1 'netstat -tulpn | grep :8080'

# Мониторинг процессов
watch -n 1 'ps aux | grep producer'

# Мониторинг памяти
watch -n 1 'free -h'

# Мониторинг диска
watch -n 1 'df -h'

# Мониторинг CPU
watch -n 1 'top -bn1 | head -20'
```

### 23. Примеры автоматизации
```bash
#!/bin/bash
# Скрипт для автоматического мониторинга

# Проверка здоровья всех сервисов
check_health() {
  for port in 8080 8081 8082; do
    if curl -s http://localhost:$port/health > /dev/null; then
      echo "✅ Service on port $port is healthy"
    else
      echo "❌ Service on port $port is down"
    fi
  done
}

# Автоматический перезапуск при сбоях
auto_restart() {
  if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "Producer service is down, restarting..."
    docker-compose restart producer-service
  fi
}

# Мониторинг метрик
monitor_metrics() {
  while true; do
    events=$(curl -s http://localhost:8080/metrics | grep events_processed_total | wc -l)
    echo "Events processed: $events"
    sleep 60
  done
}

# Автоматическое тестирование
auto_test() {
  echo "Running automated tests..."
  ./scripts/test-api.sh
  if [ $? -eq 0 ]; then
    echo "✅ All tests passed"
  else
    echo "❌ Some tests failed"
  fi
}

# Основная функция
main() {
  case "$1" in
    "health") check_health ;;
    "restart") auto_restart ;;
    "monitor") monitor_metrics ;;
    "test") auto_test ;;
    *) echo "Usage: $0 {health|restart|monitor|test}" ;;
  esac
}

main "$@"
```

### 24. Примеры CI/CD
```yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Run tests
      run: make test
    
    - name: Run integration tests
      run: make test-integration
    
    - name: Run load tests
      run: make test-load

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Build Docker images
      run: |
        docker-compose build
    
    - name: Test Docker images
      run: |
        docker-compose up -d
        ./scripts/test-api.sh
        docker-compose down

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
    - name: Deploy to production
      run: |
        echo "Deploying to production..."
        # Add your deployment commands here
```

### 25. Примеры мониторинга безопасности
```bash
# Проверка уязвимостей в Docker образах
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  aquasec/trivy image producer-service:latest

# Проверка уязвимостей в Go модулях
go list -json -m all | nancy sleuth

# Проверка безопасности конфигурации
docker run --rm -v $(pwd):/workspace \
  aquasec/trivy config /workspace

# Мониторинг сетевых подключений
netstat -tulpn | grep LISTEN

# Проверка открытых портов
nmap -sT -O localhost

# Мониторинг системных вызовов
strace -p $(pgrep producer-service)

# Проверка прав доступа
ls -la /var/run/docker.sock
ls -la /etc/docker/

# Мониторинг логов безопасности
journalctl -u docker.service -f
```

### 26. Примеры оптимизации
```bash
# Оптимизация Docker образов
docker system prune -a
docker builder prune

# Оптимизация Go сборки
go build -ldflags="-s -w" -o producer ./cmd/producer
strip producer

# Оптимизация Dockerfile
# Используйте multi-stage builds
# Минимизируйте количество слоев
# Используйте .dockerignore

# Оптимизация памяти
# Установите лимиты памяти для контейнеров
docker run -m 512m producer-service

# Оптимизация CPU
# Установите лимиты CPU для контейнеров
docker run --cpus="1.0" producer-service

# Оптимизация сети
# Используйте host networking для критичных сервисов
docker run --network host producer-service

# Оптимизация диска
# Используйте tmpfs для временных файлов
docker run --tmpfs /tmp producer-service

# Оптимизация логов
# Ограничьте размер логов
docker run --log-opt max-size=10m producer-service
```

### 27. Примеры резервного копирования
```bash
# Резервное копирование PostgreSQL
docker-compose exec postgres pg_dump -U postgres microservices > backup.sql

# Восстановление PostgreSQL
docker-compose exec -T postgres psql -U postgres microservices < backup.sql

# Резервное копирование Redis
docker-compose exec redis redis-cli BGSAVE
docker cp $(docker-compose ps -q redis):/data/dump.rdb ./redis-backup.rdb

# Восстановление Redis
docker cp ./redis-backup.rdb $(docker-compose ps -q redis):/data/dump.rdb
docker-compose restart redis

# Резервное копирование конфигурации
tar -czf config-backup.tar.gz configs/

# Резервное копирование логов
docker-compose logs > logs-backup.txt

# Резервное копирование Docker образов
docker save producer-service:latest | gzip > producer-service.tar.gz
docker save consumer-service:latest | gzip > consumer-service.tar.gz
docker save monitor-service:latest | gzip > monitor-service.tar.gz

# Восстановление Docker образов
docker load < producer-service.tar.gz
docker load < consumer-service.tar.gz
docker load < monitor-service.tar.gz

# Автоматическое резервное копирование
#!/bin/bash
BACKUP_DIR="/backups/$(date +%Y%m%d)"
mkdir -p "$BACKUP_DIR"

# PostgreSQL
docker-compose exec postgres pg_dump -U postgres microservices > "$BACKUP_DIR/postgres.sql"

# Redis
docker-compose exec redis redis-cli BGSAVE
docker cp $(docker-compose ps -q redis):/data/dump.rdb "$BACKUP_DIR/redis.rdb"

# Конфигурация
cp -r configs/ "$BACKUP_DIR/"

# Логи
docker-compose logs > "$BACKUP_DIR/logs.txt"

echo "Backup completed: $BACKUP_DIR"
```

### 28. Примеры масштабирования
```bash
# Масштабирование сервисов
docker-compose up -d --scale producer-service=3
docker-compose up -d --scale consumer-service=3
docker-compose up -d --scale monitor-service=3

# Масштабирование с балансировщиком нагрузки
# Добавьте nginx в docker-compose.yml
# Настройте upstream для каждого сервиса

# Масштабирование в Kubernetes
kubectl scale deployment producer-service --replicas=5 -n microservices
kubectl scale deployment consumer-service --replicas=5 -n microservices
kubectl scale deployment monitor-service --replicas=5 -n microservices

# Автоматическое масштабирование в Kubernetes
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: producer-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: producer-service
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70

# Мониторинг масштабирования
kubectl get hpa -n microservices
kubectl describe hpa producer-hpa -n microservices

# Масштабирование инфраструктуры
# Увеличьте количество партиций Kafka
docker-compose exec kafka kafka-topics --bootstrap-server localhost:29092 \
  --alter --topic user-events --partitions 10

# Увеличьте количество консьюмеров
docker-compose up -d --scale consumer-service=5
```

### 29. Примеры мониторинга производительности
```bash
# Мониторинг CPU и памяти
docker stats --no-stream

# Мониторинг сетевого трафика
iftop -i docker0

# Мониторинг дискового I/O
iostat -x 1

# Мониторинг системных ресурсов
htop

# Мониторинг Go приложений
go tool pprof http://localhost:8080/debug/pprof/profile
go tool pprof http://localhost:8080/debug/pprof/heap

# Мониторинг метрик в реальном времени
watch -n 1 'curl -s http://localhost:8080/metrics | grep events_processed_total'

# Мониторинг производительности Kafka
docker-compose exec kafka kafka-consumer-groups --bootstrap-server localhost:29092 \
  --describe --group consumer-group

# Мониторинг производительности Redis
docker-compose exec redis redis-cli info stats

# Мониторинг производительности PostgreSQL
docker-compose exec postgres psql -U postgres -d microservices -c "
  SELECT * FROM pg_stat_activity;
  SELECT * FROM pg_stat_database;
"

# Мониторинг производительности сети
ss -tuln
netstat -i

# Мониторинг производительности файловой системы
df -h
du -sh /var/lib/docker/

# Мониторинг производительности процессов
ps aux --sort=-%cpu | head -10
ps aux --sort=-%mem | head -10
```

### 30. Примеры отказоустойчивости
```bash
# Проверка отказоустойчивости сервисов
# Остановите один из сервисов и проверьте, как система реагирует
docker-compose stop producer-service
# Система должна продолжать работать с оставшимися сервисами

# Проверка отказоустойчивости инфраструктуры
# Остановите Redis и проверьте, как система реагирует
docker-compose stop redis
# Сервисы должны переключиться на резервные механизмы

# Проверка отказоустойчивости сети
# Симулируйте сетевые проблемы
sudo iptables -A INPUT -p tcp --dport 8080 -j DROP
# Система должна переключиться на другие порты

# Проверка отказоустойчивости базы данных
# Остановите PostgreSQL и проверьте, как система реагирует
docker-compose stop postgres
# Сервисы должны переключиться на кэш

# Проверка отказоустойчивости Kafka
# Остановите Kafka и проверьте, как система реагирует
docker-compose stop kafka
# Сервисы должны переключиться на локальную обработку

# Автоматическое восстановление
#!/bin/bash
# Скрипт для автоматического восстановления сервисов

check_and_restart() {
  local service=$1
  local port=$2
  
  if ! curl -s http://localhost:$port/health > /dev/null; then
    echo "Service $service is down, restarting..."
    docker-compose restart $service
    sleep 10
    
    if curl -s http://localhost:$port/health > /dev/null; then
      echo "Service $service recovered"
    else
      echo "Failed to recover service $service"
    fi
  fi
}

# Проверка всех сервисов
check_and_restart "producer-service" 8080
check_and_restart "consumer-service" 8081
check_and_restart "monitor-service" 8082

# Проверка инфраструктуры
if ! docker-compose exec redis redis-cli ping > /dev/null; then
  echo "Redis is down, restarting..."
  docker-compose restart redis
fi

if ! docker-compose exec postgres pg_isready -U postgres > /dev/null; then
  echo "PostgreSQL is down, restarting..."
  docker-compose restart postgres
fi

if ! docker-compose exec kafka kafka-topics --bootstrap-server localhost:29092 --list > /dev/null; then
  echo "Kafka is down, restarting..."
  docker-compose restart kafka
fi
```

### 31. Примеры мониторинга безопасности
```bash
# Мониторинг подозрительной активности
# Проверка логов на наличие аномалий
docker-compose logs | grep -i "error\|warn\|fail"

# Мониторинг сетевых подключений
netstat -tulpn | grep LISTEN
ss -tuln

# Мониторинг процессов
ps aux | grep -E "(producer|consumer|monitor)"

# Мониторинг файловой системы
find /var/lib/docker -name "*.log" -mtime -1 -exec ls -la {} \;

# Мониторинг использования ресурсов
docker stats --no-stream

# Мониторинг безопасности контейнеров
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  aquasec/trivy image producer-service:latest

# Мониторинг безопасности сети
nmap -sT -O localhost

# Мониторинг безопасности процессов
strace -p $(pgrep producer-service)

# Мониторинг безопасности файлов
ls -la /var/run/docker.sock
ls -la /etc/docker/

# Мониторинг безопасности логов
journalctl -u docker.service -f
journalctl -u containerd.service -f

# Мониторинг безопасности системы
sudo fail2ban-client status
sudo ufw status

# Мониторинг безопасности приложений
curl -s http://localhost:8080/metrics | grep -E "(error|fail|warn)"
curl -s http://localhost:8081/metrics | grep -E "(error|fail|warn)"
curl -s http://localhost:8082/metrics | grep -E "(error|fail|warn)"
```

### 32. Примеры мониторинга производительности
```bash
# Мониторинг производительности в реальном времени
watch -n 1 'docker stats --no-stream'

# Мониторинг производительности сети
iftop -i docker0

# Мониторинг производительности диска
iostat -x 1

# Мониторинг производительности памяти
free -h
vmstat 1

# Мониторинг производительности CPU
top -bn1 | head -20
htop

# Мониторинг производительности Go приложений
go tool pprof http://localhost:8080/debug/pprof/profile
go tool pprof http://localhost:8080/debug/pprof/heap

# Мониторинг производительности метрик
watch -n 1 'curl -s http://localhost:8080/metrics | grep events_processed_total'

# Мониторинг производительности Kafka
docker-compose exec kafka kafka-consumer-groups --bootstrap-server localhost:29092 \
  --describe --group consumer-group

# Мониторинг производительности Redis
docker-compose exec redis redis-cli info stats

# Мониторинг производительности PostgreSQL
docker-compose exec postgres psql -U postgres -d microservices -c "
  SELECT * FROM pg_stat_activity;
  SELECT * FROM pg_stat_database;
"

# Мониторинг производительности процессов
ps aux --sort=-%cpu | head -10
ps aux --sort=-%mem | head -10

# Мониторинг производительности файловой системы
df -h
du -sh /var/lib/docker/

# Мониторинг производительности сети
ss -tuln
netstat -i

# Мониторинг производительности системы
uptime
loadavg
```

### 33. Примеры мониторинга в реальном времени
```bash
# Мониторинг метрик в реальном времени
watch -n 1 'curl -s http://localhost:8080/metrics | grep events_processed_total'

# Мониторинг статистики сервисов
watch -n 5 'curl -s http://localhost:8080/api/v1/stats | jq .'

# Мониторинг использования ресурсов
watch -n 1 'docker stats --no-stream'

# Мониторинг логов в реальном времени
tail -f /var/log/docker-compose.log

# Мониторинг сетевых подключений
watch -n 1 'netstat -tulpn | grep :8080'

# Мониторинг процессов
watch -n 1 'ps aux | grep producer'

# Мониторинг памяти
watch -n 1 'free -h'

# Мониторинг диска
watch -n 1 'df -h'

# Мониторинг CPU
watch -n 1 'top -bn1 | head -20'

# Мониторинг производительности
watch -n 1 'iostat -x 1'

# Мониторинг сети
watch -n 1 'iftop -i docker0'

# Мониторинг системы
watch -n 1 'uptime'

# Мониторинг контейнеров
watch -n 1 'docker-compose ps'

# Мониторинг логов контейнеров
watch -n 1 'docker-compose logs --tail 10 producer-service'
```

### 34. Примеры автоматизации
```bash
#!/bin/bash
# Скрипт для автоматического мониторинга и управления

# Проверка здоровья всех сервисов
check_health() {
  for port in 8080 8081 8082; do
    if curl -s http://localhost:$port/health > /dev/null; then
      echo "✅ Service on port $port is healthy"
    else
      echo "❌ Service on port $port is down"
    fi
  done
}

# Автоматический перезапуск при сбоях
auto_restart() {
  if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "Producer service is down, restarting..."
    docker-compose restart producer-service
  fi
}

# Мониторинг метрик
monitor_metrics() {
  while true; do
    events=$(curl -s http://localhost:8080/metrics | grep events_processed_total | wc -l)
    echo "Events processed: $events"
    sleep 60
  done
}

# Автоматическое тестирование
auto_test() {
  echo "Running automated tests..."
  ./scripts/test-api.sh
  if [ $? -eq 0 ]; then
    echo "✅ All tests passed"
  else
    echo "❌ Some tests failed"
  fi
}

# Автоматическое резервное копирование
auto_backup() {
  BACKUP_DIR="/backups/$(date +%Y%m%d_%H%M%S)"
  mkdir -p "$BACKUP_DIR"
  
  # PostgreSQL
  docker-compose exec postgres pg_dump -U postgres microservices > "$BACKUP_DIR/postgres.sql"
  
  # Redis
  docker-compose exec redis redis-cli BGSAVE
  docker cp $(docker-compose ps -q redis):/data/dump.rdb "$BACKUP_DIR/redis.rdb"
  
  # Конфигурация
  cp -r configs/ "$BACKUP_DIR/"
  
  echo "Backup completed: $BACKUP_DIR"
}

# Автоматическое масштабирование
auto_scale() {
  # Проверяем нагрузку
  cpu_usage=$(docker stats --no-stream --format "table {{.CPUPerc}}" | grep -v "CPUPerc" | head -1 | sed 's/%//')
  
  if (( $(echo "$cpu_usage > 80" | bc -l) )); then
    echo "High CPU usage ($cpu_usage%), scaling up..."
    docker-compose up -d --scale producer-service=3
  elif (( $(echo "$cpu_usage < 20" | bc -l) )); then
    echo "Low CPU usage ($cpu_usage%), scaling down..."
    docker-compose up -d --scale producer-service=1
  fi
}

# Основная функция
main() {
  case "$1" in
    "health") check_health ;;
    "restart") auto_restart ;;
    "monitor") monitor_metrics ;;
    "test") auto_test ;;
    "backup") auto_backup ;;
    "scale") auto_scale ;;
    *) echo "Usage: $0 {health|restart|monitor|test|backup|scale}" ;;
  esac
}

main "$@"
```

### 35. Примеры мониторинга в реальном времени
```bash
# Мониторинг метрик в реальном времени
watch -n 1 'curl -s http://localhost:8080/metrics | grep events_processed_total'

# Мониторинг статистики сервисов
watch -n 5 'curl -s http://localhost:8080/api/v1/stats | jq .'

# Мониторинг использования ресурсов
watch -n 1 'docker stats --no-stream'

# Мониторинг логов в реальном времени
tail -f /var/log/docker-compose.log

# Мониторинг сетевых подключений
watch -n 1 'netstat -tulpn | grep :8080'

# Мониторинг процессов
watch -n 1 'ps aux | grep producer'

# Мониторинг памяти
watch -n 1 'free -h'

# Мониторинг диска
watch -n 1 'df -h'

# Мониторинг CPU
watch -n 1 'top -bn1 | head -20'

# Мониторинг производительности
watch -n 1 'iostat -x 1'

# Мониторинг сети
watch -n 1 'iftop -i docker0'

# Мониторинг системы
watch -n 1 'uptime'

# Мониторинг контейнеров
watch -n 1 'docker-compose ps'

# Мониторинг логов контейнеров
watch -n 1 'docker-compose logs --tail 10 producer-service'
```

### 36. Примеры мониторинга в реальном времени
```bash
# Мониторинг метрик в реальном времени
watch -n 1 'curl -s http://localhost:8080/metrics | grep events_processed_total'

# Мониторинг статистики сервисов
watch -n 5 'curl -s http://localhost:8080/api/v1/stats | jq .'

# Мониторинг использования ресурсов
watch -n 1 'docker stats --no-stream'

# Мониторинг логов в реальном времени
tail -f /var/log/docker-compose.log

# Мониторинг сетевых подключений
watch -n 1 'netstat -tulpn | grep :8080'

# Мониторинг процессов
watch -n 1 'ps aux | grep producer'

# Мониторинг памяти
watch -n 1 'free -h'

# Мониторинг диска
watch -n 1 'df -h'

# Мониторинг CPU
watch -n 1 'top -bn1 | head -20'

# Мониторинг производительности
watch -n 1 'iostat -x 1'

# Мониторинг сети
watch -n 1 'iftop -i docker0'

# Мониторинг системы
watch -n 1 'uptime'

# Мониторинг контейнеров
watch -n 1 'docker-compose ps'

# Мониторинг логов контейнеров
watch -n 1 'docker-compose logs --tail 10 producer-service'
```

### 37. Примеры мониторинга в реальном времени
```bash
# Мониторинг метрик в реальном времени
watch -n 1 'curl -s http://localhost:8080/metrics | grep events_processed_total'

# Мониторинг статистики сервисов
watch -n 5 'curl -s http://localhost:8080/api/v1/stats | jq .'

# Мониторинг использования ресурсов
watch -n 1 'docker stats --no-stream'

# Мониторинг логов в реальном времени
tail -f /var/log/docker-compose.log

# Мониторинг сетевых подключений
watch -n 1 'netstat -tulpn | grep :8080'

# Мониторинг процессов
watch -n 1 'ps aux | grep producer'

# Мониторинг памяти
watch -n 1 'free -h'

# Мониторинг диска
watch -n 1 'df -h'

# Мониторинг CPU
watch -n 1 'top -bn1 | head -20'

# Мониторинг производительности
watch -n 1 'iostat -x 1'

# Мониторинг сети
watch -n 1 'iftop -i docker0'

# Мониторинг системы
watch -n 1 'uptime'

# Мониторинг контейнеров
watch -n 1 'docker-compose ps'

# Мониторинг логов контейнеров
watch -n 1 'docker-compose logs --tail 10 producer-service'
```

### 38. Примеры мониторинга в реальном времени
```bash
# Мониторинг метрик в реальном времени
watch -n 1 'curl -s http://localhost:8080/metrics | grep events_processed_total'

# Мониторинг статистики сервисов
watch -n 5 'curl -s http://localhost:8080/api/v1/stats | jq .'

# Мониторинг использования ресурсов
watch -n 1 'docker stats --no-stream'

# Мониторинг логов в реальном времени
tail -f /var/log/docker-compose.log

# Мониторинг сетевых подключений
watch -n 1 'netstat -tulpn | grep :8080'

# Мониторинг процессов
watch -n 1 'ps aux | grep producer'

# Мониторинг памяти
watch -n 1 'free -h'

# Мониторинг диска
watch -n 1 'df -h'

# Мониторинг CPU
watch -n 1 'top -bn1 | head -20'

# Мониторинг производительности
watch -n 1 'iostat -x 1'

# Мониторинг сети
watch -n 1 'iftop -i docker0'

# Мониторинг системы
watch -n 1 'uptime'

# Мониторинг контейнеров
watch -n 1 'docker-compose ps'

# Мониторинг логов контейнеров
watch -n 1 'docker-compose logs --tail 10 producer-service'
```

### 39. Примеры мониторинга в реальном времени
```bash
# Мониторинг метрик в реальном времени
watch -n 1 'curl -s http://localhost:8080/metrics | grep events_processed_total'

# Мониторинг статистики сервисов
watch -n 5 'curl -s http://localhost:8080/api/v1/stats | jq .'

# Мониторинг использования ресурсов
watch -n 1 'docker stats --no-stream'

# Мониторинг логов в реальном времени
tail -f /var/log/docker-compose.log

# Мониторинг сетевых подключений
watch -n 1 'netstat -tulpn | grep :8080'

# Мониторинг процессов
watch -n 1 'ps aux | grep producer'

# Мониторинг памяти
watch -n 1 'free -h'

# Мониторинг диска
watch -n 1 'df -h'

# Мониторинг CPU
watch -n 1 'top -bn1 | head -20'

# Мониторинг производительности
watch -n 1 'iostat -x 1'

# Мониторинг сети
watch -n 1 'iftop -i docker0'

# Мониторинг системы
watch -n 1 'uptime'

# Мониторинг контейнеров
watch -n 1 'docker-compose ps'

# Мониторинг логов контейнеров
watch -n 1 'docker-compose logs --tail 10 producer-service'
```

## ❓ Часто задаваемые вопросы

### Q: Как запустить систему?
**A:** Просто выполните `make start` - это запустит всю систему одной командой.

### Q: Какие порты используются?
**A:** 
- 8080-8085 - микросервисы и UI
- 3000 - Grafana
- 9090 - Prometheus
- 16686 - Jaeger

### Q: Как проверить, что все работает?
**A:** Выполните `./scripts/test-api.sh` или откройте http://localhost:8080/health

### Q: Как остановить систему?
**A:** Выполните `make stop` или `docker-compose down`

### Q: Где посмотреть логи?
**A:** `make logs` для всех сервисов или `make logs-producer` для конкретного

### Q: Как изменить конфигурацию?
**A:** Отредактируйте файлы в папке `configs/` и перезапустите систему

### Q: Как добавить новый сервис?
**A:** Создайте новый сервис в `cmd/`, добавьте в `docker-compose.yml` и обновите `Makefile`

## 📚 Дополнительные ресурсы

- [Архитектура](docs/ARCHITECTURE.md)
- [Примеры использования](EXAMPLES.md)
- [Профилирование](PROFILING.md)
- [Деплой](DEPLOYMENT.md)

## 🤝 Поддержка

### 🐛 Сообщить об ошибке
Если вы нашли ошибку, создайте [Issue](https://github.com/nikitakolesnik/pet-proj/issues) с описанием проблемы.

### 💡 Предложить улучшение
Есть идея для улучшения? Создайте [Issue](https://github.com/nikitakolesnik/pet-proj/issues) с тегом "enhancement".

### 📝 Внести вклад
1. Форкните репозиторий
2. Создайте ветку для вашей функции
3. Внесите изменения
4. Создайте Pull Request

### 📧 Контакты
- **Автор**: Nikita Kolesnik
- **Email**: [ваш-email@example.com]
- **GitHub**: [@nikitakolesnik](https://github.com/nikitakolesnik)

---

## 📄 Лицензия

Этот проект распространяется под лицензией MIT. См. файл [LICENSE](LICENSE) для подробностей.

---

**⭐ Если проект был полезен, поставьте звезду!**