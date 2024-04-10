package importers

import (
	"du-service/utils"
	"fmt"
	"net/http"
)

func ImportsHandler(eventEmitter *utils.EventEmitter, wp *utils.WorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Create a channel to signal when the handler exits
		// This guarantees graceful exit of the handler
		exit := make(chan struct{})
		defer close(exit)

		// Subscribe to events
		subscription := eventEmitter.Subscribe()
		defer eventEmitter.Unsubscribe(subscription)

		// Test event emitter
		eventEmitter.Broadcast(utils.Event{Name: "worker", Data: "Event from worker"})

		// Write SSE events to the client
		go func() {
			select {
			case event := <-subscription:
				// Format the SSE event
				message := "event: " + event.Name + "\n"
				message += "data: " + event.Data + "\n\n"

				// Write the message to the client
				w.Write([]byte(message))
				w.(http.Flusher).Flush()
			case <-exit:
				return
			}
		}()

		eventEmitter.Broadcast(utils.Event{Name: "worker", Data: "2nd event from worker"})

		if r.Method == http.MethodGet {
			fmt.Println("GET request received")
		
			// Extract custom headers from the request
			customHeader := r.Header.Get("X-Custom-Header")

			if customHeader != "" {
				fmt.Println("Custom header value:", customHeader)
			} else {
				fmt.Println("Custom header not found")
			}
		
			// Process the request further if needed
			// You can access other request properties as needed
			// For syntax:
			// queryString := r.URL.Query().Get("paramName")
		
			// Return a response
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("GET request handled successfully"))
		}

		// Create channel for response
		// responseChan := make(chan string) // TODO: change to event struct

		// Submit a task to the worker pool
		fmt.Println("Submitting task to worker pool")

		wp.Submit(func(respChan chan string) { // Here too
			// Logic to be executed by the worker
			respChan <- "Event from worker"
		})

		// Wait for the response from the worker
		// response := <-responseChan

		// Do things with the response

		// Send status
		w.WriteHeader(http.StatusOK)
	}
}
