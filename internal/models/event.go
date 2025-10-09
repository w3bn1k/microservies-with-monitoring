package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	UserID    string                 `json:"user_id"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
}

func NewEvent(eventType, userID, source string, data map[string]interface{}) *Event {
	return &Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		UserID:    userID,
		Data:      data,
		Timestamp: time.Now(),
		Source:    source,
	}
}

const (
	EventTypeUserAction    = "user_action"
	EventTypeSystemMetric  = "system_metric"
	EventTypeBusinessEvent = "business_event"
	EventTypeError         = "error"
)

const (
	SourceProducer = "producer"
	SourceConsumer = "consumer"
	SourceMonitor  = "monitor"
	SourceGateway  = "api-gateway"
)
