package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

// ClientConfig конфигурация для gRPC клиента
type ClientConfig struct {
	Address           string
	MaxRecvMsgSize    int
	MaxSendMsgSize    int
	KeepAliveTime     time.Duration
	KeepAliveTimeout  time.Duration
	ConnectionTimeout time.Duration
	Logger            *logrus.Logger
}

// DefaultClientConfig возвращает конфигурацию по умолчанию
func DefaultClientConfig(address string, logger *logrus.Logger) *ClientConfig {
	return &ClientConfig{
		Address:           address,
		MaxRecvMsgSize:    4 * 1024 * 1024, // 4MB
		MaxSendMsgSize:    4 * 1024 * 1024, // 4MB
		KeepAliveTime:     30 * time.Second,
		KeepAliveTimeout:  5 * time.Second,
		ConnectionTimeout: 10 * time.Second,
		Logger:            logger,
	}
}

// Client обертка для gRPC клиента с connection pooling
type Client struct {
	conn   *grpc.ClientConn
	config *ClientConfig
	logger *logrus.Logger
}

// NewClient создает новый gRPC клиент с настройками для масштабируемости
func NewClient(config *ClientConfig) (*Client, error) {
	// Настройки keepalive для поддержания соединений
	kaep := keepalive.ClientParameters{
		Time:                config.KeepAliveTime,
		Timeout:             config.KeepAliveTimeout,
		PermitWithoutStream: true,
	}

	// Опции для клиента
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Для production использовать TLS
		grpc.WithKeepaliveParams(kaep),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(config.MaxRecvMsgSize),
			grpc.MaxCallSendMsgSize(config.MaxSendMsgSize),
		),
		// Retry политика встроена в gRPC
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(true),
		),
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, config.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", config.Address, err)
	}

	config.Logger.WithField("address", config.Address).Info("gRPC client connected")

	return &Client{
		conn:   conn,
		config: config,
		logger: config.Logger,
	}, nil
}

// GetConnection возвращает gRPC соединение
func (c *Client) GetConnection() *grpc.ClientConn {
	return c.conn
}

// Close закрывает соединение
func (c *Client) Close() error {
	c.logger.Info("Closing gRPC client connection")
	return c.conn.Close()
}

// WithRequestID добавляет request ID в контекст для трейсинга
func WithRequestID(ctx context.Context, requestID string) context.Context {
	md := metadata.New(map[string]string{
		"x-request-id": requestID,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

// GetRequestID извлекает request ID из контекста
func GetRequestID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md, _ = metadata.FromOutgoingContext(ctx)
	}
	if md != nil {
		if values := md.Get("x-request-id"); len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

