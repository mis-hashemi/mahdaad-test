package sms_client

import "context"

type SMSClient interface {
	Send(ctx context.Context, phone, message string) error
}
