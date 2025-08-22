package sms_client

import (
	"context"
	"time"
)

type ExternalSMSService struct { /* config, client */
}

func (s *ExternalSMSService) Send(ctx context.Context, phone, message string) error {
	time.Sleep(5 * time.Millisecond)
	return nil
}
