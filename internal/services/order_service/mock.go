package order_service

import (
	"context"
	"fmt"
)

// MockOrderService implements OrderService for testing
type MockOrderService struct {
	CreateErr bool
	CancelErr bool

	CancelCalled bool
}

func (m *MockOrderService) CreateOrder(ctx context.Context, orderID string) error {
	if m.CreateErr {
		return fmt.Errorf("create order failed")
	}
	return nil
}

func (m *MockOrderService) CancelOrder(ctx context.Context, orderID string) error {
	m.CancelCalled = true
	if m.CancelErr {
		return fmt.Errorf("cancel order failed")
	}
	return nil
}
