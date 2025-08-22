package order_flow_service

import (
	"context"
	"github.com/mis-hashemi/mahdaad-test/internal/services/inventory_service"
	"github.com/mis-hashemi/mahdaad-test/internal/services/order_service"
	"github.com/mis-hashemi/mahdaad-test/internal/services/payment_service"
	"github.com/mis-hashemi/mahdaad-test/pkg/logger"
	"testing"
)

func TestProcessOrder_InventoryFail(t *testing.T) {
	ctx := context.Background()
	l := logger.New(nil)

	orderSrv := &order_service.MockOrderService{}
	invSrv := &inventory_service.MockInventoryService{DeductErr: true}
	paySrv := &payment_service.MockPaymentService{}

	of := New(l, orderSrv, invSrv, paySrv)
	err := of.ProcessOrder(ctx, "ORD-123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !orderSrv.CancelCalled {
		t.Error("expected CancelOrder to be called")
	}
	if invSrv.RestoreCalled {
		t.Error("expected RestoreInventory to be called")
	}
	if paySrv.RefundCalled {
		t.Error("payment refund should not be called")
	}
}

func TestProcessOrder_InventoryFail_CancelOrderFail(t *testing.T) {
	ctx := context.Background()
	l := logger.New(nil)

	orderSrv := &order_service.MockOrderService{CancelErr: true}
	invSrv := &inventory_service.MockInventoryService{DeductErr: true}
	paySrv := &payment_service.MockPaymentService{}

	of := New(l, orderSrv, invSrv, paySrv)
	err := of.ProcessOrder(ctx, "ORD-123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !orderSrv.CancelCalled {
		t.Error("expected CancelOrder to be called")
	}
	if invSrv.RestoreCalled {
		t.Error("expected RestoreInventory to be called")
	}
	if paySrv.RefundCalled {
		t.Error("payment refund should not be called")
	}
}

func TestProcessOrder_CreateOrderFail(t *testing.T) {
	ctx := context.Background()
	l := logger.New(nil)

	orderSrv := &order_service.MockOrderService{CreateErr: true}
	invSrv := &inventory_service.MockInventoryService{}
	paySrv := &payment_service.MockPaymentService{}

	of := New(l, orderSrv, invSrv, paySrv)
	err := of.ProcessOrder(ctx, "ORD-123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if orderSrv.CancelCalled {
		t.Error("CancelOrder should not be called")
	}
	if invSrv.RestoreCalled {
		t.Error("RestoreInventory should not be called")
	}
	if paySrv.RefundCalled {
		t.Error("RefundPayment should not be called")
	}
}

func TestProcessOrder_PaymentFail(t *testing.T) {
	ctx := context.Background()
	l := logger.New(nil)

	orderSrv := &order_service.MockOrderService{}
	invSrv := &inventory_service.MockInventoryService{}
	paySrv := &payment_service.MockPaymentService{ChargeErr: true}

	of := New(l, orderSrv, invSrv, paySrv)
	err := of.ProcessOrder(ctx, "ORD-123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !orderSrv.CancelCalled {
		t.Error("expected CancelOrder to be called")
	}
	if !invSrv.RestoreCalled {
		t.Error("expected RestoreInventory to be called")
	}
	if paySrv.RefundCalled {
		t.Error("expected RefundPayment to be called")
	}
}

func TestProcessOrder_PaymentFail_RefundFail(t *testing.T) {
	ctx := context.Background()
	l := logger.New(nil)

	orderSrv := &order_service.MockOrderService{}
	invSrv := &inventory_service.MockInventoryService{}
	paySrv := &payment_service.MockPaymentService{ChargeErr: true, RefundErr: true}

	of := New(l, orderSrv, invSrv, paySrv)
	err := of.ProcessOrder(ctx, "ORD-123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !orderSrv.CancelCalled {
		t.Error("expected CancelOrder to be called")
	}
	if !invSrv.RestoreCalled {
		t.Error("expected RestoreInventory to be called")
	}
	if paySrv.RefundCalled {
		t.Error("expected RefundPayment to be called even if fail")
	}
}
