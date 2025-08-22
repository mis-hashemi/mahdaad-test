package circuitbreaker_test

import (
	"context"
	"errors"
	"github.com/mis-hashemi/mahdaad-test/pkg/circuitbreaker"
	"github.com/sony/gobreaker"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCircuitBreaker_Success(t *testing.T) {
	cb := circuitbreaker.NewCircuitBreaker()

	res, err := cb.Execute(context.Background(), func() (interface{}, error) {
		return 42, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 42, res)
	assert.Equal(t, gobreaker.StateClosed, cb.State())
}

func TestCircuitBreaker_OpenAfterFailures(t *testing.T) {
	cb := circuitbreaker.NewCircuitBreaker(
		circuitbreaker.WithFailureThreshold(2),
		circuitbreaker.WithTimeout(1*time.Second),
	)

	for i := 0; i < 2; i++ {
		_, err := cb.Execute(context.Background(), func() (interface{}, error) {
			return nil, errors.New("fail")
		})
		assert.Error(t, err)
	}

	_, err := cb.Execute(context.Background(), func() (interface{}, error) {
		return nil, nil
	})
	assert.Equal(t, circuitbreaker.ErrCircuitOpen, err)
	assert.Equal(t, gobreaker.StateOpen, cb.State())
}

func TestCircuitBreaker_ContextTimeout(t *testing.T) {
	cb := circuitbreaker.NewCircuitBreaker()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := cb.Execute(ctx, func() (interface{}, error) {
		time.Sleep(50 * time.Millisecond)
		return "ok", nil
	})

	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestCircuitBreaker_StateChangeCallback(t *testing.T) {
	called := false
	cb := circuitbreaker.NewCircuitBreaker(
		circuitbreaker.WithFailureThreshold(1),
		circuitbreaker.WithStateChangeCallback(func(name, from, to string) {
			called = true
		}),
	)

	_, _ = cb.Execute(context.Background(), func() (interface{}, error) {
		return nil, errors.New("fail")
	})

	assert.True(t, called)
}
