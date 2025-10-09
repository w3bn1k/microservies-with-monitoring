package kafka

import (
	"context"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type HealthChecker struct {
	client  sarama.Client
	logger  *logrus.Logger
	timeout time.Duration
}

func NewHealthChecker(brokers []string, logger *logrus.Logger) (*HealthChecker, error) {
	config := sarama.NewConfig()
	config.Net.DialTimeout = 5 * time.Second
	config.Net.ReadTimeout = 5 * time.Second
	config.Net.WriteTimeout = 5 * time.Second
	config.Metadata.Retry.Max = 3
	config.Metadata.Retry.Backoff = 100 * time.Millisecond

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, err
	}

	return &HealthChecker{
		client:  client,
		logger:  logger,
		timeout: 5 * time.Second,
	}, nil
}

// проверяет состояние Kafka кластера
func (h *HealthChecker) CheckHealth(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	// Проверяем подключение к кластеру
	if err := h.client.RefreshMetadata(); err != nil {
		h.logger.WithError(err).Error("Failed to refresh Kafka metadata")
		return err
	}

	// Проверяем, что можем получить список топиков
	topics, err := h.client.Topics()
	if err != nil {
		h.logger.WithError(err).Error("Failed to get Kafka topics")
		return err
	}

	// Проверяем, что кластер не пустой
	if len(topics) == 0 {
		h.logger.Warn("Kafka cluster has no topics")
	}

	// Проверяем состояние брокеров
	brokers := h.client.Brokers()
	if len(brokers) == 0 {
		h.logger.Error("No Kafka brokers available")
		return err
	}

	// Проверяем подключение к каждому брокеру
	for _, broker := range brokers {
		if err := broker.Open(nil); err != nil {
			h.logger.WithError(err).WithField("broker_id", broker.ID()).Error("Failed to connect to Kafka broker")
			return err
		}
		broker.Close()
	}

	return nil
}

func (h *HealthChecker) Close() error {
	return h.client.Close()
}
