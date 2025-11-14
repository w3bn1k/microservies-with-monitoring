package grpc

import (
	"context"
	"fmt"
	"time"

	"pet-proj/internal/services"
	"pet-proj/proto/common"
	"pet-proj/proto/consumer"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ConsumerHandler обрабатывает gRPC запросы для Consumer Service
type ConsumerHandler struct {
	consumer.UnimplementedConsumerServiceServer
	consumerService *services.ConsumerService
	logger          *logrus.Logger
}

// NewConsumerHandler создает новый handler
func NewConsumerHandler(consumerService *services.ConsumerService, logger *logrus.Logger) *ConsumerHandler {
	return &ConsumerHandler{
		consumerService: consumerService,
		logger:          logger,
	}
}

// GetProcessedEvent получает обработанное событие по ID
func (h *ConsumerHandler) GetProcessedEvent(ctx context.Context, req *consumer.GetProcessedEventRequest) (*consumer.GetProcessedEventResponse, error) {
	if req == nil || req.EventId == "" {
		return nil, status.Error(codes.InvalidArgument, "event_id is required")
	}

	processedEvent, err := h.consumerService.GetProcessedEvent(ctx, req.EventId)
	if err != nil {
		h.logger.WithError(err).WithField("event_id", req.EventId).Error("Failed to get processed event")
		return nil, status.Error(codes.NotFound, "processed event not found")
	}

	// Извлекаем данные из processedEvent
	var event *common.Event
	var processedAt string
	var statusStr string

	if eventData, ok := processedEvent["event"].(map[string]interface{}); ok {
		event = mapToProtoEvent(eventData)
	}

	if pa, ok := processedEvent["processed_at"].(string); ok {
		processedAt = pa
	} else if pa, ok := processedEvent["processed_at"].(time.Time); ok {
		processedAt = pa.Format(time.RFC3339)
	}

	if s, ok := processedEvent["status"].(string); ok {
		statusStr = s
	}

	return &consumer.GetProcessedEventResponse{
		Success:     true,
		Event:       event,
		ProcessedAt: processedAt,
		Status:      statusStr,
		Message:     "Processed event retrieved successfully",
	}, nil
}

// GetStats возвращает статистику транзакций
func (h *ConsumerHandler) GetStats(ctx context.Context, req *consumer.GetStatsRequest) (*consumer.GetStatsResponse, error) {
	stats, err := h.consumerService.GetStats(ctx)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get stats")
		return nil, status.Error(codes.Internal, "failed to get stats")
	}

	protoStats := &common.Stats{
		Service:   "consumer",
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Metrics:   make(map[string]string),
	}

	// Конвертируем stats в map[string]string
	if stats != nil {
		for k, v := range stats {
			protoStats.Metrics[k] = fmt.Sprintf("%v", v)
		}
	}

	return &consumer.GetStatsResponse{
		Success: true,
		Stats:   protoStats,
	}, nil
}

// HealthCheck проверяет состояние сервиса
func (h *ConsumerHandler) HealthCheck(ctx context.Context, req *consumer.HealthCheckRequest) (*consumer.HealthCheckResponse, error) {
	health := &common.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   "consumer",
		Services:  make(map[string]string),
	}

	return &consumer.HealthCheckResponse{
		Status: health,
	}, nil
}

// mapToProtoEvent конвертирует map[string]interface{} в proto Event
func mapToProtoEvent(data map[string]interface{}) *common.Event {
	event := &common.Event{
		Data: make(map[string]string),
	}

	if id, ok := data["id"].(string); ok {
		event.Id = id
	}
	if t, ok := data["type"].(string); ok {
		event.Type = t
	}
	if uid, ok := data["user_id"].(string); ok {
		event.UserId = uid
	}
	if source, ok := data["source"].(string); ok {
		event.Source = source
	}
	if ts, ok := data["timestamp"].(string); ok {
		event.Timestamp = ts
	} else if ts, ok := data["timestamp"].(time.Time); ok {
		event.Timestamp = ts.Format(time.RFC3339)
	}

	// Конвертируем data map
	if dataMap, ok := data["data"].(map[string]interface{}); ok {
		for k, v := range dataMap {
			event.Data[k] = fmt.Sprintf("%v", v)
		}
	}

	return event
}

