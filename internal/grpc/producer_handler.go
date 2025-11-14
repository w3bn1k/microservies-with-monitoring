package grpc

import (
	"context"
	"fmt"
	"time"

	"pet-proj/internal/models"
	"pet-proj/internal/services"
	"pet-proj/proto/common"
	"pet-proj/proto/producer"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProducerHandler обрабатывает gRPC запросы для Producer Service
type ProducerHandler struct {
	producer.UnimplementedProducerServiceServer
	eventService *services.EventService
	logger       *logrus.Logger
}

// NewProducerHandler создает новый handler
func NewProducerHandler(eventService *services.EventService, logger *logrus.Logger) *ProducerHandler {
	return &ProducerHandler{
		eventService: eventService,
		logger:       logger,
	}
}

// SendEvent отправляет событие в Kafka и кэширует в Redis
func (h *ProducerHandler) SendEvent(ctx context.Context, req *producer.SendEventRequest) (*producer.SendEventResponse, error) {
	if req == nil || req.Event == nil {
		return nil, status.Error(codes.InvalidArgument, "event is required")
	}

	start := time.Now()

	// Конвертируем proto Event в models.Event
	event := protoToEvent(req.Event)

	// Отправляем событие
	if err := h.eventService.SendEvent(ctx, event); err != nil {
		h.logger.WithError(err).Error("Failed to send event via gRPC")
		return nil, status.Error(codes.Internal, "failed to send event")
	}

	duration := time.Since(start).Milliseconds()

	return &producer.SendEventResponse{
		Success:    true,
		EventId:    event.ID,
		DurationMs: duration,
		Message:    "Event sent successfully",
	}, nil
}

// GetEvent получает событие по ID из кэша
func (h *ProducerHandler) GetEvent(ctx context.Context, req *producer.GetEventRequest) (*producer.GetEventResponse, error) {
	if req == nil || req.EventId == "" {
		return nil, status.Error(codes.InvalidArgument, "event_id is required")
	}

	event, err := h.eventService.GetEvent(ctx, req.EventId)
	if err != nil {
		h.logger.WithError(err).WithField("event_id", req.EventId).Error("Failed to get event")
		return nil, status.Error(codes.NotFound, "event not found")
	}

	return &producer.GetEventResponse{
		Success: true,
		Event:   eventToProto(event),
		Message: "Event retrieved successfully",
	}, nil
}

// GetStats возвращает статистику сервиса
func (h *ProducerHandler) GetStats(ctx context.Context, req *producer.GetStatsRequest) (*producer.GetStatsResponse, error) {
	stats := &common.Stats{
		Service:   "producer",
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Metrics:   make(map[string]string),
	}

	return &producer.GetStatsResponse{
		Success: true,
		Stats:   stats,
	}, nil
}

// HealthCheck проверяет состояние сервиса
func (h *ProducerHandler) HealthCheck(ctx context.Context, req *producer.HealthCheckRequest) (*producer.HealthCheckResponse, error) {
	health := &common.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   "producer",
		Services:  make(map[string]string),
	}

	return &producer.HealthCheckResponse{
		Status: health,
	}, nil
}

// protoToEvent конвертирует proto Event в models.Event
func protoToEvent(protoEvent *common.Event) *models.Event {
	event := &models.Event{
		ID:        protoEvent.Id,
		Type:      protoEvent.Type,
		UserID:    protoEvent.UserId,
		Data:      make(map[string]interface{}),
		Source:    protoEvent.Source,
	}

	// Конвертируем map[string]string в map[string]interface{}
	for k, v := range protoEvent.Data {
		event.Data[k] = v
	}

	// Парсим timestamp
	if protoEvent.Timestamp != "" {
		if t, err := time.Parse(time.RFC3339, protoEvent.Timestamp); err == nil {
			event.Timestamp = t
		} else {
			event.Timestamp = time.Now()
		}
	} else {
		event.Timestamp = time.Now()
	}

	return event
}

// eventToProto конвертирует models.Event в proto Event
func eventToProto(event *models.Event) *common.Event {
	protoEvent := &common.Event{
		Id:        event.ID,
		Type:      event.Type,
		UserId:    event.UserID,
		Data:      make(map[string]string),
		Timestamp: event.Timestamp.Format(time.RFC3339),
		Source:    event.Source,
	}

	// Конвертируем map[string]interface{} в map[string]string
	for k, v := range event.Data {
		if str, ok := v.(string); ok {
			protoEvent.Data[k] = str
		} else {
			protoEvent.Data[k] = fmt.Sprintf("%v", v)
		}
	}

	return protoEvent
}

