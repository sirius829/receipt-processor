package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", processReceiptHandler).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")
	router.Use(loggingMiddleware)
	return router
}

func TestReceiptProcessingAndPointsRetrieval(t *testing.T) {
	receipts = make(map[string]Receipt)

	router := setupRouter()

	tests := []struct {
		name           string
		receiptJSON    string
		expectedPoints int
	}{
		{
			name: "Example 1",
			receiptJSON: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [
					{"shortDescription": "Mountain Dew 12PK", "price": "6.49"},
					{"shortDescription": "Emils Cheese Pizza", "price": "12.25"},
					{"shortDescription": "Knorr Creamy Chicken", "price": "1.26"},
					{"shortDescription": "Doritos Nacho Cheese", "price": "3.35"},
					{"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ", "price": "12.00"}
				],
				"total": "35.35"
			}`,
			expectedPoints: 28,
		},
		{
			name: "Example 2",
			receiptJSON: `{
				"retailer": "M&M Corner Market",
				"purchaseDate": "2022-03-20",
				"purchaseTime": "14:33",
				"items": [
					{"shortDescription": "Gatorade", "price": "2.25"},
					{"shortDescription": "Gatorade", "price": "2.25"},
					{"shortDescription": "Gatorade", "price": "2.25"},
					{"shortDescription": "Gatorade", "price": "2.25"}
				],
				"total": "9.00"
			}`,
			expectedPoints: 109,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(tc.receiptJSON)))
			if err != nil {
				t.Fatalf("Failed to create POST request: %v", err)
			}

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("Expected status 200 for POST, got %d", rr.Code)
			}

			var postResp map[string]string

			if err := json.Unmarshal(rr.Body.Bytes(), &postResp); err != nil {
				t.Fatalf("Failed to unmarshal POST response: %v", err)
			}

			receiptID, ok := postResp["id"]
			if !ok || receiptID == "" {
				t.Fatalf("Response missing receipt ID: %+v", postResp)
			}

			reqGet, err := http.NewRequest("GET", "/receipts/"+receiptID+"/points", nil)
			if err != nil {
				t.Fatalf("Failed to create Get request: %v", err)
			}

			rrGet := httptest.NewRecorder()
			router.ServeHTTP(rrGet, reqGet)

			if rrGet.Code != http.StatusOK {
				t.Fatalf("Expected status 200 for GET, got %d", rrGet.Code)
			}

			var pointsResp map[string]int
			if err := json.Unmarshal(rrGet.Body.Bytes(), &pointsResp); err != nil {
				t.Fatalf("Failed to unmarshal GET response: %v", err)
			}

			points, exists := pointsResp["points"]
			if !exists {
				t.Fatalf("GET response missing 'points': %+v", pointsResp)
			}

			if points != tc.expectedPoints {
				t.Fatalf("Expected %d points, got %d", tc.expectedPoints, points)
			}
		})
	}
}
