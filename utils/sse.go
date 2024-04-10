package utils

import (
	"net/http"
)

// Utility function to send Server-Sent Events to subscribers
func SendSSE(w http.ResponseWriter, event string, data string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	message := "event: " + event + "\n"
	message += "data: " + data + "\n\n"

	w.Write([]byte("event: " + event + "\n"))
	w.Write([]byte("data: " + data + "\n\n"))
	w.(http.Flusher).Flush()
}