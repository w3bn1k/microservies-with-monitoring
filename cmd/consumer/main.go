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
	"pet-proj/pkg/postgres"
	"pet-proj/pkg/redis"
	"pet-proj/proto/consumer"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name:     getEnv("SERVICE_NAME", "consumer"),
			Port:     getEnvAsInt("SERVICE_PORT", 8080),
			GRPCPort: getEnvAsInt("GRPC_PORT", 9091),
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

	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.WithFields(logrus.Fields{
		"version":    version,
		"build_time": buildTime,
		"service":    "consumer",
	}).Info("Starting Consumer Service")

	postgresClient, err := postgres.NewClient(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		logrus.StandardLogger(),
	)
	if err != nil {
		logrus.Fatalf("Failed to create PostgreSQL client: %v", err)
	}
	defer postgresClient.Close()

	if err := postgresClient.CreateTransactionTable(); err != nil {
		logrus.Fatalf("Failed to create transaction table: %v", err)
	}

	redisClient := redis.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB, logrus.StandardLogger())
	defer redisClient.Close()

	ctx := context.Background()
	if err := redisClient.Ping(ctx); err != nil {
		logrus.Fatalf("Failed to connect to Redis: %v", err)
	}

	kafkaConsumer, err := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID, logrus.StandardLogger())
	if err != nil {
		logrus.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()

	consumerService := services.NewConsumerService(redisClient, postgresClient, logrus.StandardLogger())
	kafkaConsumer.SetHandler(consumerService)

	// Настраиваем gRPC сервер
	grpcConfig := grpc.DefaultServerConfig(cfg.Service.GRPCPort, logrus.StandardLogger())
	grpcServer := grpc.NewServer(grpcConfig)

	// Регистрируем gRPC handlers
	consumerHandler := grpchandler.NewConsumerHandler(consumerService, logrus.StandardLogger())
	consumer.RegisterConsumerServiceServer(grpcServer.GetServer(), consumerHandler)

	// Запускаем gRPC сервер
	if err := grpcServer.Start(); err != nil {
		logrus.Fatalf("Failed to start gRPC server: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		grpcServer.Stop(ctx)
	}()

	logrus.Infof("Consumer Service started on gRPC port %d", cfg.Service.GRPCPort)

	go func() {
		logrus.Info("Starting Kafka consumer")
		if err := kafkaConsumer.Start(ctx); err != nil {
			logrus.WithError(err).Error("Kafka consumer stopped")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := grpcServer.Stop(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	logrus.Info("Server exited")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}
