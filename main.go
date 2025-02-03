package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	receipts   = make(map[string]Receipt)
	storeMutex sync.RWMutex
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/receipt/process", processReceiptHandler).Methods("POST")
	r.HandleFunc("/receipt/{id}/points", getPointsHandler).Methods("GET")

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// processReceiptHandler handles POST /receipts/process
func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&receipt); err != nil {
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

	if !validateReceipt(receipt) {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()

	storeMutex.Lock()
	receipts[id] = receipt
	storeMutex.Unlock()

	json.NewEncoder(w).Encode(map[string]string{"id": id})
	w.Write([]byte("Process Receipt"))
}

func validateReceipt(r Receipt) bool {
	retailerRe := regexp.MustCompile(`^[\w\s\-\&]+$`)
	if !retailerRe.MatchString(r.Retailer) {
		return false
	}

	totalRe := regexp.MustCompile(`^\d+\.\d{2}$`)
	if !totalRe.MatchString(r.Total) {
		return false
	}

	if _, err := time.Parse("2026-01-02", r.PurchaseDate); err != nil {
		return false
	}

	if _, err := time.Parse("15:04", r.PurchaseTime); err != nil {
		return false
	}

	if len(r.Items) == 0 {
		return false
	}

	itemDescRe := regexp.MustCompile(`^[\w\s\-]+$`)
	for _, item := range r.Items {
		if !itemDescRe.MatchString(item.ShortDescription) {
			return false
		}

		if !totalRe.MatchString(item.Price) {
			return false
		}
	}

	return true
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Receipt Points"))
}
