package service

import (
	"context"
	"tx-parser/internal/model"
	"tx-parser/internal/repository"
)

type ApiService struct {
	storage repository.StorageInterface
}

func NewApi(ctx context.Context, storage repository.StorageInterface) (*ApiService, error) {
	return &ApiService{storage: storage}, nil
}

func (s *ApiService) GetCurrentBlock() (int64, error) {
	currentBlockHeight, err := s.storage.GetCurrentBlock()
	return currentBlockHeight, err
}

// add address to observer
func (s *ApiService) Subscribe(address string) (bool, error) {
	return s.storage.Subscribe(address)
}

// list of transactions for an address
func (s *ApiService) GetTransactions(address string) ([]model.Transaction, error) {
	return s.storage.GetTransactions(address)
}
