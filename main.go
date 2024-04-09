package main

import (
	"du-service/health"
	"fmt"
	"log"
	"net/http"
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
	manager := NewManager()

	http.HandleFunc("/health", health.HealthCheckHandler)
	http.HandleFunc("/ws", manager.serveWS)
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
