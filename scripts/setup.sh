#!/bin/bash

# Setup script for microservices project

set -e

echo "🚀 Setting up microservices project..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ Go version $GO_VERSION is too old. Please install Go $REQUIRED_VERSION or later."
    exit 1
fi

echo "✅ Go version $GO_VERSION is compatible"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker."
    exit 1
fi

echo "✅ Docker is installed"

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose."
    exit 1
fi

echo "✅ Docker Compose is installed"

# Download Go dependencies
echo "📦 Downloading Go dependencies..."
go mod download

# Install development tools
echo "🛠️ Installing development tools..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/golang/mock/mockgen@latest
go install go.k6.io/k6@latest

echo "✅ Development tools installed"

# Create necessary directories
echo "📁 Creating directories..."
mkdir -p logs
mkdir -p data/postgres
mkdir -p data/redis
mkdir -p data/kafka

echo "✅ Directories created"

# Set up environment variables
echo "🔧 Setting up environment variables..."
if [ ! -f .env ]; then
    cat > .env << EOF
# Service configuration
SERVICE_NAME=microservices
SERVICE_PORT=8080
LOG_LEVEL=info

# Kafka configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=user-events
KAFKA_GROUP_ID=consumer-group

# Redis configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# PostgreSQL configuration
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=microservices
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password

# Monitoring configuration
PROMETHEUS_PORT=9090
JAEGER_ENDPOINT=http://localhost:14268/api/traces
EOF
    echo "✅ Environment file created"
else
    echo "✅ Environment file already exists"
fi

# Start infrastructure
echo "🐳 Starting infrastructure services..."
docker-compose -f docker-compose.infrastructure.yml up -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 30

# Check if services are running
echo "🔍 Checking service health..."

# Check Kafka
if docker-compose -f docker-compose.infrastructure.yml exec -T kafka kafka-broker-api-versions --bootstrap-server localhost:9092 > /dev/null 2>&1; then
    echo "✅ Kafka is ready"
else
    echo "❌ Kafka is not ready"
fi

# Check Redis
if docker-compose -f docker-compose.infrastructure.yml exec -T redis redis-cli ping > /dev/null 2>&1; then
    echo "✅ Redis is ready"
else
    echo "❌ Redis is not ready"
fi

# Check PostgreSQL
if docker-compose -f docker-compose.infrastructure.yml exec -T postgres pg_isready -U postgres > /dev/null 2>&1; then
    echo "✅ PostgreSQL is ready"
else
    echo "❌ PostgreSQL is not ready"
fi

echo ""
echo "🎉 Setup completed successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Run 'make build' to build all services"
echo "2. Run 'make run-producer' to start the producer service"
echo "3. Run 'make run-consumer' to start the consumer service"
echo "4. Run 'make run-monitor' to start the monitor service"
echo ""
echo "🌐 Access points:"
echo "- Producer Service: http://localhost:8080"
echo "- Consumer Service: http://localhost:8081"
echo "- Monitor Service: http://localhost:8082"
echo "- Prometheus: http://localhost:9090"
echo "- Grafana: http://localhost:3000 (admin/admin)"
echo "- Jaeger: http://localhost:16686"
echo "- Kafka UI: http://localhost:8080"
echo "- Redis Commander: http://localhost:8081"
echo "- pgAdmin: http://localhost:8082 (admin@admin.com/admin)"
echo ""
echo "📚 For more information, see README.md"
