package model

// Block represents the structure of an Ethereum block
type Block struct {
	Number       string
	Transactions []Transaction
}
