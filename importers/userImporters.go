package importers

import (
	"du-service/utils"
	"fmt"

	// "fmt"
	"net/http"
	"strings"

	// "time"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"sync"
)

type User struct {
	id string
	firstName string
	lastName string
	email string
}

type ImportError struct {
	email string
	errorType string
	errorMessage string
}

func ImportsHandler(eventEmitter *utils.EventEmitter, wp *utils.WorkerPool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		importErrors := make([]ImportError, 0)

        // Initialize SSE headers...
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		
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

        var wg sync.WaitGroup
        csvReader := csv.NewReader(file)

		// Define channel to receive results from worker pool
		resultChan := make(chan interface{})

        for {
            line, err := csvReader.Read()

			wg.Add(1)

            if err == io.EOF {
                break // End of file
            }

            if err != nil {
                // Handle error: Reading failure
                break
            }

			// Validate user has correct and valid inputs
            user, err := validateUser(line)
			// Increment the wait group counter
			if err != nil {
				// Handle error: Invalid user
				importError := ImportError{
					email: user.email,
					errorType: "Invalid User",
					errorMessage: err.Error(),
				}

				importErrors = append(importErrors, importError)
				wg.Done()

				continue
			} else {
				task := utils.Task{
					Func: func() interface{} {
						// Decrement the wait group counter when the task is done
						defer wg.Done()
		
						processUser(user, importType)
		
						// After processing, if condition met, emit SSE
						if shouldEmitEvent(user) {
							emitSSE(eventEmitter, utils.Event{
								Name: "Progress Update",
								Data: "50%",
							})
						}

						return nil
					},
					Result: resultChan,
				}
				// Submit the user processing task to the worker pool
				wp.Submit(task)

				result := <-resultChan
				fmt.Println(result)
			}
        }

        // Wait for all processing to complete
        wg.Wait()

        // After all users are processed, emit a final SSE message
        emitSSE(eventEmitter, utils.Event{
			Name: "importComplete",
			Data: "All users processed",
		})

		// Build downloadable report and upload to S3

        // Response to client that the process is finished
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("All users processed"))
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

func processUser(user User, importType string) {
    // Implement user processing logic

	// Validate user data - firstName, lastName, email
	fmt.Printf("Processing user: %s %s %s %s\n", user.id, user.firstName, user.lastName, user.email)

	// Create user in KC

	// Create user in Dynamo

	// Based on type of import, call specific functions
}

func shouldEmitEvent(user User) bool {
    // Logic to decide if an SSE event should be emitted
    return true
}

func emitSSE(eventEmitter *utils.EventEmitter, data utils.Event) {
    // Implement SSE event emission, consider thread-safety
	eventEmitter.Broadcast(data)
}
