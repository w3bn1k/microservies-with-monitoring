package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nikitakolesnik/pet-proj/internal/models"
	"github.com/nikitakolesnik/pet-proj/pkg/kafka"
	"github.com/nikitakolesnik/pet-proj/pkg/monitoring"
	"github.com/nikitakolesnik/pet-proj/pkg/postgres"
	"github.com/nikitakolesnik/pet-proj/pkg/redis"
	"github.com/sirupsen/logrus"
)

// мониторит состояние системы и записывает метрики каждую минуту
type MonitorService struct {
	postgresClient postgres.ClientInterface
	redisClient    redis.ClientInterface
	kafkaHealth    *kafka.HealthChecker
	postgresHealth *postgres.HealthChecker
	logger         *logrus.Logger
}

func NewMonitorService(postgresClient postgres.ClientInterface, redisClient redis.ClientInterface, kafkaBrokers []string, postgresDSN string, logger *logrus.Logger) (*MonitorService, error) {
	// Создаем Kafka health checker
	kafkaHealth, err := kafka.NewHealthChecker(kafkaBrokers, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka health checker: %w", err)
	}

	// Создаем бд health checker
	postgresHealth, err := postgres.NewHealthChecker(postgresDSN, logger)
	if err != nil {
		kafkaHealth.Close()
		return nil, fmt.Errorf("failed to create PostgreSQL health checker: %w", err)
	}

	return &MonitorService{
		postgresClient: postgresClient,
		redisClient:    redisClient,
		kafkaHealth:    kafkaHealth,
		postgresHealth: postgresHealth,
		logger:         logger,
	}, nil
}

// запускает мониторинг системы каждую минуту
func (s *MonitorService) StartMonitoring(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	s.logger.Info("Starting monitoring service")

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Stopping monitoring service")
			return
		case <-ticker.C:
			s.recordSystemMetrics(ctx)
		}
	}
}

// собирает метрики системы и записывает транзакцию
func (s *MonitorService) recordSystemMetrics(ctx context.Context) {
	start := time.Now()

	// Проверяем состояние всех компонентов системы
	kafkaStatus := s.checkKafkaHealth(ctx)
	redisStatus := s.checkRedisHealth(ctx)
	postgresStatus := s.checkPostgresHealth(ctx)

	// Создаем транзакцию с результатами мониторинга
	transaction := &models.Transaction{
		Timestamp:   time.Now(),
		KafkaStatus: kafkaStatus,
		RedisStatus: redisStatus,
		Duration:    time.Since(start).Milliseconds(),
		Service:     models.ServiceMonitor,
		EventID:     "monitor-" + time.Now().Format("20060102150405"),
		ErrorMsg:    "",
	}

	// Сохраняем транзакцию в бд
	if err := s.postgresClient.InsertTransaction(transaction); err != nil {
		s.logger.WithError(err).Error("Failed to save monitor transaction")
	}

	// Обновляем метрики Prometheus
	monitoring.TransactionsTotal.WithLabelValues(models.ServiceMonitor, kafkaStatus, redisStatus).Inc()

	s.logger.WithFields(logrus.Fields{
		"kafka_status":    kafkaStatus,
		"redis_status":    redisStatus,
		"postgres_status": postgresStatus,
		"duration_ms":     time.Since(start).Milliseconds(),
	}).Info("System metrics recorded")
}

// проверяет состояние Kafka
func (s *MonitorService) checkKafkaHealth(ctx context.Context) string {
	if err := s.kafkaHealth.CheckHealth(ctx); err != nil {
		s.logger.WithError(err).Error("Kafka health check failed")
		return models.StatusBad
	}
	return models.StatusOK
}

// проверяет состояние Redis через ping
func (s *MonitorService) checkRedisHealth(ctx context.Context) string {
	if err := s.redisClient.Ping(ctx); err != nil {
		s.logger.WithError(err).Error("Redis health check failed")
		return models.StatusBad
	}
	return models.StatusOK
}

// проверяет состояние бд
func (s *MonitorService) checkPostgresHealth(ctx context.Context) string {
	if err := s.postgresHealth.CheckHealth(ctx); err != nil {
		s.logger.WithError(err).Error("PostgreSQL health check failed")
		return models.StatusBad
	}
	return models.StatusOK
}

// возвращает общее состояние системы
func (s *MonitorService) GetSystemHealth(ctx context.Context) map[string]interface{} {
	kafkaStatus := s.checkKafkaHealth(ctx)
	redisStatus := s.checkRedisHealth(ctx)
	postgresStatus := s.checkPostgresHealth(ctx)

	overallStatus := "healthy"
	if kafkaStatus == models.StatusBad || redisStatus == models.StatusBad || postgresStatus == models.StatusBad {
		overallStatus = "unhealthy"
	}

	return map[string]interface{}{
		"overall_status": overallStatus,
		"services": map[string]string{
			"kafka":    kafkaStatus,
			"redis":    redisStatus,
			"postgres": postgresStatus,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}
}

// возвращает список транзакций из бд
func (s *MonitorService) GetTransactions(ctx context.Context, limit int) ([]*models.Transaction, error) {
	return s.postgresClient.GetTransactions(limit)
}

// возвращает статистику транзакций из бд
func (s *MonitorService) GetTransactionStats(ctx context.Context) (map[string]interface{}, error) {
	return s.postgresClient.GetTransactionStats()
}

// Close закрывает все соединения
func (s *MonitorService) Close() error {
	var err error
	if s.kafkaHealth != nil {
		if closeErr := s.kafkaHealth.Close(); closeErr != nil {
			err = closeErr
		}
	}
	if s.postgresHealth != nil {
		if closeErr := s.postgresHealth.Close(); closeErr != nil {
			err = closeErr
		}
	}
	return err
}
