package postgres

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// проверяет состояние бд
type HealthChecker struct {
	db      *sql.DB
	logger  *logrus.Logger
	timeout time.Duration
}

// создает новый health checker для бд
func NewHealthChecker(dsn string, logger *logrus.Logger) (*HealthChecker, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &HealthChecker{
		db:      db,
		logger:  logger,
		timeout: 5 * time.Second,
	}, nil
}

// проверяет состояние бд
func (h *HealthChecker) CheckHealth(ctx context.Context) error {
	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		h.logger.WithError(err).Error("Failed to ping PostgreSQL database")
		return err
	}

	// Проверяем, что можем выполнить простой запрос
	var result int
	query := "SELECT 1"
	if err := h.db.QueryRowContext(ctx, query).Scan(&result); err != nil {
		h.logger.WithError(err).Error("Failed to execute test query on PostgreSQL")
		return err
	}

	// Проверяем, что база данных не в режиме только для чтения
	var isReadOnly bool
	query = "SELECT pg_is_in_recovery()"
	if err := h.db.QueryRowContext(ctx, query).Scan(&isReadOnly); err != nil {
		h.logger.WithError(err).Error("Failed to check PostgreSQL recovery status")
		return err
	}

	if isReadOnly {
		h.logger.Warn("PostgreSQL database is in read-only mode")
	}

	// Проверяем доступность таблицы transactions
	query = "SELECT COUNT(*) FROM transactions LIMIT 1"
	if err := h.db.QueryRowContext(ctx, query).Scan(&result); err != nil {
		h.logger.WithError(err).Error("Failed to access transactions table")
		return err
	}

	return nil
}

func (h *HealthChecker) Close() error {
	return h.db.Close()
}
