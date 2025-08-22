package payment_service

import (
	"context"
	"fmt"
)

type MockPaymentService struct {
	ChargeErr bool
	RefundErr bool

	RefundCalled bool
}

func (m *MockPaymentService) Charge(ctx context.Context, orderID string) error {
	if m.ChargeErr {
		return fmt.Errorf("charge payment failed")
	}
	return nil
}

func (m *MockPaymentService) Refund(ctx context.Context, orderID string) error {
	m.RefundCalled = true
	if m.RefundErr {
		return fmt.Errorf("refund payment failed")
	}
	return nil
}
