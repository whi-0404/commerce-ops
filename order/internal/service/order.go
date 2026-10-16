package service

import (
	"context"
	"errors"
	"time"
)

type Order struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	TotalAmount int64  `json:"total_amount"`
	Currency    string `json:"currency"`
}

var ErrOrderNotFound = errors.New("order not found")

type OrderService struct{}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) GetOrder(
	ctx context.Context,
	orderID string,
) (*Order, error) {
	switch orderID {
	case "ord_slow":
		if err := waitFor(ctx, 2*time.Second); err != nil {
			return nil, err
		}

	case "ord_timeout":
		if err := waitFor(ctx, 5*time.Second); err != nil {
			return nil, err
		}
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	switch orderID {
	case "ord_123", "ord_slow", "ord_timeout":
		return &Order{
			ID:          orderID,
			Status:      "pending",
			TotalAmount: 250000,
			Currency:    "VND",
		}, nil

	default:
		return nil, ErrOrderNotFound
	}
}
func waitFor(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-timer.C:
		return nil
	}
}
