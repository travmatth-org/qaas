package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Travmatth/faas/internal/config"
	"github.com/Travmatth/faas/internal/handlerfuncs"
	"github.com/Travmatth/faas/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

type Server struct {
	config      *config.Config
	router      *mux.Router
	server      *http.Server
	log         zerolog.Logger
	stopTimeout time.Duration
}

func New(c *config.Config, log zerolog.Logger) *Server {
	router := mux.NewRouter()
	server := &http.Server{
		Addr:         c.GetAddress(),
		WriteTimeout: c.GetWriteTimeout(),
		ReadTimeout:  c.GetReadTimeout(),
		IdleTimeout:  c.GetIdleTimeout(),
		Handler:      router,
	}
	return &Server{c, router, server, log, c.GetStopTimeout()}
}

func (s *Server) middleware() middleware.Middleware {
	return middleware.Chain(
		hlog.NewHandler(s.log),
		hlog.RequestIDHandler("req_id", "Request-Id"),
		hlog.RemoteAddrHandler("ip"),
		hlog.RequestHandler("dest"),
		hlog.RefererHandler("referer"),
	)
}

func (s *Server) RegisterHandlers() {
	mw := s.middleware()
	homeRoute := mw(handlerfuncs.Home(s.config)).ServeHTTP
	notFoundRoute := mw(http.HandlerFunc(handlerfuncs.NotFoundHandler))
	s.router.HandleFunc("/", homeRoute)
	s.router.NotFoundHandler = notFoundRoute
}

func (s *Server) AcceptConnections() {
	errCh, sigCh := make(chan error, 1), make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go s.listenAndServe(errCh)
	l := s.log.Fatal()
	select {
	case err := <-errCh:
		l.Err(err).Msg("Error occurred, shutting down")
		os.Exit(1)
	case sig := <-sigCh:
		ctx, cancel := context.WithTimeout(context.Background(), s.stopTimeout)
		defer cancel()
		err := s.server.Shutdown(ctx)
		if err != nil {
			l = l.Err(err)
		}
		l.Str("signal", sig.String()).Msg("Received signal, shutting down")
		os.Exit(0)
	}
}

func (s *Server) listenAndServe(ch chan error) {
	addr, static := s.config.GetAddress(), s.config.GetStaticRoot()
	s.log.Info().Str("addr", addr).Str("static", static).Msg("Started")
	ch <- s.server.ListenAndServe()
}
