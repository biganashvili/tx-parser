package repository

import "tx-parser/internal/model"

// storeInterface defines the behavior required for interacting with Ethereum-related data storage.
type StorageInterface interface {
	// GetCurrentBlock retrieves the number of the current Ethereum block being tracked.
	GetCurrentBlock() (int64, error)

	// Subscribe adds an Ethereum address to be tracked. Returns true if the subscription is new, false if the address is already subscribed.
	Subscribe(address string) (bool, error)

	// GetTransactions retrieves the list of transactions associated with a specific address.
	GetTransactions(address string) ([]model.Transaction, error)

	// SaveBlock persists the latest block number.
	SaveBlock(blockNumber int64) error

	// GetAllSubscriptions retrieves all Ethereum addresses that are currently being tracked.
	GetAllSubscriptions() (map[string]bool, error)

	// SaveTransaction saves a transaction associated with an Ethereum address.
	SaveTransaction(address string, tx model.Transaction) error
}
