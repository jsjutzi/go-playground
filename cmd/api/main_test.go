package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to the handler function
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function directly and pass in the ResponseRecorder and Request.
	healthCheckHandler(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	wantStatus := "healthy"

	// Decode the response body into a map
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("error decoding response body: %v", err)
	}

	// Check if the status field matches the expected value
	if gotStatus := response["status"]; gotStatus != wantStatus {
		t.Errorf("handler returned unexpected status: got %v want %v",
			gotStatus, wantStatus)
	}
}
