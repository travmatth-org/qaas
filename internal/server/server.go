package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Travmatth/faas/internal/config"
	"github.com/Travmatth/faas/internal/logger"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"
)

const (
	// OK returned by server.AcceptConnections when signal directs shutdown
	ok = iota
	// Error returned by server.AcceptConnections when error forces shutdown
	fail
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
	httpListener   *net.Listener
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

func (s *Server) logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoReq(r).Msg("Received request")
		h.ServeHTTP(w, r)
	})
}

func (s *Server) wrapRoute(h http.HandlerFunc) http.HandlerFunc {
	return alice.New(
		hlog.NewHandler(*logger.GetLogger()),
		hlog.RequestIDHandler("req_id", "Request-Id"),
		hlog.RemoteAddrHandler("ip"),
		hlog.RequestHandler("dest"),
		hlog.RefererHandler("referer"),
		s.logMiddleware,
	).ThenFunc(h).ServeHTTP
}

// RegisterHandlers attemtps to prepare and register the specified routes with
// the given middlewware on the server instance. Returns error if unable to
// register handlers
func (s *Server) RegisterHandlers() error {
	index, notfound := s.GetIndexHTML(), s.Get404()
	// register endpoints
	s.HandleFunc("/", s.wrapRoute(s.ServeStatic(index)))
	// 404 endpoint
	s.NotFoundHandler = s.wrapRoute(s.ServeStatic(notfound))
	logger.Info().Msg("Registered home and 404 html pages to endpoints")
	return nil
}

// AcceptConnections listens on the configured address and ports for http
// traffic. Simultaneously listens for incoming os signals, will return on
// either a server error or a shutdown signal
func (s *Server) AcceptConnections() int {
	// register and intercept shutdown signals
	signal.Notify(s.signalChannel, os.Interrupt)

	// start listener and notify on success
	if ln, err := net.Listen("tcp", s.Port); err != nil {
		logger.Error().Err(err).Msg("Error starting server")
		return fail
	} else {
		s.httpListener = &ln
		close(s.startedChannel)
	}

	// process incoming requests, close on err or force shutdown on signal
	go s.StartServing()
	select {
	case err := <-s.errorChannel:
		logger.Error().Err(err).Msg("Error occurred, shutting down")
		return fail
	case sig := <-s.signalChannel:
		ctx, cancel := context.WithTimeout(context.Background(), s.stopTimeout)
		defer cancel()
		logger.Error().
			Err(s.Shutdown(ctx)).
			Str("signal", sig.String()).
			Msg("Received signal, shutting down")
		return ok
	}
}

func (s *Server) StartServing() {
	if s.httpListener == nil {
		s.errorChannel <- errors.New("Not listening on port")
		return
	}
	addr, dir := s.GetAddress(), s.Static
	logger.Info().Str("addr", addr).Str("static", dir).Msg("Started")
	s.errorChannel <- s.Serve(*s.httpListener)
}
