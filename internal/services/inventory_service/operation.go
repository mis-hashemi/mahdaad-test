package inventory_service

import (
	"context"
	"log/slog"
)

func (s *Service) Deduct(ctx context.Context, orderID string) error {
	s.log.Info("InventoryDeduct called", slog.String("orderID", orderID))
	return nil
}

func (s *Service) Restore(ctx context.Context, orderID string) error {
	s.log.Info("InventoryRestore called", slog.String("orderID", orderID))
	return nil
}
