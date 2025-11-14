package grpc

import (
	"context"
	"fmt"
	"time"

	"pet-proj/internal/models"
	"pet-proj/internal/services"
	"pet-proj/proto/common"
	"pet-proj/proto/monitor"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MonitorHandler обрабатывает gRPC запросы для Monitor Service
type MonitorHandler struct {
	monitor.UnimplementedMonitorServiceServer
	monitorService *services.MonitorService
	logger         *logrus.Logger
}

// NewMonitorHandler создает новый handler
func NewMonitorHandler(monitorService *services.MonitorService, logger *logrus.Logger) *MonitorHandler {
	return &MonitorHandler{
		monitorService: monitorService,
		logger:         logger,
	}
}

// GetHealth возвращает общее состояние системы
func (h *MonitorHandler) GetHealth(ctx context.Context, req *monitor.GetHealthRequest) (*monitor.GetHealthResponse, error) {
	health := h.monitorService.GetSystemHealth(ctx)

	protoHealth := &common.HealthStatus{
		Status:    health["overall_status"].(string),
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   "monitor",
		Services:  make(map[string]string),
	}

	if services, ok := health["services"].(map[string]string); ok {
		protoHealth.Services = services
	}

	return &monitor.GetHealthResponse{
		Success: true,
		Health:  protoHealth,
	}, nil
}

// GetTransactions возвращает список транзакций
func (h *MonitorHandler) GetTransactions(ctx context.Context, req *monitor.GetTransactionsRequest) (*monitor.GetTransactionsResponse, error) {
	limit := 100 // default
	if req != nil && req.Limit > 0 {
		limit = int(req.Limit)
	}

	transactions, err := h.monitorService.GetTransactions(ctx, limit)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get transactions")
		return nil, status.Error(codes.Internal, "failed to get transactions")
	}

	protoTransactions := make([]*common.Transaction, 0, len(transactions))
	for _, tx := range transactions {
		protoTransactions = append(protoTransactions, transactionToProto(tx))
	}

	return &monitor.GetTransactionsResponse{
		Success:      true,
		Transactions: protoTransactions,
		Message:      fmt.Sprintf("Retrieved %d transactions", len(protoTransactions)),
	}, nil
}

// GetStats возвращает статистику транзакций
func (h *MonitorHandler) GetStats(ctx context.Context, req *monitor.GetStatsRequest) (*monitor.GetStatsResponse, error) {
	stats, err := h.monitorService.GetTransactionStats(ctx)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get stats")
		return nil, status.Error(codes.Internal, "failed to get stats")
	}

	protoStats := &common.Stats{
		Service:   "monitor",
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Metrics:   make(map[string]string),
	}

	if stats != nil {
		for k, v := range stats {
			protoStats.Metrics[k] = fmt.Sprintf("%v", v)
		}
	}

	return &monitor.GetStatsResponse{
		Success: true,
		Stats:   protoStats,
	}, nil
}

// GetDashboard возвращает данные для дашборда
func (h *MonitorHandler) GetDashboard(ctx context.Context, req *monitor.GetDashboardRequest) (*monitor.GetDashboardResponse, error) {
	// Получаем health
	health := h.monitorService.GetSystemHealth(ctx)
	protoHealth := &common.HealthStatus{
		Status:    health["overall_status"].(string),
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   "monitor",
		Services:  make(map[string]string),
	}
	if services, ok := health["services"].(map[string]string); ok {
		protoHealth.Services = services
	}

	// Получаем последние транзакции
	transactions, err := h.monitorService.GetTransactions(ctx, 10)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get transactions for dashboard")
	}

	protoTransactions := make([]*common.Transaction, 0, len(transactions))
	for _, tx := range transactions {
		protoTransactions = append(protoTransactions, transactionToProto(tx))
	}

	// Получаем статистику
	stats, err := h.monitorService.GetTransactionStats(ctx)
	protoStats := &common.Stats{
		Service:   "monitor",
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Metrics:   make(map[string]string),
	}
	if stats != nil {
		for k, v := range stats {
			protoStats.Metrics[k] = fmt.Sprintf("%v", v)
		}
	}

	return &monitor.GetDashboardResponse{
		Success:            true,
		Health:            protoHealth,
		RecentTransactions: protoTransactions,
		Stats:             protoStats,
		Message:           "Dashboard data retrieved successfully",
	}, nil
}

// HealthCheck проверяет состояние сервиса
func (h *MonitorHandler) HealthCheck(ctx context.Context, req *monitor.HealthCheckRequest) (*monitor.HealthCheckResponse, error) {
	health := &common.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   "monitor",
		Services:  make(map[string]string),
	}

	return &monitor.HealthCheckResponse{
		Status: health,
	}, nil
}

// transactionToProto конвертирует models.Transaction в proto Transaction
func transactionToProto(tx *models.Transaction) *common.Transaction {
	return &common.Transaction{
		Id:          tx.ID,
		Timestamp:   tx.Timestamp.Format(time.RFC3339),
		KafkaStatus: tx.KafkaStatus,
		RedisStatus: tx.RedisStatus,
		DurationMs:  tx.Duration,
		Service:     tx.Service,
		EventId:     tx.EventID,
		ErrorMsg:    tx.ErrorMsg,
		CreatedAt:   tx.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   tx.UpdatedAt.Format(time.RFC3339),
	}
}

