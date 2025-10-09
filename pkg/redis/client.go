package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type Client struct {
	client *redis.Client
	logger *logrus.Logger
}

func NewClient(addr, password string, db int, logger *logrus.Logger) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Client{
		client: rdb,
		logger: logger,
	}
}

// сохраняет значение в Redis с TTL
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		c.logger.WithError(err).Error("Failed to marshal value for Redis")
		return err
	}

	err = c.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		c.logger.WithError(err).WithField("key", key).Error("Failed to set Redis key")
		return err
	}

	c.logger.WithField("key", key).Debug("Successfully set Redis key")
	return nil
}

// получает значение из Redis и десериализует в dest
func (c *Client) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			c.logger.WithField("key", key).Debug("Redis key not found")
			return err
		}
		c.logger.WithError(err).WithField("key", key).Error("Failed to get Redis key")
		return err
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		c.logger.WithError(err).WithField("key", key).Error("Failed to unmarshal Redis value")
		return err
	}

	c.logger.WithField("key", key).Debug("Successfully retrieved Redis key")
	return nil
}

// удаляет ключ из Redis
func (c *Client) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		c.logger.WithError(err).WithField("key", key).Error("Failed to delete Redis key")
		return err
	}

	c.logger.WithField("key", key).Debug("Successfully deleted Redis key")
	return nil
}

// проверяет существование ключа в Redis
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		c.logger.WithError(err).WithField("key", key).Error("Failed to check Redis key existence")
		return false, err
	}
	return count > 0, nil
}

// проверяет подключение к Redis
func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.client.Close()
}
