package service

import "tx-parser/internal/model"

type ApiInterface interface {
	GetCurrentBlock() (int64, error)

	Subscribe(address string) (bool, error)

	GetTransactions(address string) ([]model.Transaction, error)
}
