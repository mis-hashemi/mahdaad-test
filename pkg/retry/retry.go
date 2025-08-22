package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
)

type Retry struct {
	backOff    backoff.BackOff
	maxRetries uint64
}

type Opt func(*Retry)

func WithConstantBackoff(interval time.Duration) Opt {
	return func(r *Retry) {
		r.backOff = backoff.NewConstantBackOff(interval)
	}
}

func WithExponentialBackoff(initial, max time.Duration) Opt {
	return func(r *Retry) {
		exp := backoff.NewExponentialBackOff()
		exp.InitialInterval = initial
		exp.MaxInterval = max
		r.backOff = exp
	}
}

func WithMaxRetries(n uint64) Opt {
	return func(r *Retry) {
		r.maxRetries = n
	}
}

func New(opts ...Opt) *Retry {
	r := &Retry{
		backOff: backoff.NewExponentialBackOff(),
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

func (r *Retry) Do(ctx context.Context, operation func() error) error {
	bo := r.backOff
	if r.maxRetries > 0 {
		bo = backoff.WithMaxRetries(bo, r.maxRetries)
	}
	bo = backoff.WithContext(bo, ctx)

	return backoff.RetryNotify(
		func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return operation()
			}
		},
		bo,
		func(err error, next time.Duration) {
			fmt.Printf("retry after %s, error: %v\n", next, err)
		},
	)
}
