package test_client

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type SMSClient interface {
	Send(ctx context.Context, phone, msg string) error
}

// Normal SMS client (always succeeds)
func New() (SMSClient, error) {
	return &externalSMSService{}, nil
}

type externalSMSService struct{}

func (s *externalSMSService) Send(ctx context.Context, phone, message string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Millisecond):
		fmt.Println("send", phone, message)
		return nil
	}
}

// Always failing client
func NewFailSMSClient() (SMSClient, error) {
	return &failSMSClient{}, nil
}

type failSMSClient struct{}

func (c *failSMSClient) Send(ctx context.Context, phone, msg string) error {
	return errors.New("external service failed")
}

// Slow client (simulate timeout)
func NewSlowSMSClient() (SMSClient, error) {
	return &slowSMSClient{}, nil
}

type slowSMSClient struct{}

func (c *slowSMSClient) Send(ctx context.Context, phone, msg string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(200 * time.Millisecond):
		return nil
	}
}

// Flaky client: fails first N times, then succeeds
func NewFlakySMSClient(failures int) SMSClient {
	return &flakySMSClient{failuresLeft: failures}
}

type flakySMSClient struct {
	failuresLeft int
}

func (c *flakySMSClient) Send(ctx context.Context, phone, msg string) error {
	if c.failuresLeft > 0 {
		c.failuresLeft--
		return errors.New("temporary failure")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Millisecond):
		fmt.Println("send", phone, msg)
		return nil
	}
}
