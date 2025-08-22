package timeout

import (
	"context"
	"time"
)

type Timeout struct {
	duration time.Duration
	fallback func() (any, error)
}

type Option func(*Timeout)

func WithDuration(d time.Duration) Option {
	return func(t *Timeout) {
		t.duration = d
	}
}

func WithFallback(f func() (any, error)) Option {
	return func(t *Timeout) {
		t.fallback = f
	}
}

func New(opts ...Option) *Timeout {
	t := &Timeout{
		duration: 5 * time.Second,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Execute runs the function with timeout
func (t *Timeout) Execute(fn func(ctx context.Context) (any, error)) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t.duration)
	defer cancel()

	ch := make(chan struct {
		res any
		err error
	}, 1)

	go func() {
		res, err := fn(ctx)
		ch <- struct {
			res any
			err error
		}{res, err}
	}()

	select {
	case <-ctx.Done():
		if t.fallback != nil {
			return t.fallback()
		}
		return nil, ctx.Err()
	case result := <-ch:
		return result.res, result.err
	}
}
