package order_flow_service

import (
	"context"
	"github.com/mis-hashemi/mahdaad-test/pkg/logger"
)

type OrderService interface {
	CreateOrder(ctx context.Context, orderID string) error
	CancelOrder(ctx context.Context, orderID string) error
}

type InventoryService interface {
	Deduct(ctx context.Context, orderID string) error
	Restore(ctx context.Context, orderID string) error
}

type PaymentService interface {
	Charge(ctx context.Context, orderID string) error
	Refund(ctx context.Context, orderID string) error
}

type Service struct {
	log       *logger.Logger
	order     OrderService
	inventory InventoryService
	payment   PaymentService
}

func New(log *logger.Logger, o OrderService, i InventoryService, p PaymentService) *Service {
	return &Service{log, o, i, p}
}
