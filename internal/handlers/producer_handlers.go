package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nikitakolesnik/pet-proj/internal/models"
	"github.com/nikitakolesnik/pet-proj/internal/services"
	"github.com/nikitakolesnik/pet-proj/pkg/monitoring"
	"github.com/sirupsen/logrus"
)

// обрабатывает HTTP запросы для Producer Service
type ProducerHandlers struct {
	eventService *services.EventService
	logger       *logrus.Logger
}

func NewProducerHandlers(eventService *services.EventService, logger *logrus.Logger) *ProducerHandlers {
	return &ProducerHandlers{
		eventService: eventService,
		logger:       logger,
	}
}

// принимает событие через HTTP POST и отправляет его в систему
func (h *ProducerHandlers) SendEvent(c *gin.Context) {
	start := time.Now()

	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		h.logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		monitoring.HTTPRequestsTotal.WithLabelValues("POST", "/api/v1/events", "400").Inc()
		return
	}

	// Устанавливаем источник события, если не указан
	if event.Source == "" {
		event.Source = models.SourceProducer
	}

	if err := h.eventService.SendEvent(c.Request.Context(), &event); err != nil {
		h.logger.WithError(err).Error("Failed to send event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send event"})
		monitoring.HTTPRequestsTotal.WithLabelValues("POST", "/api/v1/events", "500").Inc()
		return
	}

	duration := time.Since(start).Seconds()
	monitoring.HTTPRequestDuration.WithLabelValues("POST", "/api/v1/events").Observe(duration)
	monitoring.HTTPRequestsTotal.WithLabelValues("POST", "/api/v1/events", "200").Inc()

	c.JSON(http.StatusOK, gin.H{
		"message":     "Event sent successfully",
		"event_id":    event.ID,
		"duration_ms": int64(duration * 1000),
	})
}

// возвращает событие по ID из кэша Redis
func (h *ProducerHandlers) GetEvent(c *gin.Context) {
	start := time.Now()
	eventID := c.Param("id")

	event, err := h.eventService.GetEvent(c.Request.Context(), eventID)
	if err != nil {
		h.logger.WithError(err).WithField("event_id", eventID).Error("Failed to get event")
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		monitoring.HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/events/:id", "404").Inc()
		return
	}

	duration := time.Since(start).Seconds()
	monitoring.HTTPRequestDuration.WithLabelValues("GET", "/api/v1/events/:id").Observe(duration)
	monitoring.HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/events/:id", "200").Inc()

	c.JSON(http.StatusOK, event)
}

// возвращает базовую статистику сервиса
func (h *ProducerHandlers) GetStats(c *gin.Context) {
	start := time.Now()

	stats := gin.H{
		"service":   "producer",
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(start).String(),
	}

	duration := time.Since(start).Seconds()
	monitoring.HTTPRequestDuration.WithLabelValues("GET", "/api/v1/stats").Observe(duration)
	monitoring.HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/stats", "200").Inc()

	c.JSON(http.StatusOK, stats)
}

// проверяет состояние сервиса
func (h *ProducerHandlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "producer",
	})
}
