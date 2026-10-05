package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/whi-0404/commerce-ops/order/internal/config"
	"github.com/whi-0404/commerce-ops/order/internal/server"
)

func main() {
	// 1. Load config từ env
	cfg := config.Load()

	// 2. Setup zerolog — ConsoleWriter cho dev, JSON cho production
	var writer io.Writer = os.Stdout
	if os.Getenv("APP_ENV") != "production" {
		writer = zerolog.ConsoleWriter{Out: os.Stdout}
	}
	logger := zerolog.New(writer).With().Timestamp().Str("service", "order").Logger()

	// 3. Setup router
	router := server.NewRouter()

	// 4. Khởi tạo server
	srv := server.New(cfg.HTTP, router)

	// 5. Chạy server trong goroutine riêng
	go func() {
		if err := srv.Start(logger); err != nil {
			logger.Fatal().Err(err).Msg("server error")
		}
	}()

	// 6. Graceful shutdown: chờ OS signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx, logger); err != nil {
		logger.Fatal().Err(err).Msg("shutdown error")
	}
}