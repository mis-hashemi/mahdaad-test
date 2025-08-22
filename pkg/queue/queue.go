package queue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Message struct {
	Key  string
	Data interface{}
}

type Handler func(Message) error

type SimpleQueue struct {
	ch chan Message
	wg sync.WaitGroup
}

func New(buffer int) *SimpleQueue {
	return &SimpleQueue{ch: make(chan Message, buffer)}
}

func (q *SimpleQueue) Publish(msg Message) {
	q.ch <- msg
}

func (q *SimpleQueue) StartConsumer(ctx context.Context, handler Handler) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-q.ch:
				for attempt := 1; ; attempt++ {
					err := handler(msg)
					if err == nil {
						break
					}
					fmt.Printf("delivery failed for key %s (attempt %d): %v\n", msg.Key, attempt, err)
					time.Sleep(time.Duration(attempt) * time.Second)
				}
			}
		}
	}()
}

func (q *SimpleQueue) Wait() {
	q.wg.Wait()
}
