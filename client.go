package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	clientInterval = 10 * time.Second
	heartBeatInterval   =  (clientInterval * 9) / 10
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager

	// Egress used to avoid concurrent writes on the websocket connection
	egress chan Event // Do not like, fix later
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (c *Client) readMessage() {
	defer func() {
		//cleanup connection
		c.manager.removeClient(c)
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(clientInterval)); err != nil {	
		log.Println(err)
		return
	}
	
	c.connection.SetPongHandler(c.heartBeatHandler)

	for {
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}

			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling message: %v", err)
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			log.Printf("error routing event: %v", err)
		}

		log.Println(messageType)
		log.Println(string(payload))
	}
}

func (c *Client) writeMessage() {
	defer func() {
		//cleanup connection
		c.connection.Close()
	}()

	ticker := time.NewTicker(heartBeatInterval)

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				// The manager closed the channel
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println(err)
				}

				return
			}

			data, err := json.Marshal(message)

			if err != nil {
				log.Println(err)
				return
			}

			err = c.connection.WriteMessage(websocket.TextMessage, data)

			if err != nil {
				log.Println(err)
				return
			}

			log.Println("Message sent")

		case <-ticker.C:
			log.Println("client ping")

			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (c *Client) heartBeatHandler(msg string) error {
	log.Println("heartbeat")
	return c.connection.SetReadDeadline(time.Now().Add(clientInterval))
}
