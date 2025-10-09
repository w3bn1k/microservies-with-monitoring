package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	consumer sarama.ConsumerGroup
	topic    string
	logger   *logrus.Logger
	handler  MessageHandler
}

type MessageHandler interface {
	HandleMessage(ctx context.Context, message *sarama.ConsumerMessage) error
}

func NewConsumer(brokers []string, topic, groupID string, logger *logrus.Logger) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		topic:    topic,
		logger:   logger,
	}, nil
}

func (c *Consumer) SetHandler(handler MessageHandler) {
	c.handler = handler
}

func (c *Consumer) Start(ctx context.Context) error {
	topics := []string{c.topic}

	go func() {
		for err := range c.consumer.Errors() {
			c.logger.WithError(err).Error("Kafka consumer error")
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := c.consumer.Consume(ctx, topics, c)
			if err != nil {
				c.logger.WithError(err).Error("Error from consumer")
				return err
			}
		}
	}
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			if c.handler != nil {
				if err := c.handler.HandleMessage(session.Context(), message); err != nil {
					c.logger.WithError(err).Error("Failed to handle message")
				}
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
