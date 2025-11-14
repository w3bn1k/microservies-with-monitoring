#!/bin/bash

# Универсальный скрипт запуска системы
# Поддерживает Docker Compose и Kubernetes

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функция показа справки
show_help() {
    echo -e "${BLUE}🚀 Универсальный запуск микросервисной системы${NC}"
    echo ""
    echo "Использование: $0 [ОПЦИЯ]"
    echo ""
    echo "ОПЦИИ:"
    echo "  docker     Запустить в Docker Compose (по умолчанию)"
    echo "  help       Показать эту справку"
    echo ""
    echo "Примеры:"
    echo "  $0            # Запустить в Docker Compose"
    echo "  $0 docker     # Запустить в Docker Compose"
}

# Функция запуска Docker Compose
start_docker() {
    echo -e "${BLUE}🐳 Запуск системы в Docker Compose...${NC}"
    
    # Проверяем наличие docker-compose
    if ! command -v docker-compose &> /dev/null; then
        echo -e "${RED}❌ docker-compose не найден. Установите Docker Compose.${NC}"
        exit 1
    fi
    
    # Проверяем, что Docker запущен
    if ! docker info &> /dev/null; then
        echo -e "${RED}❌ Docker не запущен. Запустите Docker Desktop.${NC}"
        exit 1
    fi
    
    # Генерируем gRPC код из proto файлов
    echo -e "${BLUE}📦 Генерация gRPC кода из proto файлов...${NC}"
    if ! command -v protoc &> /dev/null; then
        echo -e "${YELLOW}⚠️  protoc не найден. Пропускаем генерацию proto.${NC}"
    else
        cd "$(dirname "$0")/.." || exit 1
        export PATH=$PATH:$(go env GOPATH)/bin
        mkdir -p proto/common proto/producer proto/consumer proto/monitor
        protoc --proto_path=proto --go_out=. --go_opt=module=pet-proj \
            --go-grpc_out=. --go-grpc_opt=module=pet-proj \
            proto/common.proto proto/producer.proto proto/consumer.proto proto/monitor.proto
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✅ gRPC код сгенерирован!${NC}"
        else
            echo -e "${YELLOW}⚠️  Ошибка генерации proto, продолжаем...${NC}"
        fi
    fi
    
    # Запускаем систему
    docker-compose up -d
    
    echo -e "${GREEN}✅ Система запущена в Docker Compose!${NC}"
    echo -e "${BLUE}📊 Доступные сервисы:${NC}"
    echo -e "  ${GREEN}Producer API:${NC}     http://localhost:8080"
    echo -e "  ${GREEN}Consumer API:${NC}     http://localhost:8081"
    echo -e "  ${GREEN}Monitor API:${NC}      http://localhost:8082"
    echo -e "  ${GREEN}Kafka UI:${NC}         http://localhost:8083"
    echo -e "  ${GREEN}Redis Commander:${NC}  http://localhost:8084"
    echo -e "  ${GREEN}Prometheus:${NC}       http://localhost:9090"
    echo -e "  ${GREEN}Grafana:${NC}          http://localhost:3000"
    echo -e "  ${GREEN}Jaeger:${NC}           http://localhost:16686"
    echo -e "  ${GREEN}pgAdmin:${NC}          http://localhost:8085"
    
    echo ""
    echo -e "${YELLOW}💡 Для остановки: make stop${NC}"
    echo -e "${YELLOW}💡 Для просмотра логов: make logs${NC}"
}


# Основная логика
case "${1:-docker}" in
    "docker")
        start_docker
        ;;
    "help"|"-h"|"--help")
        show_help
        ;;
    *)
        echo -e "${RED}❌ Неизвестная опция: $1${NC}"
        echo ""
        show_help
        exit 1
        ;;
esac
