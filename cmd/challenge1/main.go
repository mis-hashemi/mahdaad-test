package main

import (
	"context"
	"fmt"
	"github.com/mis-hashemi/mahdaad-test/adapter/sms_client/test_client"
	"github.com/mis-hashemi/mahdaad-test/internal/sms_service"
	"github.com/mis-hashemi/mahdaad-test/pkg/circuitbreaker"
	"github.com/mis-hashemi/mahdaad-test/pkg/retry"
	"time"
)

func main() {
	cb := circuitbreaker.NewCircuitBreaker(
		circuitbreaker.WithName("SMS CB"),
		circuitbreaker.WithFailureThreshold(3),
		circuitbreaker.WithTimeout(3*time.Second),
	)

	r := retry.New(
		retry.WithExponentialBackoff(10*time.Millisecond, 50*time.Millisecond),
	)

	client, _ := test_client.New()

	handler := sms_service.New(client, cb, r)

	ctx := context.Background()
	phone := "+989135718862"
	msg := "Hello from demo!"

	err := handler.SendSMS(ctx, phone, msg)
	if err != nil {
		fmt.Println("SendSMS error:", err)
	}

	fmt.Println("CircuitBreaker state:", cb.State())
}
