# Makefile для микросервисного приложения

.PHONY: help build test clean docker-build docker-run

# Переменные
APP_NAME := pet-proj
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "latest")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | awk '{print $$3}')

# Цвета для вывода
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

help: ## Показать справку
	@echo "$(BLUE)Доступные команды:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Go команды
build: ## Собрать все сервисы
	@echo "$(BLUE)Сборка всех сервисов...$(NC)"
	@mkdir -p bin
	@echo "$(YELLOW)Сборка Producer Service...$(NC)"
	@go build -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)" -o bin/producer-service ./cmd/producer
	@echo "$(YELLOW)Сборка Consumer Service...$(NC)"
	@go build -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)" -o bin/consumer-service ./cmd/consumer
	@echo "$(YELLOW)Сборка Monitor Service...$(NC)"
	@go build -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)" -o bin/monitor-service ./cmd/monitor
	@echo "$(GREEN)Сборка завершена!$(NC)"

run-producer: ## Запустить Producer Service
	@echo "$(BLUE)Запуск Producer Service...$(NC)"
	@go run ./cmd/producer

run-consumer: ## Запустить Consumer Service
	@echo "$(BLUE)Запуск Consumer Service...$(NC)"
	@go run ./cmd/consumer

run-monitor: ## Запустить Monitor Service
	@echo "$(BLUE)Запуск Monitor Service...$(NC)"
	@go run ./cmd/monitor

run-gateway: ## Запустить API Gateway
	@echo "$(BLUE)Запуск API Gateway...$(NC)"
	@go run ./cmd/api-gateway

# Production commands
start: generate-proto ## Запустить всю систему (Docker Compose)
	@./scripts/start.sh docker

start-docker: generate-proto ## Запустить всю систему в Docker Compose
	@./scripts/start.sh docker

stop: ## Остановить всю систему
	@echo "$(BLUE)🛑 Остановка системы...$(NC)"
	@docker-compose down
	@echo "$(GREEN)✅ Система остановлена!$(NC)"

restart: stop start ## Перезапустить всю систему

status: ## Показать статус всех сервисов
	@echo "$(BLUE)📊 Статус сервисов:$(NC)"
	@docker-compose ps

logs: ## Показать логи всех сервисов
	@echo "$(BLUE)📋 Логи сервисов:$(NC)"
	@docker-compose logs -f

logs-producer: ## Показать логи Producer
	@docker-compose logs -f producer-service

logs-consumer: ## Показать логи Consumer
	@docker-compose logs -f consumer-service

logs-monitor: ## Показать логи Monitor
	@docker-compose logs -f monitor-service

check-monitoring: ## Проверить статус всех систем мониторинга
	@./scripts/check-monitoring.sh

generate-data: ## Сгенерировать тестовые данные
	@./scripts/generate-test-data.sh


# Тестирование
test: ## Запустить все тесты
	@echo "$(BLUE)Запуск тестов...$(NC)"
	@go test -v ./...

test-unit: ## Запустить unit тесты
	@echo "$(BLUE)Запуск unit тестов...$(NC)"
	@go test -v ./internal/...

test-integration: ## Запустить integration тесты
	@echo "$(BLUE)Запуск integration тестов...$(NC)"
	@go test -v ./tests/integration/...

test-e2e: ## Запустить e2e тесты
	@echo "$(BLUE)Запуск e2e тестов...$(NC)"
	@go test -v ./tests/e2e/...

test-coverage: ## Запустить тесты с покрытием
	@echo "$(BLUE)Запуск тестов с покрытием...$(NC)"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Отчет о покрытии сохранен в coverage.html$(NC)"

benchmark: ## Запустить бенчмарки
	@echo "$(BLUE)Запуск бенчмарков...$(NC)"
	@go test -bench=. -benchmem ./internal/...

# Линтинг и форматирование
fmt: ## Форматировать код
	@echo "$(BLUE)Форматирование кода...$(NC)"
	@go fmt ./...

lint: ## Запустить линтер
	@echo "$(BLUE)Запуск линтера...$(NC)"
	@golangci-lint run

lint-fix: ## Исправить ошибки линтера
	@echo "$(BLUE)Исправление ошибок линтера...$(NC)"
	@golangci-lint run --fix

# Зависимости
deps: ## Установить зависимости
	@echo "$(BLUE)Установка зависимостей...$(NC)"
	@go mod download
	@go mod tidy

