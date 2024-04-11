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

	// Define worker pool - 20 workers
	workerPool := utils.NewWorkerPool(20)

	// Define routes

	// Health check endpoint
	http.HandleFunc("/health", health.HealthCheckHandler)

	// Bulk import endpoints
	http.HandleFunc("/import-users-mp", importers.ImportsHandler(eventEmitter, workerPool))
	http.HandleFunc("/import-users-cohort", importers.ImportsHandler(eventEmitter, workerPool))
	http.HandleFunc("/import-users-license", importers.ImportsHandler(eventEmitter, workerPool))

	// Import Single User endpoints
	http.HandleFunc("/import-user-mp", importers.ImportsHandler(eventEmitter, workerPool))
	http.HandleFunc("/import-user-cohort", importers.ImportsHandler(eventEmitter, workerPool))
	http.HandleFunc("/import-user-license", importers.ImportsHandler(eventEmitter, workerPool))

	// Bulk Removal Endpoints
	http.HandleFunc("/remove-users-mp", importers.ImportsHandler(eventEmitter, workerPool))
	http.HandleFunc("/revoke-users-license", importers.ImportsHandler(eventEmitter, workerPool))
}
