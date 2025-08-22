package sms_service

import (
	"context"
	"github.com/mis-hashemi/mahdaad-test/adapter/sms_client"
	"github.com/mis-hashemi/mahdaad-test/pkg/circuitbreaker"
	"github.com/mis-hashemi/mahdaad-test/pkg/retry"
	"github.com/mis-hashemi/mahdaad-test/pkg/richerror"
)

type Service struct {
	client sms_client.SMSClient
	cb     *circuitbreaker.CircuitBreaker
	retry  *retry.Retry
}

func New(client sms_client.SMSClient, cb *circuitbreaker.CircuitBreaker, r *retry.Retry) *Service {
	return &Service{
		client: client,
		cb:     cb,
		retry:  r,
	}
}

// SendSMS sends SMS using retry and circuit breaker
func (h *Service) SendSMS(ctx context.Context, phone, msg string) error {
	return h.retry.Do(ctx, func() error {
		_, err := h.cb.Execute(ctx, func() (interface{}, error) {
			if err := h.client.Send(ctx, phone, msg); err != nil {
				return nil, richerror.New("SMSHandler.SendSMS").
					WithErr(err).
					WithKind(richerror.KindUnexpected)
			}
			return nil, nil
		})
		return err
	})
}
