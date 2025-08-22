package timeout_test

import (
	"context"
	"errors"
	"github.com/mis-hashemi/mahdaad-test/pkg/timeout"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeout_Success(t *testing.T) {
	to := timeout.New(timeout.WithDuration(100 * time.Millisecond))

	res, err := to.Execute(func(ctx context.Context) (any, error) {
		return "ok", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "ok", res)
}

func TestTimeout_TimeoutFallback(t *testing.T) {
	to := timeout.New(
		timeout.WithDuration(10*time.Millisecond),
		timeout.WithFallback(func() (any, error) {
			return "fallback", nil
		}),
	)

	res, err := to.Execute(func(ctx context.Context) (any, error) {
		time.Sleep(50 * time.Millisecond)
		return "slow", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "fallback", res)
}

func TestTimeout_TimeoutError(t *testing.T) {
	to := timeout.New(timeout.WithDuration(10 * time.Millisecond))

	res, err := to.Execute(func(ctx context.Context) (any, error) {
		time.Sleep(50 * time.Millisecond)
		return "slow", nil
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestTimeout_ErrorFromFunc(t *testing.T) {
	to := timeout.New(timeout.WithDuration(50 * time.Millisecond))

	res, err := to.Execute(func(ctx context.Context) (any, error) {
		return nil, errors.New("boom")
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.EqualError(t, err, "boom")
}
