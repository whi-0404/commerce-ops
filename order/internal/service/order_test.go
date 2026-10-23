package service

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestOrderService_GetOrder(t *testing.T) {
	svc := NewOrderService()

	test := []struct {
		name        string
		orderID     string
		wantErr     error
		wantOrderID string
	}{
		{
			name:        "existing order",
			orderID:     "ord_123",
			wantOrderID: "ord_123",
		},
		{
			name:    "non-existing order",
			orderID: "ord_999",
			wantErr: ErrOrderNotFound,
		},
	}

	for _, item := range test {
		t.Run(item.name, func(t *testing.T) {
			ctx := context.Background()

			order, err := svc.GetOrder(ctx, item.orderID)

			if item.wantErr != nil {
				if errors.Is(err, item.wantErr) {
					t.Fatalf(
						"expected error %v, got %v",
						item.wantErr,
						err,
					)
				}

				if order != nil {
					t.Fatalf("expected nil order, got %+v", order)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if order.ID != item.wantOrderID {
				t.Errorf(
					"expected order ID %q, got %q",
					item.wantOrderID,
					order.ID,
				)
			}
		})
	}
}

func TestOrderService_GetOrder_CanceledContext(t *testing.T) {
	svc := NewOrderService()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	order, err := svc.GetOrder(ctx, "ord_123")

	if !errors.Is(err, context.Canceled) {
		t.Fatalf(
			"expected context.Canceled, got %v",
			err,
		)
	}

	if order != nil {
		t.Fatalf("expected nil order, got %+v", order)
	}
}

func TestOrderService_GetOrder_DeadlineExceeded(t *testing.T) {
	svc := NewOrderService()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		50*time.Millisecond,
	)
	defer cancel()

	start := time.Now()

	order, err := svc.GetOrder(ctx, "ord_timeout")
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf(
			"expected context.DeadlineExceeded, got %v",
			err,
		)
	}

	if order != nil {
		t.Fatalf("expected nil order, got %+v", order)
	}

	if elapsed >= time.Second {
		t.Fatalf(
			"expected early cancellation, took %v",
			elapsed,
		)
	}
}
