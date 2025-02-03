package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/receipt/process", processReceiptHandler).Methods("POST")
	r.HandleFunc("/receipt/{id}/points", getPointsHandler).Methods("GET")

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Receipt Processor Service"))
	// })

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Process Receipt"))
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Receipt Points"))
}
