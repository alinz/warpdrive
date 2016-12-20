package warpify

import (
	"sync"
)

// Callback is an interface for calling back the event
type Callback interface {
	Do(kind int, value string)
}

type pubSub interface {
	Subscribe(int, Callback)
	Unsubscribe(int)
	Publish(int, string)
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

func (s *simplePubSub) Publish(kind int, value string) {
	s.ctx.Lock()
	defer s.ctx.Unlock()

	callback, ok := s.callbacks[kind]
	if ok {
		// we are executing the callback inside a goroutine to
		// make sure it does not block the main process
		go func() {
			callback.Do(kind, value)
		}()
	}
}

func newPubSub() pubSub {
	return &simplePubSub{
		callbacks: make(map[int]Callback),
	}
}
