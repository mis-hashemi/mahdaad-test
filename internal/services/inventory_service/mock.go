package inventory_service

import (
	"context"
	"fmt"
)

type MockInventoryService struct {
	DeductErr  bool
	RestoreErr bool

	RestoreCalled bool
}

func (m *MockInventoryService) Deduct(ctx context.Context, orderID string) error {
	if m.DeductErr {
		return fmt.Errorf("deduct inventory failed")
	}
	return nil
}

func (m *MockInventoryService) Restore(ctx context.Context, orderID string) error {
	m.RestoreCalled = true
	if m.RestoreErr {
		return fmt.Errorf("restore inventory failed")
	}
	return nil
}
