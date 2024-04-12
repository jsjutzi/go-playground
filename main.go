package main

import (
	"du-service/health"
	"du-service/importers"
	"du-service/utils"
	"fmt"
	"log"
	"net/http"
)

func main() {
	setupAPI()
	fmt.Printf("Starting server at port 8080\n")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func setupAPI() {
	// Create new event emitter
	eventEmitter := utils.NewEventEmitter()
	go eventEmitter.Start()

	// Define worker pool - 100 workers for demo purposes
	workerPool := utils.NewWorkerPool(100)

	// Define routes

	// Health check endpoint
	http.HandleFunc("/health", health.HealthCheckHandler)

	// Bulk import endpoint
	http.HandleFunc("/import-users", importers.ImportsHandler(eventEmitter, workerPool))

	// Import Single User endpoint
	http.HandleFunc("/import-user", importers.ImportsHandler(eventEmitter, workerPool))

	// Bulk Removal Endpoint
	http.HandleFunc("/remove-users", importers.ImportsHandler(eventEmitter, workerPool))
}
