package model

// Block represents an Ethereum block
type Block struct {
	Number       string
	Transactions []Transaction
}
