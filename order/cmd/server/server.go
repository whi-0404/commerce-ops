package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/whi-0404/commerce-ops/order/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg config.HTTPConfig, handler http.Handler) *Server {
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	return &Server{
		httpServer: httpServer,
	}
}

func (s *Server) Start(logger *slog.Logger) error {
	logger.Info("HTTP server starting",
		"addr", s.httpServer.Addr,
	)

	err := s.httpServer.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context, logger *slog.Logger) error {
	logger.Info("HTTP server shutting down")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	logger.Info("HTTP server stopped")

	return nil
}
