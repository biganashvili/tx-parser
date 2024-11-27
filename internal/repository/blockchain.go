package repository

import "tx-parser/internal/model"

type BlockchainInterface interface {
	GetCurrentBlock() (int64, error)
	GetBlockByNumber(int64) (model.Block, error)
}
