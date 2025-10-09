package kafka

import "context"

type ProducerInterface interface {
	SendMessage(ctx context.Context, key string, value interface{}) error
	Close() error
}

type ConsumerInterface interface {
	SetHandler(handler MessageHandler)
	Start(ctx context.Context) error
	Close() error
}
