package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/whi-0404/commerce-ops/order/internal/middleware"
	"github.com/whi-0404/commerce-ops/order/internal/response"
	"github.com/whi-0404/commerce-ops/order/internal/service"
)

type OrderGetter interface {
	GetOrder(
		ctx context.Context,
		orderID string,
	) (*service.Order, error)
}

type OrderHandler struct {
	service OrderGetter
	logger  *slog.Logger
}

func NewOrderHandler(
	orderService OrderGetter,
	logger *slog.Logger,
) *OrderHandler {
	return &OrderHandler{
		service: orderService,
		logger:  logger,
	}
}

func (h *OrderHandler) GetOrder(
	w http.ResponseWriter,
	r *http.Request,
) {
	requestID := middleware.GetRequestID(r.Context())

	orderID := r.PathValue("id")

	if orderID == "" {
		_ = response.WriteError(
			w,
			http.StatusBadRequest,
			"INVALID_ORDER_ID",
			"Order ID is required",
			requestID,
		)
		return
	}

	order, err := h.service.GetOrder(r.Context(), orderID)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			_ = response.WriteError(
				w,
				http.StatusNotFound,
				"ORDER_NOT_FOUND",
				"Order not found",
				requestID,
			)
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			h.logger.Warn(
				"request deadline exceeded",
				"request_id", requestID,
				"order_id", orderID,
				"error", err,
			)
			_ = response.WriteError(
				w,
				http.StatusServiceUnavailable,
				"REQUEST_DEADLINE_EXCEEDED",
				"Request could not be completed within the allowed time",
				requestID,
			)
			return
		}

		if errors.Is(err, context.Canceled) {
			// Client có thể đã ngắt kết nối.
			// Không cố ghi thêm response nếu request đã bị hủy.
			h.logger.Info(
				"request canceled",
				"request_id", requestID,
				"order_id", orderID,
			)
			return
		}

		h.logger.Error(
			"failed to get order",
			"request_id", requestID,
			"order_id", orderID,
			"error", err,
		)

		_ = response.WriteError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"An unexpected error occurred",
			requestID,
		)
		return
	}

	_ = response.WriteSuccess(
		w,
		http.StatusOK,
		order,
	)
}
