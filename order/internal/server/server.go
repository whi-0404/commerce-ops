package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
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

func (s *Server) Start(logger zerolog.Logger) error {
	logger.Info().Str("addr", s.httpServer.Addr).Msg("HTTP server starting")

	err := s.httpServer.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context, logger zerolog.Logger) error {
	logger.Info().Msg("HTTP server shutting down")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	logger.Info().Msg("HTTP server stopped")

	return nil
}
