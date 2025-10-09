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
	"github.com/nikitakolesnik/pet-proj/pkg/redis"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	// Загружаем конфигурацию из переменных окружения
	cfg := &config.Config{
		Service: config.ServiceConfig{
			Name: getEnv("SERVICE_NAME", "producer"),
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

	// Настраиваем HTTP сервер
	router := gin.New()

	// Добавляем middleware
	router.Use(middleware.Logger(logrus.StandardLogger()))
	router.Use(middleware.Recovery(logrus.StandardLogger()))
	router.Use(middleware.Metrics())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())

	// Настраиваем маршруты
	producerHandlers := handlers.NewProducerHandlers(eventService, logrus.StandardLogger())
	setupProducerRoutes(router, producerHandlers)

	// Запускаем HTTP сервер
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Service.Port),
		Handler: router,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		logrus.Infof("Starting producer service on port %d", cfg.Service.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ждем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// Graceful shutdown с таймаутом 30 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
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

// настраивает HTTP маршруты для Producer Service
func setupProducerRoutes(router *gin.Engine, handlers *handlers.ProducerHandlers) {
	api := router.Group("/api/v1")
	{
		api.POST("/events", handlers.SendEvent)
		api.GET("/events/:id", handlers.GetEvent)
		api.GET("/stats", handlers.GetStats)
	}

	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(monitoring.Handler()))

	// Go profiling endpoints
	router.GET("/debug/pprof/*path", gin.WrapH(http.DefaultServeMux))
}
