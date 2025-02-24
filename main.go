package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BidRequest struct {
	ImpressionID string `json:"impression_id"`
	Timeout      int    `json:"timeout_ms"`
}

type BidResponse struct {
	BidderID int     `json:"bidder_id"`
	Amount   float64 `json:"amount"`
}

// HandleBidRequest simulates a Prebid request
func HandleBidRequest(w http.ResponseWriter, r *http.Request) {
	var req BidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Simulate bid processing
	resp := BidResponse{
		BidderID: 1,
		Amount:   5.25,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/bid", HandleBidRequest)
	fmt.Println("Serveur démarré sur :8080")
	http.ListenAndServe(":8080", nil)
}
