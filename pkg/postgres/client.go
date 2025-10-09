package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/nikitakolesnik/pet-proj/internal/models"
	"github.com/sirupsen/logrus"
)

type Client struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewClient(host string, port int, database, username, password string, logger *logrus.Logger) (*Client, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Client{
		db:     db,
		logger: logger,
	}, nil
}

func (c *Client) CreateTransactionTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		kafka_status VARCHAR(10) NOT NULL CHECK (kafka_status IN ('ok', 'bad')),
		redis_status VARCHAR(10) NOT NULL CHECK (redis_status IN ('ok', 'bad')),
		duration_ms BIGINT NOT NULL CHECK (duration_ms >= 0),
		service VARCHAR(50) NOT NULL,
		event_id VARCHAR(100),
		error_msg TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);`

	_, err := c.db.Exec(query)
	if err != nil {
		c.logger.WithError(err).Error("Failed to create transactions table")
		return err
	}

	c.logger.Info("Transactions table created successfully")
	return nil
}

func (c *Client) InsertTransaction(tx *models.Transaction) error {
	query := `
	INSERT INTO transactions (timestamp, kafka_status, redis_status, duration_ms, service, event_id, error_msg)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := c.db.Exec(query, tx.Timestamp, tx.KafkaStatus, tx.RedisStatus,
		tx.Duration, tx.Service, tx.EventID, tx.ErrorMsg)

	if err != nil {
		c.logger.WithError(err).Error("Failed to insert transaction")
		return err
	}

	c.logger.WithField("event_id", tx.EventID).Debug("Transaction inserted successfully")
	return nil
}

func (c *Client) GetTransactions(limit int) ([]*models.Transaction, error) {
	query := `SELECT id, timestamp, kafka_status, redis_status, duration_ms, service, event_id, error_msg
			  FROM transactions ORDER BY timestamp DESC LIMIT $1`

	rows, err := c.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		tx := &models.Transaction{}
		err := rows.Scan(&tx.ID, &tx.Timestamp, &tx.KafkaStatus, &tx.RedisStatus,
			&tx.Duration, &tx.Service, &tx.EventID, &tx.ErrorMsg)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (c *Client) GetTransactionStats() (map[string]interface{}, error) {
	query := `
	SELECT 
		service,
		COUNT(*) as total_transactions,
		COUNT(CASE WHEN kafka_status = 'ok' THEN 1 END) as kafka_success,
		COUNT(CASE WHEN kafka_status = 'bad' THEN 1 END) as kafka_failures,
		COUNT(CASE WHEN redis_status = 'ok' THEN 1 END) as redis_success,
		COUNT(CASE WHEN redis_status = 'bad' THEN 1 END) as redis_failures,
		AVG(duration_ms) as avg_duration_ms
	FROM transactions
	WHERE timestamp >= NOW() - INTERVAL '1 hour'
	GROUP BY service`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]interface{})
	for rows.Next() {
		var service string
		var total, kafkaSuccess, kafkaFailures, redisSuccess, redisFailures int64
		var avgDuration sql.NullFloat64

		err := rows.Scan(&service, &total, &kafkaSuccess, &kafkaFailures, &redisSuccess, &redisFailures, &avgDuration)
		if err != nil {
			return nil, err
		}

		stats[service] = map[string]interface{}{
			"total_transactions": total,
			"kafka_success":      kafkaSuccess,
			"kafka_failures":     kafkaFailures,
			"redis_success":      redisSuccess,
			"redis_failures":     redisFailures,
			"avg_duration_ms":    avgDuration.Float64,
		}
	}

	return stats, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}
