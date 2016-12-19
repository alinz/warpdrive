package warpify

import (
	"sync"
)

// EventKind is a type to define Event types
type EventKind int

// Event is a struct for passing to callback
type Event struct {
	Kind  EventKind
	Value interface{}
}

// Callback is an interface for calling back the event
type Callback interface {
	Do(event *Event)
}

type pubSub interface {
	Subscribe(EventKind, Callback)
	Unsubscribe(EventKind)
	Publish(*Event)
}

type simplePubSub struct {
	ctx       sync.Mutex
	callbacks map[EventKind]Callback
}

func (s *simplePubSub) Subscribe(eventKind EventKind, callback Callback) {
	s.ctx.Lock()
	defer s.ctx.Unlock()
	s.callbacks[eventKind] = callback
}

func (s *simplePubSub) Unsubscribe(eventKind EventKind) {
	s.ctx.Lock()
	defer s.ctx.Unlock()
	delete(s.callbacks, eventKind)
}

func (s *simplePubSub) Publish(event *Event) {
	s.ctx.Lock()
	defer s.ctx.Unlock()

	callback, ok := s.callbacks[event.Kind]
	if ok {
		// we are executing the callback inside a goroutine to
		// make sure it does not block the main process
		go func() {
			callback.Do(event)
		}()
	}
}

func newPubSub() pubSub {
	return &simplePubSub{
		callbacks: make(map[EventKind]Callback),
	}
}

func createEvent(kind EventKind, value interface{}) *Event {
	return &Event{
		Kind:  kind,
		Value: value,
	}
}
