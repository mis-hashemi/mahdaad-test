package circuitbreaker

import (
	"context"
	"errors"
	"log"
	"time"

	gobreaker "github.com/sony/gobreaker"
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type CircuitBreaker struct {
	cb *gobreaker.CircuitBreaker
}

type Opt func(*gobreaker.Settings)

func NewCircuitBreaker(opts ...Opt) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        "default",
		MaxRequests: 1,
		Timeout:     5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf("circuit breaker %s state change: %s -> %s", name, from, to)
		},
	}

	for _, opt := range opts {
		opt(&settings)
	}

	cb := gobreaker.NewCircuitBreaker(settings)
	return &CircuitBreaker{cb: cb}
}

func WithName(name string) Opt {
	return func(s *gobreaker.Settings) {
		s.Name = name
	}
}

func WithFailureThreshold(threshold uint32) Opt {
	return func(s *gobreaker.Settings) {
		s.ReadyToTrip = func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= threshold
		}
	}
}

func WithTimeout(timeout time.Duration) Opt {
	return func(s *gobreaker.Settings) {
		s.Timeout = timeout
	}
}

func WithMaxRequests(max uint32) Opt {
	return func(s *gobreaker.Settings) {
		s.MaxRequests = max
	}
}

func WithInterval(interval time.Duration) Opt {
	return func(s *gobreaker.Settings) {
		s.Interval = interval
	}
}

func WithStateChangeCallback(fn func(name string, from, to string)) Opt {
	return func(s *gobreaker.Settings) {
		s.OnStateChange = func(name string, fromState, toState gobreaker.State) {
			fn(name, fromState.String(), toState.String())
		}
	}
}

func (c *CircuitBreaker) Execute(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	resultChan := make(chan struct {
		res interface{}
		err error
	})

	go func() {
		res, err := c.cb.Execute(fn)
		if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
			resultChan <- struct {
				res interface{}
				err error
			}{nil, ErrCircuitOpen}
			return
		}
		resultChan <- struct {
			res interface{}
			err error
		}{res, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-resultChan:
		return r.res, r.err
	}
}

func (c *CircuitBreaker) State() gobreaker.State {
	return c.cb.State()
}
