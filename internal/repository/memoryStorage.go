package repository

import (
	"strings"
	"sync"
	"tx-parser/internal/model"
)

// MemoryStorage is an in-memory storage for block-related data
type MemoryStorage struct {
	mu           sync.RWMutex
	currentBlock int64           // Latest block number processed by the listener
	subscribers  map[string]bool // Subscribed addresses
	transactions map[string]map[string]model.Transaction
}

// NewMemoryStorage initializes an in-memory block storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		transactions: make(map[string]map[string]model.Transaction),
		subscribers:  make(map[string]bool),
		currentBlock: 0, // Initialize to 0 or any appropriate value
	}
}

// SaveBlock stores the block number for a given address
func (s *MemoryStorage) SaveBlock(blockNum int64) error {
	s.mu.Lock() // Use lock to protect write access
	defer s.mu.Unlock()
	s.currentBlock = blockNum
	return nil
}

// GetCurrentBlock retrieves the latest block number
func (s *MemoryStorage) GetCurrentBlock() (int64, error) {
	s.mu.RLock() // Use read lock for thread-safe reading
	defer s.mu.RUnlock()
	return s.currentBlock, nil
}

// SaveTransaction stores a transaction for an address
func (s *MemoryStorage) SaveTransaction(address string, tx model.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.transactions[address] == nil {
		s.transactions[address] = make(map[string]model.Transaction)
	}
	s.transactions[address][tx.Hash] = tx
	return nil
}

// Get All Subscription
func (s *MemoryStorage) GetAllSubscriptions() (map[string]bool, error) {
	return s.subscribers, nil
}

// GetTransactions retrieves all transactions for a given address
func (s *MemoryStorage) GetTransactions(address string) ([]model.Transaction, error) {
	s.mu.RLock() // Read lock for concurrent reads
	defer s.mu.RUnlock()
	res := []model.Transaction{}
	for _, v := range s.transactions[address] {
		res = append(res, v)
	}
	return res, nil
}

// Subscribe adds an address to the list of observed addresses
func (s *MemoryStorage) Subscribe(address string) (bool, error) {
	if _, exists := s.subscribers[strings.ToLower(address)]; !exists {
		s.subscribers[strings.ToLower(address)] = true // Mark the address as subscribed
		return true, nil
	}

	return false, nil
}