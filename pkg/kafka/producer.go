package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
	logger   *logrus.Logger
}

func NewProducer(brokers []string, topic string, logger *logrus.Logger) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 10 * time.Second

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topic:    topic,
		logger:   logger,
	}, nil
}

func (p *Producer) SendMessage(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		p.logger.WithError(err).Error("Failed to marshal message")
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		p.logger.WithError(err).Error("Failed to send message to Kafka")
		return err
	}

	p.logger.WithFields(logrus.Fields{
		"topic":     p.topic,
		"partition": partition,
		"offset":    offset,
		"key":       key,
	}).Info("Message sent to Kafka successfully")

	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
