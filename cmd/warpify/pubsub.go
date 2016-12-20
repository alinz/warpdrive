package warpify

import (
	"sync"
)

// Event is a struct for passing to callback
type Event struct {
	Kind  int
	Value string
}

// Callback is an interface for calling back the event
type Callback interface {
	Do(event *Event)
}

type pubSub interface {
	Subscribe(int, Callback)
	Unsubscribe(int)
	Publish(*Event)
}

type simplePubSub struct {
	ctx       sync.Mutex
	callbacks map[int]Callback
}

func (s *simplePubSub) Subscribe(int int, callback Callback) {
	s.ctx.Lock()
	defer s.ctx.Unlock()
	s.callbacks[int] = callback
}

func (s *simplePubSub) Unsubscribe(int int) {
	s.ctx.Lock()
	defer s.ctx.Unlock()
	delete(s.callbacks, int)
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
		callbacks: make(map[int]Callback),
	}
}

func createEvent(kind int, value string) *Event {
	return &Event{
		Kind:  kind,
		Value: value,
	}
}
