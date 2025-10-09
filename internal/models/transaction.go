package models

import "time"

type Transaction struct {
	ID          int64     `json:"id" db:"id"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	KafkaStatus string    `json:"kafka_status" db:"kafka_status"`
	RedisStatus string    `json:"redis_status" db:"redis_status"`
	Duration    int64     `json:"duration_ms" db:"duration_ms"`
	Service     string    `json:"service" db:"service"`
	EventID     string    `json:"event_id" db:"event_id"`
	ErrorMsg    string    `json:"error_msg" db:"error_msg"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type TransactionStatus struct {
	KafkaStatus string
	RedisStatus string
	Duration    int64
	ErrorMsg    string
}

const (
	StatusOK  = "ok"
	StatusBad = "bad"
)

const (
	ServiceProducer = "producer"
	ServiceConsumer = "consumer"
	ServiceMonitor  = "monitor"
	ServiceGateway  = "api-gateway"
)
