package postgres

import (
	"pet-proj/internal/models"
)

type ClientInterface interface {
	InsertTransaction(tx *models.Transaction) error
	GetTransactions(limit int) ([]*models.Transaction, error)
	GetTransactionStats() (map[string]interface{}, error)
	CreateTransactionTable() error
	Close() error
}
