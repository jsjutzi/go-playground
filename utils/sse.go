package utils

import (
	"sync"
)

type Event struct {
	Name string
	Data string
	Error error
}

type EventEmitter struct {
	subscriptions map[chan Event]struct{} // map of subscribers to events
	mu sync.RWMutex // Protect map from concurrent access
	addChan chan chan Event // Channel through which channels can be passed, and the channels passed can pass events
	removeChan chan chan Event
	broadcastChan chan Event
}

func NewEventEmitter() *EventEmitter {
	emitter := &EventEmitter{
		subscriptions: make(map[chan Event]struct{}),
		addChan: make(chan chan Event),
		removeChan: make(chan chan Event),
		broadcastChan: make(chan Event),
	}

	return emitter
}

func (emitter *EventEmitter) Start() {
	for {
		select {
		case subscription := <-emitter.addChan:
			emitter.mu.Lock()
			emitter.subscriptions[subscription] = struct{}{}
			emitter.mu.Unlock()

		case subscription := <-emitter.removeChan:
			emitter.mu.Lock()
            delete(emitter.subscriptions, subscription)
            close(subscription)
			emitter.mu.Unlock()
			
        case event := <-emitter.broadcastChan:
			emitter.mu.RLock()
            for sub := range emitter.subscriptions {
                sub <- event
            }
			emitter.mu.RUnlock()
        }
	}
}

func (emitter *EventEmitter) Subscribe() chan Event {
    subscription := make(chan Event)
    emitter.addChan <- subscription
    return subscription
}

func (emitter *EventEmitter) Unsubscribe(subscription chan Event) {
    emitter.removeChan <- subscription
}

func (emitter *EventEmitter) Broadcast(event Event) {
    emitter.broadcastChan <- event
}



// Utility function to emit Server-Sent Events to subscribers
// TODO: Implement wrapper so handlers don't have to implement this themselves?
// func EventHandler(emitter *EventEmitter) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//         // Set headers for SSE
//         w.Header().Set("Content-Type", "text/event-stream")
//         w.Header().Set("Cache-Control", "no-cache")
//         w.Header().Set("Connection", "keep-alive")

//         // Subscribe to events
//         subscription := emitter.Subscribe()
//         defer emitter.Unsubscribe(subscription)

//         // Write SSE events to the client
//         for event := range subscription {
//             // Format the SSE event
//             message := "event: " + event.Name + "\n"
//             message += "data: " + event.Data + "\n\n"

//             // Write the message to the client
//             w.Write([]byte(message))
//             w.(http.Flusher).Flush()
//         }
//     }
// }