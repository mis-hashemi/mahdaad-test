package sms_service

import (
	"context"
	"errors"
	"github.com/mis-hashemi/mahdaad-test/adapter/sms_client"
	"github.com/mis-hashemi/mahdaad-test/adapter/sms_client/test_client"
	"github.com/mis-hashemi/mahdaad-test/pkg/circuitbreaker"
	"github.com/mis-hashemi/mahdaad-test/pkg/retry"
	"github.com/mis-hashemi/mahdaad-test/pkg/richerror"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func newTestHandler(client sms_client.SMSClient) *Service {
	cb := circuitbreaker.NewCircuitBreaker(
		circuitbreaker.WithName("sms_service"),
		circuitbreaker.WithFailureThreshold(3),
	)
	r := retry.New(retry.WithExponentialBackoff(10*time.Millisecond, 50*time.Millisecond),
		retry.WithMaxRetries(10),
	)
	return New(client, cb, r)
}

func TestSMSHandler_Success(t *testing.T) {
	client, _ := test_client.New()
	handler := newTestHandler(client)

	ctx := context.Background()
	err := handler.SendSMS(ctx, "+123", "Hello")
	assert.NoError(t, err)
}

func TestSMSHandler_AlwaysFail(t *testing.T) {
	client, _ := test_client.NewFailSMSClient()
	handler := newTestHandler(client)

	ctx := context.Background()
	err := handler.SendSMS(ctx, "+123", "Hello")
	assert.Error(t, err)
}

func TestSMSHandler_Timeout(t *testing.T) {
	client, _ := test_client.NewSlowSMSClient()
	handler := newTestHandler(client)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := handler.SendSMS(ctx, "+123", "Hello")
	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestSMSHandler_CircuitBreaker(t *testing.T) {
	client, _ := test_client.NewFailSMSClient()
	handler := newTestHandler(client)

	ctx := context.Background()
	for i := 0; i < 5; i++ {
		_ = handler.SendSMS(ctx, "+123", "Hello")
	}

	err := handler.SendSMS(ctx, "+123", "Hello")
	assert.Error(t, err)
	assert.Equal(t, circuitbreaker.ErrCircuitOpen, err)
}

func TestSMSHandler_RetrySuccessAfterFail(t *testing.T) {
	// client fails first 2 calls, then succeeds
	client := test_client.NewFlakySMSClient(2)
	handler := newTestHandler(client)

	ctx := context.Background()
	err := handler.SendSMS(ctx, "+123", "Hello")
	assert.NoError(t, err)
}

func TestSMSHandler_ContextCanceled(t *testing.T) {
	client, _ := test_client.NewSlowSMSClient()
	handler := newTestHandler(client)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	err := handler.SendSMS(ctx, "+123", "Hello")
	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	var richErr *richerror.RichError
	if errors.As(err, &richErr) {
		assert.Equal(t, richerror.KindUnexpected, richErr.Kind())
	}
}
