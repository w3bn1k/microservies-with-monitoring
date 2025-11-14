package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"pet-proj/internal/models"
	"pet-proj/pkg/monitoring"
	"pet-proj/pkg/postgres"
	"pet-proj/pkg/redis"
	"github.com/sirupsen/logrus"
)

// обрабатывает сообщения из Kafka и сохраняет транзакции в бд
type ConsumerService struct {
	redisClient    redis.ClientInterface
	postgresClient postgres.ClientInterface
	logger         *logrus.Logger
}

func NewConsumerService(redisClient redis.ClientInterface, postgresClient postgres.ClientInterface, logger *logrus.Logger) *ConsumerService {
	return &ConsumerService{
		redisClient:    redisClient,
		postgresClient: postgresClient,
		logger:         logger,
	}
}

// обрабатывает сообщение из Kafka и создает транзакцию
func (s *ConsumerService) HandleMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	start := time.Now()

	var event models.Event
	if err := json.Unmarshal(message.Value, &event); err != nil {
		s.logger.WithError(err).Error("Failed to unmarshal message")
		monitoring.KafkaMessageDuration.WithLabelValues("user-events").Observe(time.Since(start).Seconds())
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"event_id":   event.ID,
		"event_type": event.Type,
		"partition":  message.Partition,
		"offset":     message.Offset,
	}).Info("Processing message")

	// Обрабатываем событие и получаем статусы операций
	kafkaStatus, redisStatus, err := s.ProcessEvent(ctx, &event)

	// Создаем транзакцию с результатами обработки
	transaction := &models.Transaction{
		Timestamp:   time.Now(),
		KafkaStatus: kafkaStatus,
		RedisStatus: redisStatus,
		Duration:    time.Since(start).Milliseconds(),
		Service:     models.ServiceConsumer,
		EventID:     event.ID,
		ErrorMsg:    "",
	}

	if err != nil {
		transaction.ErrorMsg = err.Error()
		s.logger.WithError(err).Error("Failed to process event")
	}

	// Сохраняем транзакцию в бд
	if err := s.postgresClient.InsertTransaction(transaction); err != nil {
		s.logger.WithError(err).Error("Failed to save transaction")
	}

	// Обновляем метрики производительности
	monitoring.KafkaMessageDuration.WithLabelValues("user-events").Observe(time.Since(start).Seconds())
	monitoring.EventsProcessedTotal.WithLabelValues(event.Type, models.ServiceConsumer, "success").Inc()
	monitoring.TransactionsTotal.WithLabelValues(models.ServiceConsumer, kafkaStatus, redisStatus).Inc()

	return nil
}

// обрабатывает событие и кэширует результат в Redis
func (s *ConsumerService) ProcessEvent(ctx context.Context, event *models.Event) (string, string, error) {
	// Кэшируем обработанное событие в Redis на 30 минут
	redisStatus := models.StatusOK
	cacheKey := fmt.Sprintf("processed_event:%s", event.ID)
	processedEvent := map[string]interface{}{
		"event":        event,
		"processed_at": time.Now(),
		"status":       "processed",
	}

	if err := s.redisClient.Set(ctx, cacheKey, processedEvent, 30*time.Minute); err != nil {
		redisStatus = models.StatusBad
		s.logger.WithError(err).Error("Failed to cache processed event")
	}

	// Выполняем бизнес-логику в зависимости от типа события
	switch event.Type {
	case models.EventTypeUserAction:
		s.processUserAction(event)
	case models.EventTypeSystemMetric:
		s.processSystemMetric(event)
	case models.EventTypeBusinessEvent:
		s.processBusinessEvent(event)
	default:
		s.logger.WithField("event_type", event.Type).Warn("Unknown event type")
	}

	return models.StatusOK, redisStatus, nil
}

// обрабатывает действия пользователя
func (s *ConsumerService) processUserAction(event *models.Event) {
	s.logger.WithFields(logrus.Fields{
		"event_id": event.ID,
		"user_id":  event.UserID,
		"action":   event.Data["action"],
	}).Info("Processing user action")
}

// обрабатывает системные метрики
func (s *ConsumerService) processSystemMetric(event *models.Event) {
	s.logger.WithFields(logrus.Fields{
		"event_id": event.ID,
		"cpu":      event.Data["cpu_usage"],
		"memory":   event.Data["memory_usage"],
	}).Info("Processing system metric")
}

// обрабатывает бизнес-события
func (s *ConsumerService) processBusinessEvent(event *models.Event) {
	s.logger.WithFields(logrus.Fields{
		"event_id": event.ID,
		"order_id": event.Data["order_id"],
		"amount":   event.Data["amount"],
		"currency": event.Data["currency"],
	}).Info("Processing business event")
}

// получает обработанное событие из кэша Redis
func (s *ConsumerService) GetProcessedEvent(ctx context.Context, eventID string) (map[string]interface{}, error) {
	cacheKey := fmt.Sprintf("processed_event:%s", eventID)

	var processedEvent map[string]interface{}
	if err := s.redisClient.Get(ctx, cacheKey, &processedEvent); err != nil {
		s.logger.WithError(err).WithField("event_id", eventID).Error("Failed to get processed event from cache")
		return nil, err
	}

	return processedEvent, nil
}

// возвращает статистику транзакций из бд
func (s *ConsumerService) GetStats(ctx context.Context) (map[string]interface{}, error) {
	stats, err := s.postgresClient.GetTransactionStats()
	if err != nil {
		s.logger.WithError(err).Error("Failed to get transaction stats")
		return nil, err
	}

	return stats, nil
}

// сохраняет транзакцию в бд
func (s *ConsumerService) InsertTransaction(tx *models.Transaction) error {
	return s.postgresClient.InsertTransaction(tx)
}
