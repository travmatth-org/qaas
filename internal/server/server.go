package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Travmatth/faas/internal/config"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
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
	log         zerolog.Logger
	stopTimeout time.Duration
	static      map[string][]byte
}

// New configures and returns a server instance struct.
// Accepts a config and zerolog.Logger struct for embedding
func New(c *config.Config, log zerolog.Logger) *Server {
	router := mux.NewRouter()
	server := &http.Server{
		Addr:         c.GetAddress(),
		WriteTimeout: c.GetWriteTimeout(),
		ReadTimeout:  c.GetReadTimeout(),
		IdleTimeout:  c.GetIdleTimeout(),
		Handler:      router,
	}
	m := make(map[string][]byte)
	return &Server{c, router, server, log, c.GetStopTimeout(), m}
}

func (s *Server) configureMiddleware() alice.Chain {
	return alice.New(
		hlog.NewHandler(s.log),
		hlog.RequestIDHandler("req_id", "Request-Id"),
		hlog.RemoteAddrHandler("ip"),
		hlog.RequestHandler("dest"),
		hlog.RefererHandler("referer"),
	)
}

// RegisterHandlers attemtps to prepare and register the specified routes with
// the given middlewware on the server instance. Returns error if unable to
// register handlers
func (s *Server) RegisterHandlers() error {
	mw := s.configureMiddleware()

	// register index.html
	index := s.GetIndexHTML()
	err := s.loadFileIntoMemory(index, index)
	if err != nil {
		return err
	}
	endpoint := "/"
	s.HandleFunc(endpoint, mw.ThenFunc(s.ServeStatic(index)).ServeHTTP)
	msg := "Registered static file to endpoint"
	s.log.Info().Str("endpoint", endpoint).Str("file", index).Msg(msg)

	// register 404.html
	notfound := s.Get404()
	err = s.loadFileIntoMemory(notfound, notfound)
	if err != nil {
		return err
	}
	s.NotFoundHandler = mw.ThenFunc(s.ServeStatic(notfound))
	msg = "Registered static file to 404 endpoint"
	s.log.Info().Str("file", notfound).Msg(msg)
	return nil
}

// AcceptConnections listens on the configured address and ports for http
// traffic. Simultaneously listens for incoming os signals, will return on
func (s *Server) AcceptConnections() int {
	errCh, sigCh := make(chan error, 1), make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		s.log.Info().
			Str("addr", s.GetAddress()).
			Str("static", s.GetStaticRoot()).
			Msg("Startin")
		errCh <- s.ListenAndServe()
	}()
	select {
	case err := <-errCh:
		s.log.Fatal().Err(err).Msg("Error occurred, shutting down")
		return fail
	case sig := <-sigCh:
		ctx, cancel := context.WithTimeout(context.Background(), s.stopTimeout)
		defer cancel()
		err := s.Shutdown(ctx)
		msg := "Received signal, shutting down"
		s.log.Fatal().Err(err).Str("signal", sig.String()).Msg(msg)
		return ok
	}
}
