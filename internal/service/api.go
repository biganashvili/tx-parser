package service

import "tx-parser/internal/model"

type ApiInterface interface {
	// last parsed block
	GetCurrentBlock() (int64, error)
	// add address to observer
	Subscribe(address string) (bool, error)
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) ([]model.Transaction, error)
}
