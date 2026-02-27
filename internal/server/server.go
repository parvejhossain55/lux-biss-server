package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/internal/config"
	"github.com/parvej/luxbiss_server/internal/logger"
)

type Server struct {
	httpServer *http.Server
	router     *gin.Engine
	log        *logger.Logger
	cfg        *config.ServerConfig
}

func New(cfg *config.ServerConfig, log *logger.Logger) *Server {
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		httpServer: httpServer,
		router:     router,
		log:        log,
	}
}

func (s *Server) Router() *gin.Engine {
	return s.router
}

func (s *Server) Start() error {
	// Validate TLS configuration
	if s.cfg.TLSEnabled {
		if s.cfg.TLSCertFile == "" || s.cfg.TLSKeyFile == "" {
			return fmt.Errorf("TLS enabled but cert file or key file not specified")
		}
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)

	go func() {
		s.log.Infof("Server starting on %s", s.httpServer.Addr)
		var err error
		if s.cfg.TLSEnabled {
			s.log.Info("HTTPS enabled")
			err = s.httpServer.ListenAndServeTLS(s.cfg.TLSCertFile, s.cfg.TLSKeyFile)
		} else {
			s.log.Warn("HTTP only - TLS disabled")
			err = s.httpServer.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		s.log.Infof("Received signal %s, shutting down gracefully...", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	s.log.Info("Server stopped gracefully")
	return nil
}
