package order_service

import (
	"context"
	"log/slog"
)

func (s *Service) CreateOrder(ctx context.Context, orderID string) error {
	return nil
}

func (s *Service) CancelOrder(ctx context.Context, orderID string) error {
	s.log.Info("CancelOrder called", slog.String("orderID", orderID))
	return nil
}
