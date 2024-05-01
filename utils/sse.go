package utils

import (
	"sync"
)

type Event struct {
	Name string
	Data string
	Error error
	StreamId string
}


type EventEmitter struct {
	streams map[string]map[chan Event]struct{}
	mu sync.RWMutex // Protect map from concurrent access
	addChan chan subscription // Channel through which subscriptions can be passed, and the channels passed can pass events
	removeChan chan subscription
	broadcastChan chan Event
}

type subscription struct {
	streamId string
	ch chan Event
}

func NewEventEmitter() *EventEmitter {
	emitter := &EventEmitter{
		streams: make(map[string]map[chan Event]struct{}),
		addChan: make(chan subscription),
		removeChan: make(chan subscription),
		broadcastChan: make(chan Event),
	}

	return emitter
}

func (emitter *EventEmitter) Start() {
	for {
		select {
		case sub := <-emitter.addChan:
			emitter.mu.Lock()
			if _, ok := emitter.streams[sub.streamId]; !ok {
				emitter.streams[sub.streamId] = make(map[chan Event]struct{})
			}
			emitter.streams[sub.streamId][sub.ch] = struct{}{}
			emitter.mu.Unlock()

		case sub := <-emitter.removeChan:
			emitter.mu.Lock()

			if subs, ok := emitter.streams[sub.streamId]; ok {
				delete(subs, sub.ch)
				close(sub.ch)
				if len(subs) == 0 {
					delete(emitter.streams, sub.streamId)
				}
			}

			emitter.mu.Unlock()
			
        case event := <-emitter.broadcastChan:
			emitter.mu.RLock()
            if subs, ok := emitter.streams[event.StreamId]; ok {
				for sub := range subs {
					sub <- event
				}
			}
			emitter.mu.RUnlock()
        }
	}
}

func (emitter *EventEmitter) Subscribe(streamId string) chan Event {
    ch := make(chan Event)
    emitter.addChan <- subscription{streamId, ch}
    return ch
}

func (emitter *EventEmitter) Unsubscribe(streamId string, ch chan Event) {
    emitter.removeChan <- subscription{streamId, ch}
}

func (emitter *EventEmitter) Broadcast(event Event) {
    emitter.broadcastChan <- event
}
