// Package server exposes the rate limiter over HTTP.
package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"ratelimiter"
)

// Server owns the limiter and the HTTP server around it.
type Server struct {
	cfg     Config
	logger  *slog.Logger
	limiter *ratelimiter.Limiter
	ready   atomic.Bool
}

// New builds a Server from cfg.
func New(cfg Config, logger *slog.Logger) (*Server, error) {
	limiter, err := ratelimiter.New(cfg.RateLimit, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid rate limit config: %w", err)
	}
	return &Server{cfg: cfg, logger: logger, limiter: limiter}, nil
}

// Run serves until ctx is canceled, then shuts down gracefully.
func (s *Server) Run(ctx context.Context) error {
	httpServer := &http.Server{
		Addr:              ":" + s.cfg.Port,
		Handler:           s.routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	listener, err := net.Listen("tcp", httpServer.Addr)
	if err != nil {
		return err
	}
	s.logger.Info("listening", "port", s.cfg.Port, "rateLimit", fmt.Sprintf("%+v", s.cfg.RateLimit))
	s.ready.Store(true)

	serveErr := make(chan error, 1)
	go func() { serveErr <- httpServer.Serve(listener) }()

	select {
	case err := <-serveErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}

	s.logger.Info("shutting down")
	s.ready.Store(false)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return httpServer.Shutdown(shutdownCtx)
}
