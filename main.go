package main

import (
	"du-service/health"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// type User struct {
// 	firstName string
// 	lastName  string
// 	email     string
// }

func main() {
	setupAPI()
	fmt.Printf("Starting server at port 8080\n")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func setupAPI() {
	// manager := NewManager()

	http.HandleFunc("/health", health.HealthCheckHandler)
	http.HandleFunc("/import-users", importsHandler)
}

func importsHandler(w http.ResponseWriter, r *http.Request) {
	eventsHandler(&w, r)
}

type responsedWriteWrapper struct {
	http.ResponseWriter
}

func(rw *responsedWriteWrapper) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
}

func eventsHandler(w *http.ResponseWriter, r *http.Request) {

	rw := responsedWriteWrapper{*w}

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Expose-Headers", "Content-Type")
   
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	 // Simulate sending events (replace this with real data)
		rw.Header().Set("Content-Type", "text/event-stream")
	  
	   // send a random number every 2 seconds
	for {
		rand.Seed(time.Now().UnixNano())
		fmt.Fprintf(rw, "data: %d \n\n", rand.Intn(100))
		flusher := rw.ResponseWriter.(http.Flusher) // Convert rw.ResponseWriter to http.Flusher interface
		flusher.Flush()
		time.Sleep(2 * time.Second)
	}
}

// func testChan() {
// 	var listOfUsers []User

// 	newUser := User{firstName: "Jack", lastName: "Jutzi", email: "jacks email"}
// 	newUser2 := User{firstName: "Bryce", lastName: "Hull", email: "bryces email"}
// 	newUser3 := User{firstName: "Leighton", lastName: "Tidwell", email: "tidwells email"}
// 	newUser4 := User{firstName: "Justin", lastName: "Knight", email: "knights email"}
// 	newUser5 := User{firstName: "Tom", lastName: "Smith", email: "tommys email"}

// 	listOfUsers = append(listOfUsers, newUser)
// 	listOfUsers = append(listOfUsers, newUser2)
// 	listOfUsers = append(listOfUsers, newUser3)
// 	listOfUsers = append(listOfUsers, newUser4)
// 	listOfUsers = append(listOfUsers, newUser5)

// 	var wg sync.WaitGroup
// 	wg.Add(len(listOfUsers))
// 	c := make(chan string)

// 	for index, value := range listOfUsers {
// 		go func(index int, value User) {
// 			defer wg.Done()
// 			c <- value.firstName
// 		}(index, value)
// 	}

// 	go func() {
// 		for message := range c {
// 			println(message)
// 		}
// 	}()

// 	wg.Wait() // Blocks thread until waitgroup counter is 0

// 	close(c)
// }
