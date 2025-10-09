package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nikitakolesnik/pet-proj/internal/services"
	"github.com/nikitakolesnik/pet-proj/pkg/monitoring"
	"github.com/sirupsen/logrus"
)

type MonitorHandlers struct {
	monitorService *services.MonitorService
	logger         *logrus.Logger
}

func NewMonitorHandlers(monitorService *services.MonitorService, logger *logrus.Logger) *MonitorHandlers {
	return &MonitorHandlers{
		monitorService: monitorService,
		logger:         logger,
	}
}

func (h *MonitorHandlers) GetHealth(c *gin.Context) {
	start := time.Now()

	health := h.monitorService.GetSystemHealth(c.Request.Context())

	duration := time.Since(start).Seconds()
	monitoring.HTTPRequestDuration.WithLabelValues("GET", "/api/v1/health").Observe(duration)
	monitoring.HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/health", "200").Inc()

	c.JSON(http.StatusOK, health)
}

func (h *MonitorHandlers) GetTransactions(c *gin.Context) {
	start := time.Now()

	// Parse limit parameter
	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 1000 {
			limit = parsed
		}
	}

	transactions, err := h.monitorService.GetTransactions(c.Request.Context(), limit)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get transactions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions"})
		monitoring.HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/transactions", "500").Inc()
		return
	}

	duration := time.Since(start).Seconds()
	monitoring.HTTPRequestDuration.WithLabelValues("GET", "/api/v1/transactions").Observe(duration)
	monitoring.HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/transactions", "200").Inc()

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"count":        len(transactions),
		"limit":        limit,
	})
}

func (h *MonitorHandlers) GetStats(c *gin.Context) {
	start := time.Now()

	stats, err := h.monitorService.GetTransactionStats(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get transaction stats")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		monitoring.HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/stats", "500").Inc()
		return
	}

	duration := time.Since(start).Seconds()
	monitoring.HTTPRequestDuration.WithLabelValues("GET", "/api/v1/stats").Observe(duration)
	monitoring.HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/stats", "200").Inc()

	c.JSON(http.StatusOK, stats)
}

func (h *MonitorHandlers) GetDashboard(c *gin.Context) {
	start := time.Now()

	health := h.monitorService.GetSystemHealth(c.Request.Context())

	transactions, _ := h.monitorService.GetTransactions(c.Request.Context(), 10)

	stats, _ := h.monitorService.GetTransactionStats(c.Request.Context())

	dashboard := gin.H{
		"health":       health,
		"recent_transactions": transactions,
		"stats":        stats,
		"timestamp":    time.Now().Format(time.RFC3339),
	}

	duration := time.Since(start).Seconds()
	monitoring.HTTPRequestDuration.WithLabelValues("GET", "/api/v1/dashboard").Observe(duration)
	monitoring.HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/dashboard", "200").Inc()

	c.JSON(http.StatusOK, dashboard)
}

func (h *MonitorHandlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "monitor",
	})
}
