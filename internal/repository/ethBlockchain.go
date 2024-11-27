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

// EthBlockchain
type EthBlockchain struct {
	rpcEndpoint string
}

// NewEthBlockchain
func NewEthBlockchain(url string) *EthBlockchain {
	return &EthBlockchain{
		rpcEndpoint: url,
	}
}

// GetCurrentBlock retrieves the latest block by number
func (b *EthBlockchain) GetCurrentBlock() (int64, error) {
	requestPayload := model.RpcRequest{
		JsonRPC: "2.0",
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		ID:      1,
	}

	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		return 0, err
	}

	resp, err := http.Post(b.rpcEndpoint, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error making request: %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return 0, err
	}

	var rpcResp model.RpcResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return 0, err
	}
	return hexToInt64(rpcResp.Result)
}

// GetBlockByNumber retrieves a block by its number
func (b *EthBlockchain) GetBlockByNumber(blockNumber int64) (model.Block, error) {
	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{fmt.Sprintf("0x%x", blockNumber), true},
		"id":      1,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		return model.Block{}, err
	}

	resp, err := http.Post(b.rpcEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error making request: %v", err)
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

func hexToInt64(hex string) (int64, error) {
	blockHex := strings.TrimPrefix(hex, "0x")
	return strconv.ParseInt(blockHex, 16, 64)
}
