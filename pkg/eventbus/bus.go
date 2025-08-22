package eventbus

import (
	"context"
	"sync"
)

// Event represents a generic event
type Event struct {
	Name string
	Data interface{}
}

// Handler is a callback for event consumption
type Handler func(event Event)

// Bus is an in-memory event bus
type Bus struct {
	subscribers map[string][]chan Event
	lock        sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	closed      bool
	wg          sync.WaitGroup
}

// New creates a new Bus
func New() *Bus {
	ctx, cancel := context.WithCancel(context.Background())
	return &Bus{
		subscribers: make(map[string][]chan Event),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Subscribe adds a handler to an event
func (b *Bus) Subscribe(eventName string, handler Handler) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.closed {
		return
	}

	ch := make(chan Event, 16)
	b.subscribers[eventName] = append(b.subscribers[eventName], ch)

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for e := range ch {
			handler(e)
		}
	}()
}

func (b *Bus) Publish(event Event) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.closed {
		return
	}

	if chans, ok := b.subscribers[event.Name]; ok {
		for _, ch := range chans {
			ch <- event // blocking until channel can receive
		}
	}
}

// Close gracefully shuts down all subscribers
func (b *Bus) Close() {
	b.lock.Lock()
	if b.closed {
		b.lock.Unlock()
		return
	}
	b.closed = true
	b.lock.Unlock()

	b.cancel()

	// close all channels
	b.lock.RLock()
	for _, chans := range b.subscribers {
		for _, ch := range chans {
			close(ch)
		}
	}
	b.lock.RUnlock()

	// wait for all subscriber goroutines to finish
	b.wg.Wait()
}
