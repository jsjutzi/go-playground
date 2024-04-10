package main

import (
	"du-service/health"
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
	eventEmitter := utils.NewEventEmitter()

	http.HandleFunc("/health", health.HealthCheckHandler)
	http.HandleFunc("/import-users", importsHandler(eventEmitter))
}

func importsHandler(eventEmitter *utils.EventEmitter) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Subscribe to events
        subscription := eventEmitter.Subscribe()
        defer eventEmitter.Unsubscribe(subscription)

        // Handle other logic here
    }
}
