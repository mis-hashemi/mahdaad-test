package order_flow_service

import (
	"context"
	"github.com/mis-hashemi/mahdaad-test/pkg/retry"
	"github.com/mis-hashemi/mahdaad-test/pkg/saga"
	"log/slog"
	"time"
)

func (srv *Service) ProcessOrder(ctx context.Context, orderID string) error {
	sg := saga.New(srv.log)

	for _, step := range srv.buildOrderFlowSteps(orderID) {
		sg.AddStep(step)
	}

	if err := sg.Execute(ctx); err != nil {
		srv.log.Error("saga failed", slog.Any("order_id", orderID), slog.Any("error", err))
		return err
	}

	srv.log.Info("order completed successfully", slog.String("order_id", orderID))
	return nil
}

func (srv *Service) buildOrderFlowSteps(orderID string) []saga.SagaStep {
	return []saga.SagaStep{
		{
			Name: "CreateOrder",
			Action: func(ctx context.Context) error {
				return srv.order.CreateOrder(ctx, orderID)
			},
			Compensation: func(ctx context.Context) error {
				return srv.order.CancelOrder(ctx, orderID)
			},
		},
		{
			Name: "DeductInventory",
			Action: func(ctx context.Context) error {
				return srv.inventory.Deduct(ctx, orderID)
			},
			Compensation: func(ctx context.Context) error {
				return srv.inventory.Restore(ctx, orderID)
			},
			Retry: retry.New(
				retry.WithExponentialBackoff(500*time.Millisecond, 5*time.Second),
				retry.WithMaxRetries(3),
			),
		},
		{
			Name: "ProcessPayment",
			Action: func(ctx context.Context) error {
				return srv.payment.Charge(ctx, orderID)
			},
			Compensation: func(ctx context.Context) error {
				return srv.payment.Refund(ctx, orderID)
			},
		},
	}
}
