package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Define websocket upgrader
var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	clients ClientList
	sync.RWMutex

	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()
	return m
}

func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	log.Println("new websocket connection opened")

	// Upgrade standar HTTP connection to a websocket connection
	conn, err := websocketUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, m)
	m.addClient(client)

	// Start 2 goroutines to read and write messages
	go client.readMessage()
	go client.writeMessage()
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}

		return nil
	} else {
		return errors.New("there is no such event type")
	}
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessage
}

func SendMessage(event Event, c *Client) error {
	fmt.Println((event))
	return nil
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}
