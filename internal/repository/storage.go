package repository

import "tx-parser/internal/model"

// StorageInterface
type StorageInterface interface {
	GetCurrentBlock() (int64, error)

	Subscribe(address string) (bool, error)

	GetTransactions(address string) ([]model.Transaction, error)

	SaveBlock(blockNumber int64) error

	GetAllSubscriptions() (map[string]bool, error)

	SaveTransaction(address string, tx model.Transaction) error
}
