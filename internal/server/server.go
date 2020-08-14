package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Travmatth/faas/internal/config"
	"github.com/Travmatth/faas/internal/logger"
	"github.com/Travmatth/faas/internal/middleware"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/coreos/go-systemd/v22/activation"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"
)

// Server represents the running server with embedded
// *http.Server, *mux.Router, zerolog.Logger instances
// manages the configurations, starting, stopping and routing
// of HTTP server instance
type Server struct {
	*config.Config
	*mux.Router
	*http.Server
	stopTimeout    time.Duration
	static         map[string][]byte
	signalChannel  chan os.Signal
	errorChannel   chan error
	startedChannel chan struct{}
	httpListener   net.Listener
}

// New configures and returns a server instance struct.
// Accepts a config and zerolog.Logger struct for embedding
func New(c *config.Config) *Server {
	router := mux.NewRouter()
	server := &http.Server{
		Addr:         c.GetAddress(),
		WriteTimeout: c.GetWriteTimeout(),
		ReadTimeout:  c.GetReadTimeout(),
		IdleTimeout:  c.GetIdleTimeout(),
		Handler:      router,
	}
	m := make(map[string][]byte)
	sig := make(chan os.Signal, 1)
	err := make(chan error, 1)
	started := make(chan struct{})
	return &Server{c, router, server, c.GetStopTimeout(), m, sig, err, started, nil}
}

// WrapRoute composes endpoints by wrapping destination handler with handler
// pipeline providing tracing with aws x-ray, injecting logging middleware
// with request details into the context, and error recovery middleware
func (s *Server) WrapRoute(h http.HandlerFunc) http.HandlerFunc {
	return xray.Handler(
		xray.NewFixedSegmentNamer("faas-httpd"),
		alice.New(
			s.RecoverHandler,
			hlog.NewHandler(*logger.GetLogger()),
			hlog.RequestIDHandler("req_id", "Request-Id"),
			hlog.RemoteAddrHandler("ip"),
			hlog.RequestHandler("dest"),
			hlog.RefererHandler("referer"),
			middleware.Log,
		).ThenFunc(h)).ServeHTTP
}

// RegisterHandlers attemtps to prepare and register the specified
// routes with the given middlewware on the server instance.
func (s *Server) RegisterHandlers() {
	index, notfound := s.GetIndexHTML(), s.Get404()
	s.HandleFunc("/", s.WrapRoute(s.ServeStatic(index)))
	s.NotFoundHandler = s.WrapRoute(s.ServeStatic(notfound))
	logger.Info().Msg("Registered home and 404 html pages to endpoints")
}

// OpenListener returns a listener for the server to receive traffic on, or err
func (s *Server) OpenListener() (net.Listener, error) {
	// when systemd starts a process using socket-based activation it sets
	// `LISTEN_PID` & `LISTEN_FDS`. To check if socket based activation is
	// check to see if they are set
	if os.Getenv("LISTEN_PID") == strconv.Itoa(os.Getpid()) {
		logger.Info().Msg("Activating systemd socket")
		listeners, err := activation.Listeners()
		if err != nil {
			return listeners[0], err
		} else if n := len(listeners); n != 1 {
			err = fmt.Errorf("Systemd socket error: unexepected number of listeners: %d", n)
		}
		return listeners[0], err
	}
	logger.Info().Msg("Activating non-systemd socket")
	return net.Listen("tcp", s.Port)
}

// AcceptConnections listens on the configured address and ports for http
// traffic. Simultaneously listens for incoming os signals, will return on
// either a server error or a shutdown signal
func (s *Server) AcceptConnections() error {
	// register and intercept shutdown signals
	signal.Notify(s.signalChannel, os.Interrupt)

	ln, err := s.OpenListener()
	if err != nil {
		logger.Error().Err(err).Msg("Error initializing listener for http server")
		return err
	}
	s.httpListener = ln
	close(s.startedChannel)

	// process incoming requests, close on err or force shutdown on signal
	go s.startServing()
	select {
	case err := <-s.errorChannel:
		logger.Error().Err(err).Msg("Error occurred, shutting down")
		return err
	case sig := <-s.signalChannel:
		ctx, cancel := context.WithTimeout(context.Background(), s.stopTimeout)
		defer cancel()
		shutdownErr := s.Shutdown(ctx)
		logger.Error().
			Err(shutdownErr).
			Str("signal", sig.String()).
			Msg("Received signal, shutting down")
		return shutdownErr
	}
}

// serve http on given listener, or return if no listener
func (s *Server) startServing() {
	if s.httpListener == nil {
		s.errorChannel <- errors.New("Not listening on port")
		return
	}
	addr, dir := s.GetAddress(), s.Static
	logger.Info().Str("addr", addr).Str("static", dir).Msg("Started")
	// drop permissions before serving
	_ = syscall.Umask(0022)
	s.errorChannel <- s.Serve(*s.httpListener)
}
