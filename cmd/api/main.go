package main

import (
	"context"
	"du-import-service/config"
	"du-import-service/internal/importers"
	"du-import-service/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r, err := setupAPI()

	if err != nil {
		log.Fatalf("Failed to setup API: %v", err)
	}

	sseServer := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  1 * time.Hour, // Long read timeout to keep connection open
		WriteTimeout: 0,             // No write timeout for SSE
	}

	regularServer := &http.Server{
		Addr:         ":8081",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Println("Starting SSE server on port 8080")
		if err := sseServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("SSE Server failed: %v", err)
		}
	}()

	go func() {
		log.Println("Starting regular server on port 8081")
		if err := regularServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Regular Server failed: %v", err)
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Timeout context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := sseServer.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown SSE server: %v", err)
	}

	if err := regularServer.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown regular server: %v", err)
	}
}

func setupAPI() (*mux.Router, error) {
	// Define shared clients here - these will be passed to the handlers
	// TODO: Move this to the handlers via a utility function, so each import handler
	// can create its own clients, then destroy them once import is complete?
	dynamoDbClient := config.DynamoDBClient{}
	opensearchClient := config.OpensearchClient{}
	s3Client := config.S3Client{}
	kcAdminClient := config.KcAdminClient{}

	sharedClients := &config.SharedClients{
		DynamoDbclient:   &dynamoDbClient,
		OpensearchClient: &opensearchClient,
		S3Client:         &s3Client,
		KcAdminClient:    &kcAdminClient,
	}

	// Create new event emitter
	eventEmitter := utils.NewEventEmitter()
	go eventEmitter.Start()

	// Define worker pool - 100 workers for demo purposes
	workerPool := utils.NewWorkerPool(100)

	// Define routes
	r := mux.NewRouter()

	// Health check endpoint
	r.HandleFunc("/health", healthCheckHandler).Methods("GET")

	// Bulk import endpoint
	r.HandleFunc("/import-users", importers.ImportsHandler(eventEmitter, workerPool, sharedClients)).Methods("GET")

	// // Import Single User endpoint
	// http.HandleFunc("/import-user", importers.ImportsHandler(eventEmitter, workerPool))

	// // Bulk Removal Endpoint
	// http.HandleFunc("/remove-users", importers.ImportsHandler(eventEmitter, workerPool))

	return r, nil
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := checkApplicationHealth()

	if !status {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	response := map[string]string{"status": "healthy"}

	if !status {
		response["status"] = "unhealthy"
	}

	json.NewEncoder(w).Encode(response)
}

func checkApplicationHealth() bool {
	// Implement actual checks here, e.g., database ping, etc.
	return true // return false if any check fails
}
