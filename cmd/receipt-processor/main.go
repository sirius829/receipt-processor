package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"receipt-processor/internal/handlers"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(os.Stdout, "receipt-processor: ", log.LstdFlags)

	r := mux.NewRouter()
	h := handlers.NewHandler(logger)
	r.HandleFunc("/receipts/process", h.ProcessReceiptHandler).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", h.GetPointsHandler).Methods("GET")
	r.Use(handlers.LoggingMiddleware(logger))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		logger.Println("Server starting on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server Shutdown Failed:%+v", err)
	}
	logger.Println("Server gracefully stopped")
}
