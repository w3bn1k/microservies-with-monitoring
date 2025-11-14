package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"pet-proj/internal/config"
	grpchandler "pet-proj/internal/grpc"
	"pet-proj/internal/services"
	"pet-proj/pkg/grpc"
	"pet-proj/pkg/kafka"
	"pet-proj/pkg/redis"
	"pet-proj/proto/producer"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	// Загружаем конфигурацию из переменных окружения
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name:     getEnv("SERVICE_NAME", "producer"),
			Port:     getEnvAsInt("SERVICE_PORT", 8080),
			GRPCPort: getEnvAsInt("GRPC_PORT", 9090),
		},
		Kafka: config.KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKERS", "kafka:29092")},
			Topic:   getEnv("KAFKA_TOPIC", "user-events"),
			GroupID: getEnv("KAFKA_GROUP_ID", "consumer-group"),
		},
		Redis: config.RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "redis:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			Timeout:  getEnvAsDuration("REDIS_TIMEOUT", "5s"),
		},
		Postgres: config.PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "postgres"),
			Port:     getEnvAsInt("POSTGRES_PORT", 5432),
			Database: getEnv("POSTGRES_DB", "microservices"),
			Username: getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "password"),
			SSLMode:  getEnv("POSTGRES_SSL_MODE", "disable"),
		},
		Monitoring: config.MonitoringConfig{
			PrometheusPort: getEnvAsInt("PROMETHEUS_PORT", 9090),
			JaegerEndpoint: getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
		},
	}

	// Настраиваем логгер
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.WithFields(logrus.Fields{
		"version":    version,
		"build_time": buildTime,
		"service":    "producer",
	}).Info("Starting Producer Service")

	// Инициализируем Kafka producer
	kafkaProducer, err := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic, logrus.StandardLogger())
	if err != nil {
		logrus.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Инициализируем Redis клиент
	redisClient := redis.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB, logrus.StandardLogger())
	defer redisClient.Close()

	// Проверяем подключение к Redis
	ctx := context.Background()
	if err := redisClient.Ping(ctx); err != nil {
		logrus.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Создаем сервисы
	eventService := services.NewEventService(kafkaProducer, redisClient, logrus.StandardLogger())

	// Настраиваем gRPC сервер
	grpcConfig := grpc.DefaultServerConfig(cfg.Service.GRPCPort, logrus.StandardLogger())
	grpcServer := grpc.NewServer(grpcConfig)

	// Регистрируем gRPC handlers
	producerHandler := grpchandler.NewProducerHandler(eventService, logrus.StandardLogger())
	producer.RegisterProducerServiceServer(grpcServer.GetServer(), producerHandler)

	// Запускаем gRPC сервер
	if err := grpcServer.Start(); err != nil {
		logrus.Fatalf("Failed to start gRPC server: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		grpcServer.Stop(ctx)
	}()

	logrus.Infof("Producer Service started on gRPC port %d", cfg.Service.GRPCPort)

	// Ждем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// Graceful shutdown с таймаутом 30 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := grpcServer.Stop(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	logrus.Info("Server exited")
}

// получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// получает переменную окружения как int или возвращает значение по умолчанию
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// получает переменную окружения как time.Duration или возвращает значение по умолчанию
func getEnvAsDuration(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

