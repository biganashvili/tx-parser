// internal/model/storage.go
package model

// JSON-RPC response for getting a block
type JsonRPCResponse struct {
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  Block  `json:"result"`
}

// RpcRequest represents the structure of an Ethereum JSON-RPC request
type RpcRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// RpcResponse represents the structure of an Ethereum JSON-RPC response
type RpcResponse struct {
	ID      int    `json:"id"`
	JsonRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
}
