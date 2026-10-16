package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/whi-0404/commerce-ops/order/internal/config"
	"github.com/whi-0404/commerce-ops/order/internal/logger"
	"github.com/whi-0404/commerce-ops/order/internal/server"
)

func main() {
	// 1. Load config từ env
	cfg := config.Load()

	// 2. Setup slog logger
	appLogger := logger.New()

	// 3. Setup router
	router := server.NewRouter(appLogger, cfg.HTTP.RequestTimeout)

	// 4. Khởi tạo http server
	srv := server.New(cfg.HTTP, router)

	// 5. Chạy server trong goroutine riêng
	go func() {
		if err := srv.Start(appLogger); err != nil {
			appLogger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	// 6. Graceful shutdown: chờ OS signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx, appLogger); err != nil {
		appLogger.Error("shutdown error", "err", err)
		os.Exit(1)
	}
}
