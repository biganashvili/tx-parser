package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"tx-parser/internal/repository"
	"tx-parser/internal/service"
)

var parser service.ParserInterface
var api service.ApiInterface

func main() {
	ctx := context.Background()
	var err error
	storage := repository.NewMemoryStorage()
	api, err = service.NewApi(
		ctx,
		storage,
	)
	if err != nil {
		panic(err)
	}
	go StartServer()
	_, _ = api.Subscribe("0xe52470bef1da70af094a91e326076c0bdca688ff")

	parser, err = service.NewEthParser(
		ctx,
		storage,
		repository.NewEthBlockchain("https://ethereum-rpc.publicnode.com"),
	)

	if err != nil {
		panic(err)
	}

	parser.Run(true)
}

// StartServer initializes and starts the HTTP server
func StartServer() error {

	http.HandleFunc("/currentBlock", CurrentBlockHandler)
	http.HandleFunc("/subscribe", SaveSubscriptionHandler)
	http.HandleFunc("/transactions", ListTransactionsHandler)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}
	return nil
}

// CurrentBlockHandler returns the current block number
func CurrentBlockHandler(w http.ResponseWriter, r *http.Request) {
	setJSONResponseHeaders(w)
	blockNumber, err := api.GetCurrentBlock()
	if err != nil {
		log.Printf("Error parsing hex block number: %v", err)
	}

	if err := json.NewEncoder(w).Encode(blockNumber); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("error encoding current block response: %v", err)
	}
}

// SaveSubscriptionHandler handles subscription requests for an Ethereum address
func SaveSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	setJSONResponseHeaders(w)

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address", http.StatusBadRequest)
		return
	}

	subscribed, err := api.Subscribe(address)

	status := "Already Subscribed"
	if subscribed {
		status = "Subscribed"
	}
	var response map[string]string
	if err != nil {
		response = map[string]string{"status": err.Error(), "address": address}

	} else {
		response = map[string]string{"status": status, "address": address}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("error encoding subscription response for address %s: %v", address, err)
	}
}

// ListTransactionsHandler returns a list of transactions for a given Ethereum address
func ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	setJSONResponseHeaders(w)

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address", http.StatusBadRequest)
		return
	}

	transactions, err := api.GetTransactions(strings.ToLower(address))
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Printf("db error current block response: %v", err)
	}
	// Check if transactions are empty
	if len(transactions) == 0 {
		// Return an empty array or a message, depending on your preference
		response := map[string]interface{}{
			"address":      address,
			"transactions": []string{}, // Returning an empty list to indicate no transactions
			"message":      "No transactions found for this address",
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("error encoding empty transactions response for address %s: %v", address, err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("error encoding transactions for address %s: %v", address, err)
	}
}

// setJSONResponseHeaders sets common headers for JSON responses
func setJSONResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allows cross-origin requests
}
