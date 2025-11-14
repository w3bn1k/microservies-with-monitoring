package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor логирует все gRPC запросы
func LoggingInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Извлекаем метаданные из контекста
		md, _ := metadata.FromIncomingContext(ctx)
		requestID := getRequestID(md)

		logger.WithFields(logrus.Fields{
			"method":     info.FullMethod,
			"request_id": requestID,
		}).Info("gRPC request started")

		// Выполняем запрос
		resp, err := handler(ctx, req)

		// Логируем результат
		duration := time.Since(start)
		fields := logrus.Fields{
			"method":     info.FullMethod,
			"duration_ms": duration.Milliseconds(),
			"request_id": requestID,
		}

		if err != nil {
			st, _ := status.FromError(err)
			fields["error"] = err.Error()
			fields["code"] = st.Code().String()
			logger.WithFields(fields).Error("gRPC request failed")
		} else {
			logger.WithFields(fields).Info("gRPC request completed")
		}

		return resp, err
	}
}

// MetricsInterceptor собирает метрики для Prometheus
func MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Выполняем запрос
		resp, err := handler(ctx, req)

		// Обновляем метрики
		duration := time.Since(start).Seconds()
		code := "OK"
		if err != nil {
			st, _ := status.FromError(err)
			code = st.Code().String()
		}

		// Метрики будут обновляться через monitoring пакет
		_ = duration
		_ = code
		_ = info.FullMethod

		return resp, err
	}
}

// TimeoutInterceptor устанавливает таймаут для запросов
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return handler(ctx, req)
	}
}

// RecoveryInterceptor обрабатывает паники
func RecoveryInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		defer func() {
			if r := recover(); r != nil {
				logger.WithFields(logrus.Fields{
					"method": info.FullMethod,
					"panic":  r,
				}).Error("gRPC panic recovered")
			}
		}()

		return handler(ctx, req)
	}
}

// RequestIDInterceptor добавляет request ID в контекст
func RequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		requestID := getRequestID(md)
		if requestID == "" {
			requestID = generateRequestID()
			md.Set("x-request-id", requestID)
		}

		ctx = metadata.NewIncomingContext(ctx, md)
		return handler(ctx, req)
	}
}

// getRequestID извлекает request ID из метаданных
func getRequestID(md metadata.MD) string {
	if values := md.Get("x-request-id"); len(values) > 0 {
		return values[0]
	}
	return ""
}

// generateRequestID генерирует новый request ID
func generateRequestID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix())
}

