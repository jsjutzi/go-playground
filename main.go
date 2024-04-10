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
        // Set headers for SSE
        w.Header().Set("Content-Type", "text/event-stream")
        w.Header().Set("Cache-Control", "no-cache")
        w.Header().Set("Connection", "keep-alive")

        // Subscribe to events
        sub := eventEmitter.Subscribe()
        defer eventEmitter.Unsubscribe(sub)

        // Write SSE events to the client
        for event := range sub {
            // Format the SSE event
            message := "event: " + event.Name + "\n"
            message += "data: " + event.Data + "\n\n"

            // Write the message to the client
            w.Write([]byte(message))
            w.(http.Flusher).Flush()
        }

        // If this point is reached, it means the SSE subscription has ended
        // Now handle the POST request body
        if r.Method == http.MethodPost {
            // Parse the request body
            err := r.ParseForm()
            if err != nil {
                http.Error(w, "Failed to parse request body", http.StatusBadRequest)
                return
            }

            // Access form values from the request body
            // Syntax r.Form.Get("fieldName")
            // Process the form data accordingly
            // Return a response
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("POST request handled successfully"))
        }
    }
}
