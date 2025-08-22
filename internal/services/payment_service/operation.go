package payment_service

import (
	"context"
	"log/slog"
)

func (s *Service) Charge(ctx context.Context, orderID string) error {
	s.log.Info("Payment processed for:", slog.String("orderID", orderID))
	return nil
}

func (s *Service) Refund(ctx context.Context, orderID string) error {
	s.log.Info("Payment refunded for:", slog.String("orderID", orderID))
	return nil
}
