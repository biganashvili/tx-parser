package service

import (
	"context"
	"log"
	"time"
	"tx-parser/internal/repository"
)

type EthParser struct {
	storage    repository.StorageInterface
	blockchain repository.BlockchainInterface
}

func NewEthParser(ctx context.Context, storage repository.StorageInterface, blockchain repository.BlockchainInterface) (*EthParser, error) {
	return &EthParser{storage: storage, blockchain: blockchain}, nil
}

func (ep *EthParser) Run(live bool) {
	currentBlockHeight := int64(0)
	var err error
	if live {
		//continue from blockchain current block
		currentBlockHeight, err = ep.blockchain.GetCurrentBlock()
	} else {
		//continue from db current block
		currentBlockHeight, err = ep.storage.GetCurrentBlock()
		currentBlockHeight++
	}

	if err != nil {
		log.Println(1, err)
	}

	for {

		block, err := ep.blockchain.GetBlockByNumber(currentBlockHeight)
		if err != nil {
			log.Println(2, err)
			continue
		}
		if block.Number == "" {
			log.Printf("block with height %d does not exist\n", currentBlockHeight)
			time.Sleep(5 * time.Second)
			continue
		}

		subscriptions, err := ep.storage.GetAllSubscriptions()
		if err != nil {
			log.Println(3, err)
			continue
		}
		log.Println("processing block ", currentBlockHeight)
		for _, tx := range block.Transactions {
			// log.Println(tx)
			if _, ok := subscriptions[tx.To]; ok {
				err := ep.storage.SaveTransaction(tx.To, tx)
				if err != nil {
					log.Println(4, err)
					continue
				}
			}
			if _, ok := subscriptions[tx.From]; ok {
				log.Println("from ", tx.From, tx.Value)
				err := ep.storage.SaveTransaction(tx.From, tx)
				if err != nil {
					log.Println(5, err)
					continue
				}
			}

		}
		err = ep.storage.SaveBlock(currentBlockHeight)
		if err != nil {
			log.Println(6, err)
			continue
		}
		currentBlockHeight++
	}
}
