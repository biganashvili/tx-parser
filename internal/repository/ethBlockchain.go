package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"tx-parser/internal/model"
)

// MemoryStorage is an in-memory storage for block-related data
type EthBlockchain struct {
	rpcEndpoint string
}

// NewMemoryStorage initializes an in-memory block storage
func NewEthBlockchain(url string) *EthBlockchain {
	return &EthBlockchain{
		rpcEndpoint: url,
	}
}

// GetLatestETHBlock retrieves the latest Ethereum block number via RPC
func (b *EthBlockchain) GetCurrentBlock() (int64, error) {
	requestPayload := model.RpcRequest{
		JsonRPC: "2.0",
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		ID:      1,
	}

	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		log.Printf("Error marshaling request payload: %v", err)
		return 0, err
	}

	resp, err := http.Post(b.rpcEndpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error making RPC request: %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading RPC response: %v", err)
		return 0, err
	}

	var rpcResp model.RpcResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		log.Printf("Error unmarshaling RPC response: %v", err)
		return 0, err
	}
	return hexToInt64(rpcResp.Result)
}

// GetEthBlockByNumber retrieves a block by its number via RPC
func (b *EthBlockchain) GetBlockByNumber(blockNumber int64) (model.Block, error) {
	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{fmt.Sprintf("0x%x", blockNumber), true},
		"id":      1,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Error marshaling request payload: %v", err)
		return model.Block{}, err
	}

	resp, err := http.Post(b.rpcEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error making RPC request: %v", err)
		return model.Block{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return model.Block{}, err
	}
	var rpcResp model.JsonRPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		log.Printf("Error unmarshaling block response: %v", err)
		return model.Block{}, err
	}

	return rpcResp.Result, nil
}

/*
// FilterTransactionsByAddress filters transactions for subscribed addresses
func (b *EthBlockchain) FilterTransactionsByAddress(transactions []model.Transaction) error {
	addressMap := model.SharedStore().GetAllSubscriptions()
	if len(addressMap) == 0 {
		// Map is empty, handle the empty case here
		fmt.Println("No subscriptions found.")
		return nil
	}

	for _, tx := range transactions {
		lowerFrom := strings.ToLower(tx.From)
		lowerTo := strings.ToLower(tx.To)

		// Check if From or To address is subscribed
		if addressMap[lowerFrom] {
			if err := model.SharedStore().SaveTransaction(lowerFrom, tx); err != nil {
				log.Printf("Error saving transaction for address %s: %v", lowerFrom, err)
				return err
			}
		}

		if addressMap[lowerTo] && lowerFrom != lowerTo {
			if err := model.SharedStore().SaveTransaction(lowerTo, tx); err != nil {
				log.Printf("Error saving transaction for address %s: %v", lowerTo, err)
				return err
			}
		}
	}

	return nil
}

// IncrementBlockNumber increments and stores the block number
func (b *EthBlockchain) IncrementBlockNumber(blockHex string) error {
	blockHex = strings.TrimPrefix(blockHex, "0x")
	blockNumber, err := strconv.ParseInt(blockHex, 16, 64)
	if err != nil {
		log.Printf("Error parsing hex block number: %v", err)
		return errors.New("error parsing block number")
	}

	blockNumber++
	newBlockHex := fmt.Sprintf("0x%x", blockNumber)
	if err := SaveLatestBlock(newBlockHex); err != nil {
		log.Printf("Error saving incremented block number: %v", err)
		return err
	}

	log.Printf("Successfully incremented block number to: %d %s", blockNumber, newBlockHex)
	return nil
}

func isValidEthereumAddress(address string) bool {
	return strings.HasPrefix(address, "0x") && len(address) == 42
}

// Subscribe adds an address to the list of observed addresses
func (b *EthBlockchain) Subscribe(address string) (bool, error) {
	if !isValidEthereumAddress(address) {
		fmt.Println("Invalid Ethereum address")
		return false, errors.New("invalid Ethereum address")
	}
	return model.SharedStore().Subscribe(address)
}

*/

func hexToInt64(hex string) (int64, error) {
	blockHex := strings.TrimPrefix(hex, "0x")
	return strconv.ParseInt(blockHex, 16, 64)
}
