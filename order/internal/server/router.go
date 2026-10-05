package server

import (
	"net/http"

	"github.com/whi-0404/commerce-ops/order/internal/handler"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	healthHandler := handler.NewHealthHandler()

	mux.HandleFunc("GET /healthz", healthHandler.Healthz)
	mux.HandleFunc("GET /readyz", healthHandler.Readyz)

	return mux
}
