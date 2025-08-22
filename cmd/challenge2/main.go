package main

import (
	"context"
	"github.com/mis-hashemi/mahdaad-test/internal/services/inventory_service"
	"github.com/mis-hashemi/mahdaad-test/internal/services/order_flow_service"
	"github.com/mis-hashemi/mahdaad-test/internal/services/order_service"
	"github.com/mis-hashemi/mahdaad-test/internal/services/payment_service"
	"github.com/mis-hashemi/mahdaad-test/pkg/logger"
	"log/slog"
)

func main() {
	ctx := context.Background()
	l := logger.New(nil)
	orderSrv := order_service.New(l)
	inventorySrv := inventory_service.New(l)
	paymentSrv := payment_service.New(l)
	of := order_flow_service.New(l, orderSrv, inventorySrv, paymentSrv)

	err := of.ProcessOrder(ctx, "ORD-123")
	if err != nil {
		l.Info("Transaction failed.", slog.Any("error", err))
	} else {
		l.Info("Transaction succeeded")
	}
}
