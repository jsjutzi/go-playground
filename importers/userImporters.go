package importers

import (
	"context"
	"du-service/config"
	"du-service/utils"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

type User struct {
	id string
	firstName string
	lastName string
	email string
}

// type ImportError struct {
// 	email string
// 	errorType string
// 	errorMessage string
// }

func ImportsHandler(eventEmitter *utils.EventEmitter, wp *utils.WorkerPool, sharedClients *config.SharedClients) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		// Establish SSE connection via headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ctx := r.Context() // Get request's context

		// Create a new subscription to the event emitter
		subscription := eventEmitter.Subscribe()
		defer eventEmitter.Unsubscribe(subscription)

        // Goroutine for sending SSE messages to the client
        go func() {
            for {
                select {
                case event := <-subscription:
                    message := fmt.Sprintf("event: %s\ndata: %s\n\n", event.Name, event.Data)
                    _, err := w.Write([]byte(message))
                    if err != nil {
                        return // Stop if there's an error sending SSE
                    }
                    w.(http.Flusher).Flush()
                case <-ctx.Done():
                    return // Handle cancellation
                }
            }
        }()


		// importErrors := make([]ImportError, 0)
		
		csvFilePath := "./importers/benchmark.csv"
		importType := "mp"
		// csvFilePath := r.Header.Get("X-S3-Path")
		// importName := r.Header.Get("X-Import-Name")
		// importType := r.Header.Get("X-Import-Type")
		// uploadId := r.Header.Get("X-Upload-Id")
		// missionPartnerId := r.Header.Get("X-Mission-Partner-Id")
		// groupId := r.Header.Get("X-Group-Id")

		// Access upload record from dynamo
		// Update status to processing

        file, err := os.Open(csvFilePath)

		// Handle error accessing the file
        if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Error opening CSV file"))

			// Update upload status to error
			return
        }

        defer file.Close()

        csvReader := csv.NewReader(file)
        var wg sync.WaitGroup

		// Reading and processing the CSV
		for {
			line, err := csvReader.Read()
			if err != nil {
				if err == csv.ErrFieldCount {
					// Handle expected errors such as wrong field count
					emitSSE(eventEmitter, utils.Event{Name: "File Error", Error: errors.New("invalid CSV format")})
				}
				break // Break on any error (EOF is also an error)
			}

			// This service does not modify the existing CSV, but it's best practice to create a copy of the line
			lineCopy := make([]string, len(line))
            copy(lineCopy, line) // Create a copy of line for the goroutine

			wg.Add(1) // Increment wait group counter for each goroutine

			task := utils.Task{
				Ctx: ctx,
				Func: func(ctx context.Context) utils.Event {
					defer wg.Done()
					user, validateErr := validateUser(lineCopy)
					
					if validateErr != nil {
						return utils.Event{Name: "ValidationError", Data: validateErr.Error()}
					}

					// ProcessUser returns an Event
					return processUser(user, importType)
				},
				Result: make(chan interface{}, 1), // Buffered channel to prevent blocking
			}

			go func(t utils.Task) {
				wp.Submit(t) // Submit task to worker pool
				// Handling results
				select {
                case result := <-t.Result:
					event, ok := result.(utils.Event)

					if ok {
						emitSSE(eventEmitter, event)
					}
                case <-r.Context().Done():
                    return
                }
			}(task)
		}

		wg.Wait() // Ensure all goroutines complete before handler exits

		// Build error report CSV

		eventEmitter.Broadcast(utils.Event{Name: "importComplete", Data: "All users processed"})
    }
}

func validateUser(line []string) (User, error) {
	// Implement user validation logic

	// Validate user data - firstName, lastName, email
	id := line[0]
	firstName := strings.TrimSpace(line[1])
	lastName:= strings.TrimSpace(line[2])
	email := strings.TrimSpace(line[3])

	var err error

	if firstName == "" {
		err = errors.New("first name is required")
	} else if lastName == "" {
		err = errors.New("last name is required")
	} else if email == "" {
		err = errors.New("email is required")
	} else if !utils.ValidateEmail(email) {
		err = errors.New("invalid email")
	}

	return User{
		id: id,
		firstName: firstName,
		lastName: lastName,
		email: email,
	}, err
}

func processUser(user User, importType string) utils.Event{
    // Implement user processing logic

	// Validate user data - firstName, lastName, email
	fmt.Printf("Processing user: %s %s %s %s\n", user.id, user.firstName, user.lastName, user.email)

	// Create user in KC

	// Create user in Dynamo

	// Based on type of import, call specific functions
	return utils.Event{Name: "UserProcessed", Data: user.id} 
}

// TODO: Figure out how to check progress and emit events condtionally
func shouldEmitEvent(user User) bool {
    // Logic to decide if an SSE event should be emitted
    return true
}

func emitSSE(eventEmitter *utils.EventEmitter, data utils.Event) {
    // Implement SSE event emission, consider thread-safety
	eventEmitter.Broadcast(data)
}
