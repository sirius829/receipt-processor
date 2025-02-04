package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
	r.HandleFunc("/receipts/process", processReceiptHandler).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")
	r.Use(loggingMiddleware)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request for %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
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

	if _, err := time.Parse("2006-01-02", r.PurchaseDate); err != nil {
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

// getPointsHandler handles GET /receipts/{id}/points
func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	storeMutex.RLock()
	receipt, ok := receipts[id]
	storeMutex.RUnlock()

	if !ok {
		http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
		return
	}

	points, err := calculatePoints(receipt)
	if err != nil {
		http.Error(w, "Error calculating points", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"points": points})
}

// calculatePoints calculate the points for given receipt and given condition on readme
func calculatePoints(r Receipt) (int, error) {
	points := 0

	// Rule 1: 1 point for every alphanumeric character in the retailer name
	for _, c := range r.Retailer {
		if isAlphanumeric(c) {
			points++
		}
	}

	totalVal, err := strconv.ParseFloat(r.Total, 64)
	if err != nil {
		return 0, err
	}

	// Rule 2: 50 points if the total is round dollar amount
	if math.Mod(totalVal, 1.0) == 0 {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if math.Mod(totalVal, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items
	points += (len(r.Items) / 2) * 5

	// Rule 5: For every item whose trimmed description length is a multiple of 3,
	// multiply it's price by 0.2 and round up
	for _, item := range r.Items {
		desc := strings.TrimSpace(item.ShortDescription)
		if len(desc) > 0 && len(desc)%3 == 0 {
			priceVal, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, err
			}
			bonus := int(math.Ceil(priceVal * 0.2))
			points += bonus
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd
	purchaseDate, err := time.Parse("2006-01-02", r.PurchaseDate)
	if err != nil {
		return 0, err
	}

	if purchaseDate.Day()%2 == 1 {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00 pm and before 4:00 pm
	purchaseTime, err := time.Parse("15:04", r.PurchaseTime)
	if err != nil {
		return 0, err
	}

	after2, _ := time.Parse("15:04", "14:00")
	before4, _ := time.Parse("15:04", "16:00")
	if purchaseTime.After(after2) && purchaseTime.Before(before4) {
		points += 10
	}

	return points, nil
}

// isAlphanumeric return true if the rune is a letter or digit
func isAlphanumeric(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9')
}
