package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// ServerConfig конфигурация для gRPC сервера
type ServerConfig struct {
	Port            int
	MaxRecvMsgSize  int
	MaxSendMsgSize  int
	MaxConcurrentStreams uint32
	KeepAliveTime   time.Duration
	KeepAliveTimeout time.Duration
	Logger          *logrus.Logger
}

// DefaultServerConfig возвращает конфигурацию по умолчанию
func DefaultServerConfig(port int, logger *logrus.Logger) *ServerConfig {
	return &ServerConfig{
		Port:                port,
		MaxRecvMsgSize:      4 * 1024 * 1024, // 4MB
		MaxSendMsgSize:      4 * 1024 * 1024, // 4MB
		MaxConcurrentStreams: 100,
		KeepAliveTime:       30 * time.Second,
		KeepAliveTimeout:     5 * time.Second,
		Logger:              logger,
	}
}

// Server обертка для gRPC сервера
type Server struct {
	server *grpc.Server
	config *ServerConfig
	logger *logrus.Logger
}

// NewServer создает новый gRPC сервер с interceptors
func NewServer(config *ServerConfig) *Server {
	// Настройки keepalive для предотвращения разрыва соединений
	kaep := keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second,
		PermitWithoutStream: true,
	}

	kaepServer := keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second,
		MaxConnectionAge:      30 * time.Second,
		MaxConnectionAgeGrace: 5 * time.Second,
		Time:                  config.KeepAliveTime,
		Timeout:               config.KeepAliveTimeout,
	}

	// Создаем опции для сервера
	opts := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kaepServer),
		grpc.MaxRecvMsgSize(config.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(config.MaxSendMsgSize),
		grpc.MaxConcurrentStreams(config.MaxConcurrentStreams),
		// Добавляем interceptors в правильном порядке
		grpc.UnaryInterceptor(chainUnaryInterceptors(
			RecoveryInterceptor(config.Logger),
			RequestIDInterceptor(),
			LoggingInterceptor(config.Logger),
			MetricsInterceptor(),
			TimeoutInterceptor(30 * time.Second),
		)),
	}

	server := grpc.NewServer(opts...)

	// Включаем reflection для gRPC UI инструментов (grpcurl, etc)
	reflection.Register(server)

	return &Server{
		server: server,
		config: config,
		logger: config.Logger,
	}
}

// GetServer возвращает внутренний gRPC сервер для регистрации сервисов
func (s *Server) GetServer() *grpc.Server {
	return s.server
}

// Start запускает gRPC сервер
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	s.logger.WithField("port", s.config.Port).Info("Starting gRPC server")

	go func() {
		if err := s.server.Serve(lis); err != nil {
			s.logger.WithError(err).Fatal("gRPC server failed")
		}
	}()

	return nil
}

// Stop останавливает gRPC сервер gracefully
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping gRPC server")

	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("gRPC server stopped gracefully")
		return nil
	case <-ctx.Done():
		s.logger.Warn("gRPC server forced to stop")
		s.server.Stop()
		return ctx.Err()
	}
}

// chainUnaryInterceptors объединяет несколько interceptors в один
func chainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = func(interceptor grpc.UnaryServerInterceptor, next grpc.UnaryHandler) grpc.UnaryHandler {
				return func(ctx context.Context, req interface{}) (interface{}, error) {
					return interceptor(ctx, req, info, next)
				}
			}(interceptors[i], chain)
		}
		return chain(ctx, req)
	}
}

