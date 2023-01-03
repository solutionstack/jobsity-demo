package server

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server is a light wrapper around http.Server.
type Server struct {
	*http.Server
	address  string
	listener net.Listener
	// Channel used to signal server has shutdown
	serverShutdown chan struct{}
}

// StartHTTPServer creates a Server and starts it.
func StartHTTPServer(router http.Handler, logger zerolog.Logger) error {
	srv, err := New(&http.Server{
		Handler: router,
	}, logger)
	if err != nil {
		return err
	}

	return srv.Start(logger)
}

// New creates a new Server, built off of a base http.Server.
func New(s *http.Server, logger zerolog.Logger) (*Server, error) {
	// Default server timout, in seconds
	const defaultSrvTimeout = 3 * time.Second
	const port = "8081"

	var (
		srv *Server
		err error
	)

	// ensure timeouts are set
	if s.ReadTimeout == 0 {
		s.ReadTimeout = defaultSrvTimeout
	}

	if s.WriteTimeout == 0 {
		s.WriteTimeout = defaultSrvTimeout
	}

	listener, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", port))
	if err != nil {
		logger.Error().Err(err).Msgf("%s unavailable", port)
	}
	srv = &Server{
		s,
		port,
		listener,
		make(chan struct{}),
	}

	return srv, nil
}

// Start begins serving, and listens for termination signals to shutdown gracefully.
func (srv *Server) Start(logger zerolog.Logger) error {
	var err error

	go srv.shutdown(logger)

	logger.Log().
		Str("address", srv.address).
		Int("pid", os.Getpid()).
		Msg("chat http server listening")

	err = srv.Serve(srv.listener)
	if err != nil && err != http.ErrServerClosed {
		return errors.New("server failed to start")
	}

	<-srv.serverShutdown

	return nil
}

// Shutdown server gracefully on SIGINT or SIGTERM.
func (srv *Server) shutdown(logger zerolog.Logger) {
	// Block until signal is received
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Allow up to thirty seconds for server operations to finish before
	// canceling them.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error().
			Err(err).
			Msg("Server shutdown error")
	}

	logger.Log().Msg("chat htp server shutdown")

	// Close channel to signal shutdown is complete
	close(srv.serverShutdown)
}
