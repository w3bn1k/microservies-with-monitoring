package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/nikitakolesnik/pet-proj/internal/config"
	"github.com/nikitakolesnik/pet-proj/internal/handlers"
	"github.com/nikitakolesnik/pet-proj/internal/middleware"
	"github.com/nikitakolesnik/pet-proj/internal/services"
	"github.com/nikitakolesnik/pet-proj/pkg/kafka"
	"github.com/nikitakolesnik/pet-proj/pkg/monitoring"
	"github.com/nikitakolesnik/pet-proj/pkg/postgres"
	"github.com/nikitakolesnik/pet-proj/pkg/redis"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name: getEnv("SERVICE_NAME", "consumer"),
			Port: getEnvAsInt("SERVICE_PORT", 8080),
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

	router := gin.New()

	router.Use(middleware.Logger(logrus.StandardLogger()))
	router.Use(middleware.Recovery(logrus.StandardLogger()))
	router.Use(middleware.Metrics())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())

	consumerHandlers := handlers.NewConsumerHandlers(consumerService, logrus.StandardLogger())
	setupConsumerRoutes(router, consumerHandlers)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Service.Port),
		Handler: router,
	}

	go func() {
		logrus.Infof("Starting consumer service on port %d", cfg.Service.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

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

	if err := srv.Shutdown(ctx); err != nil {
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

func setupConsumerRoutes(router *gin.Engine, handlers *handlers.ConsumerHandlers) {
	api := router.Group("/api/v1")
	{
		api.GET("/events/:id", handlers.GetProcessedEvent)
		api.GET("/stats", handlers.GetStats)
	}

	router.GET("/health", handlers.HealthCheck)

	router.GET("/metrics", gin.WrapH(monitoring.Handler()))

	router.GET("/debug/pprof/*path", gin.WrapH(http.DefaultServeMux))
}