deps-update: ## Обновить зависимости
	@echo "$(BLUE)Обновление зависимостей...$(NC)"
	@go get -u ./...
	@go mod tidy

# Docker команды
docker-build: ## Собрать Docker образы
	@echo "$(BLUE)Сборка Docker образов...$(NC)"
	@docker build -t producer-service:$(VERSION) -f deployments/docker/producer.Dockerfile .
	@docker build -t consumer-service:$(VERSION) -f deployments/docker/consumer.Dockerfile .
	@docker build -t monitor-service:$(VERSION) -f deployments/docker/monitor.Dockerfile .
	@echo "$(GREEN)Docker образы собраны!$(NC)"

docker-run: ## Запустить в Docker
	@echo "$(BLUE)Запуск в Docker...$(NC)"
	@docker-compose up -d

docker-stop: ## Остановить Docker контейнеры
	@echo "$(BLUE)Остановка Docker контейнеров...$(NC)"
	@docker-compose down

docker-logs: ## Показать логи Docker контейнеров
	@echo "$(BLUE)Логи Docker контейнеров:$(NC)"
	@docker-compose logs -f

docker-clean: ## Очистить Docker ресурсы
	@echo "$(BLUE)Очистка Docker ресурсов...$(NC)"
	@docker-compose down -v
	@docker system prune -f


# Профилирование
profile-cpu: ## Создать CPU профиль
	@echo "$(BLUE)Создание CPU профиля...$(NC)"
	@go test -cpuprofile=cpu.prof -bench=. ./internal/services/
	@go tool pprof cpu.prof

profile-mem: ## Создать Memory профиль
	@echo "$(BLUE)Создание Memory профиля...$(NC)"
	@go test -memprofile=mem.prof -bench=. ./internal/services/
	@go tool pprof mem.prof

# Генерация кода
generate: generate-proto ## Генерировать код
	@echo "$(BLUE)Генерация кода...$(NC)"
	@go generate ./...

generate-proto: ## Генерировать gRPC код из proto файлов
	@echo "$(BLUE)Генерация gRPC кода из proto файлов...$(NC)"
	@mkdir -p proto/common proto/producer proto/consumer proto/monitor
	@export PATH=$$PATH:$$(go env GOPATH)/bin && \
	protoc --proto_path=proto --go_out=. --go_opt=module=pet-proj \
		--go-grpc_out=. --go-grpc_opt=module=pet-proj \
		proto/common.proto proto/producer.proto proto/consumer.proto proto/monitor.proto
	@echo "$(GREEN)gRPC код сгенерирован!$(NC)"


# Очистка
clean: ## Очистить артефакты сборки
	@echo "$(BLUE)Очистка артефактов...$(NC)"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@rm -f *.prof
	@go clean

# Проверка
check: fmt lint test ## Запустить все проверки
	@echo "$(GREEN)Все проверки пройдены!$(NC)"

# Разработка
dev-setup: deps infra-up ## Настроить среду разработки
	@echo "$(GREEN)Среда разработки настроена!$(NC)"

dev-clean: infra-down clean ## Очистить среду разработки
	@echo "$(GREEN)Среда разработки очищена!$(NC)"

# Мониторинг
monitor: ## Открыть мониторинг
	@echo "$(BLUE)Открытие мониторинга...$(NC)"
	@open http://localhost:9090  # Prometheus
	@open http://localhost:3000  # Grafana
	@open http://localhost:16686 # Jaeger

# Документация
docs: ## Генерировать документацию
	@echo "$(BLUE)Генерация документации...$(NC)"
	@godoc -http=:6060 &
	@echo "$(GREEN)Документация доступна на http://localhost:6060$(NC)"

# Информация о проекте
info: ## Показать информацию о проекте
	@echo "$(BLUE)Информация о проекте:$(NC)"
	@echo "  Название: $(APP_NAME)"
	@echo "  Версия: $(VERSION)"
	@echo "  Время сборки: $(BUILD_TIME)"
	@echo "  Go версия: $(GO_VERSION)"
	@echo "  Архитектура: $(shell go env GOARCH)"
	@echo "  ОС: $(shell go env GOOS)"

# Установка инструментов
install-tools: ## Установить необходимые инструменты
	@echo "$(BLUE)Установка инструментов...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/golang/mock/mockgen@latest
	@go install go.k6.io/k6@latest
	@echo "$(GREEN)Инструменты установлены!$(NC)"
