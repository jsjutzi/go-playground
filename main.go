package main

import (
	"du-service/config"
	"du-service/health"
	"du-service/importers"
	"du-service/utils"
	"log"
	"net/http"
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
        ReadTimeout:  1 * time.Hour,   // Long read timeout to keep connection open
        WriteTimeout: 0,               // No write timeout for SSE
    }

    regularServer := &http.Server{
        Addr:         ":8081",
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
    }

	go func() {
		log.Println("Starting SSE server on port 8080")
        if err := sseServer.ListenAndServe(); err != nil {
            log.Fatal("SSE Server failed:", err)
        }
	}()

	log.Println("Starting regular server on port 8081")
    if err := regularServer.ListenAndServe(); err != nil {
        log.Fatal("Regular Server failed:", err)
    }
}

func setupAPI() (*mux.Router, error) {
	// Define shared clients here - these will be passed to the handlers
	dynamoDbClient := config.DynamoDBClient{}
	opensearchClient := config.OpensearchClient{}
	s3Client := config.S3Client{}
	kcAdminClient := config.KcAdminClient{}

	sharedClients := &config.SharedClients{
		DynamoDbclient: &dynamoDbClient,
		OpensearchClient: &opensearchClient,
		S3Client: &s3Client,
		KcAdminClient: &kcAdminClient,
	}

	// Create new event emitter
	eventEmitter := utils.NewEventEmitter()
	go eventEmitter.Start()

	// Define worker pool - 100 workers for demo purposes
	workerPool := utils.NewWorkerPool(100)

	// Define routes
	r := mux.NewRouter()

	// Health check endpoint
	r.HandleFunc("/health", health.HealthCheckHandler).Methods("GET")

	// Bulk import endpoint
	r.HandleFunc("/import-users", importers.ImportsHandler(eventEmitter, workerPool, sharedClients)).Methods("GET")

	// // Import Single User endpoint
	// http.HandleFunc("/import-user", importers.ImportsHandler(eventEmitter, workerPool))

	// // Bulk Removal Endpoint
	// http.HandleFunc("/remove-users", importers.ImportsHandler(eventEmitter, workerPool))

	return r, nil
}
