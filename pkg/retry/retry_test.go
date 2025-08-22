package retry_test

import (
	"context"
	"errors"
	"github.com/mis-hashemi/mahdaad-test/pkg/retry"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry_SuccessImmediately(t *testing.T) {
	r := retry.New(
		retry.WithConstantBackoff(10 * time.Millisecond),
	)

	called := 0
	err := r.Do(context.Background(), func() error {
		called++
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, called, "should succeed in first attempt")
}

func TestRetry_FailThenSuccess(t *testing.T) {
	r := retry.New(
		retry.WithConstantBackoff(10 * time.Millisecond),
	)

	called := 0
	err := r.Do(context.Background(), func() error {
		called++
		if called < 3 {
			return errors.New("temporary failure")
		}
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, called, "should retry until 3rd attempt succeeds")
}

func TestRetry_AlwaysFail(t *testing.T) {
	exp := retry.New(
		retry.WithExponentialBackoff(10*time.Millisecond, 50*time.Millisecond),
	)

	called := 0
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := exp.Do(ctx, func() error {
		called++
		return errors.New("always fail")
	})

	assert.Error(t, err)
	assert.GreaterOrEqual(t, called, 2, "should retry at least twice")
}

func TestRetry_ContextCanceled(t *testing.T) {
	r := retry.New(
		retry.WithConstantBackoff(50 * time.Millisecond),
	)

	called := 0
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := r.Do(ctx, func() error {
		called++
		return errors.New("fail")
	})

	assert.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, 0, called, "should not even attempt if context is canceled")
}
