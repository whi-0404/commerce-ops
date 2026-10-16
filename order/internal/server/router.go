package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/whi-0404/commerce-ops/order/internal/handler"
	"github.com/whi-0404/commerce-ops/order/internal/middleware"
	"github.com/whi-0404/commerce-ops/order/internal/service"
)

func NewRouter(logger *slog.Logger, requestTimeout time.Duration) http.Handler {
	mux := http.NewServeMux()

	healthHandler := handler.NewHealthHandler()

	mux.HandleFunc("GET /healthz", healthHandler.Healthz)
	mux.HandleFunc("GET /readyz", healthHandler.Readyz)

	// Order endpoints.
	orderService := service.NewOrderService()
	orderHandler := handler.NewOrderHandler(orderService, logger)

	mux.HandleFunc("GET /api/v1/orders/{id}", orderHandler.GetOrder)

	var h http.Handler = mux

	// Middleware dược bao bọc từ trong ra ngoài
	// 1. Timeout là lớp gần router nhất.
	h = middleware.RequestTimeout(requestTimeout)(h)

	// 2. Logging bao quanh timeout và router.
	h = middleware.Logging(logger)(h)

	// 3. Request ID chạy đầu tiên.
	h = middleware.RequestID(h)

	return h
}
