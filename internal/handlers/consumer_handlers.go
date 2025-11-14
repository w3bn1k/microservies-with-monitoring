package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"pet-proj/internal/services"
	"pet-proj/pkg/monitoring"
	"github.com/sirupsen/logrus"
)

type ConsumerHandlers struct {
	consumerService *services.ConsumerService
	logger          *logrus.Logger
}

func NewConsumerHandlers(consumerService *services.ConsumerService, logger *logrus.Logger) *ConsumerHandlers {
	return &ConsumerHandlers{
		consumerService: consumerService,
		logger:          logger,
	}
}

func (h *ConsumerHandlers) GetProcessedEvent(c *gin.Context) {
	start := time.Now()
	eventID := c.Param("id")

	processedEvent, err := h.consumerService.GetProcessedEvent(c.Request.Context(), eventID)
	if err != nil {
		h.logger.WithError(err).WithField("event_id", eventID).Error("Failed to get processed event")
		c.JSON(http.StatusNotFound, gin.H{"error": "Processed event not found"})
		monitoring.HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/events/:id", "404").Inc()
		return
	}

	duration := time.Since(start).Seconds()
	monitoring.HTTPRequestDuration.WithLabelValues("GET", "/api/v1/events/:id").Observe(duration)
	monitoring.HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/events/:id", "200").Inc()

	c.JSON(http.StatusOK, processedEvent)
}

func (h *ConsumerHandlers) GetStats(c *gin.Context) {
	start := time.Now()

	stats, err := h.consumerService.GetStats(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get consumer stats")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		monitoring.HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/stats", "500").Inc()
		return
	}

	duration := time.Since(start).Seconds()
	monitoring.HTTPRequestDuration.WithLabelValues("GET", "/api/v1/stats").Observe(duration)
	monitoring.HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/stats", "200").Inc()

	c.JSON(http.StatusOK, stats)
}

func (h *ConsumerHandlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "consumer",
	})
}
