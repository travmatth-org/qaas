package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Travmatth/faas/internal/config"
	"github.com/Travmatth/faas/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

const ()

type Server struct {
	*config.Config
	*mux.Router
	*http.Server
	log         zerolog.Logger
	stopTimeout time.Duration
	static      map[string][]byte
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
	m := make(map[string][]byte)
	return &Server{c, router, server, log, c.GetStopTimeout(), m}
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

func (s *Server) RegisterHandlers() error {
	mw := s.middleware()

	index := s.GetIndexHtml()
	if err := s.LoadFileIntoMemory(index, index); err != nil {
		return err
	}
	endpoint := "/"
	s.HandleFunc(endpoint, mw(s.ServeStatic(index)).ServeHTTP)
	s.log.Info().
		Str("endpoint", endpoint).
		Str("file", index).
		Msg("Registered static file to endpoint")

	notFound := s.Get404()
	if err := s.LoadFileIntoMemory(notFound, notFound); err != nil {
		return err
	}
	handler := mw(s.ServeStatic(notFound)).ServeHTTP
	s.NotFoundHandler = http.HandlerFunc(handler)
	s.log.Info().
		Str("file", notFound).
		Msg("Registered static file to 404 endpoint")

	return nil
}

func (s *Server) AcceptConnections() {
	errCh, sigCh := make(chan error, 1), make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		s.log.Info().
			Str("addr", s.GetAddress()).
			Str("static", s.GetStaticRoot()).
			Msg("Started")
		errCh <- s.ListenAndServe()
	}()
	select {
	case err := <-errCh:
		s.log.Fatal().Err(err).Msg("Error occurred, shutting down")
		os.Exit(1)
	case sig := <-sigCh:
		ctx, cancel := context.WithTimeout(context.Background(), s.stopTimeout)
		defer cancel()
		err := s.Shutdown(ctx)
		s.log.Fatal().
			Err(err).
			Str("signal", sig.String()).
			Msg("Received signal, shutting down")
		os.Exit(0)
	}
}
