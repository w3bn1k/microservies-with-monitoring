package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nikitakolesnik/pet-proj/internal/models"
	"github.com/nikitakolesnik/pet-proj/pkg/kafka"
	"github.com/nikitakolesnik/pet-proj/pkg/monitoring"
	"github.com/nikitakolesnik/pet-proj/pkg/redis"
	"github.com/sirupsen/logrus"
)

// обрабатывает события: отправляет в Kafka и кэширует в Redis
type EventService struct {
	kafkaProducer kafka.ProducerInterface
	redisClient   redis.ClientInterface
	logger        *logrus.Logger
}

func NewEventService(kafkaProducer kafka.ProducerInterface, redisClient redis.ClientInterface, logger *logrus.Logger) *EventService {
	return &EventService{
		kafkaProducer: kafkaProducer,
		redisClient:   redisClient,
		logger:        logger,
	}
}

// отправляет событие в Kafka и кэширует в Redis
func (s *EventService) SendEvent(ctx context.Context, event *models.Event) error {
	start := time.Now()

	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Отправляем событие в Kafka
	kafkaStatus := models.StatusOK
	if err := s.kafkaProducer.SendMessage(ctx, event.ID, event); err != nil {
		kafkaStatus = models.StatusBad
		s.logger.WithError(err).Error("Failed to send event to Kafka")
		monitoring.KafkaMessagesTotal.WithLabelValues("user-events", "failed").Inc()
		return err
	}
	monitoring.KafkaMessagesTotal.WithLabelValues("user-events", "success").Inc()

	// Кэшируем событие в Redis на 10 минут
	redisStatus := models.StatusOK
	cacheKey := fmt.Sprintf("event:%s", event.ID)
	if err := s.redisClient.Set(ctx, cacheKey, event, 10*time.Minute); err != nil {
		redisStatus = models.StatusBad
		s.logger.WithError(err).Error("Failed to cache event in Redis")
	}
	monitoring.RedisOperationsTotal.WithLabelValues("set", redisStatus).Inc()

	// Записываем метрики производительности
	duration := time.Since(start).Milliseconds()
	monitoring.TransactionsTotal.WithLabelValues(models.ServiceProducer, kafkaStatus, redisStatus).Inc()
	monitoring.EventsProcessedTotal.WithLabelValues(event.Type, models.ServiceProducer, "success").Inc()

	s.logger.WithFields(logrus.Fields{
		"event_id":     event.ID,
		"event_type":   event.Type,
		"kafka_status": kafkaStatus,
		"redis_status": redisStatus,
		"duration_ms":  duration,
	}).Info("Event processed successfully")

	return nil
}

// получает событие из кэша Redis по ID
func (s *EventService) GetEvent(ctx context.Context, eventID string) (*models.Event, error) {
	cacheKey := fmt.Sprintf("event:%s", eventID)

	var event models.Event
	if err := s.redisClient.Get(ctx, cacheKey, &event); err != nil {
		s.logger.WithError(err).WithField("event_id", eventID).Error("Failed to get event from cache")
		return nil, err
	}

	return &event, nil
}
