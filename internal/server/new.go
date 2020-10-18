package server

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/travmatth-org/qaas/internal/afs"
	"github.com/travmatth-org/qaas/internal/api"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/handlers"
	"github.com/travmatth-org/qaas/internal/logger"
)

type listener struct {
	http net.Listener
}

type channel struct {
	// monitor incoming os signals
	signal chan os.Signal
	// monitor outgoing errors thrown by server
	error chan error
	// monitor server startup
	started chan struct{}
}

type timeout struct {
	// duration to wait for server to gracefully shutdown
	stop time.Duration
}

// Server w/ embedded *http.Server, *mux.Router, zerolog.Logger instances.
// Manages configurations, starting, stopping and routing of server instances
type Server struct {
	*mux.Router
	*http.Server
	api      *api.API
	fs       *afs.AFS
	address  string
	static   map[string]string
	config   *config.Config
	listener listener
	channel  channel
	timeout  timeout
}

// Opt is the signature of server configuration functions
type Opt func(s *Server) (*Server, error)

// helper func to evaluate server options
func build(s *Server, opts ...Opt) (*Server, error) {
	var err error
	for _, fn := range opts {
		if s, err = fn(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// New configures and returns a server instance struct with the specified opts
func New(c *config.Config, opts ...Opt) (*Server, error) {
	router := mux.NewRouter()
	s := &Server{
		Router: router,
		Server: &http.Server{
			Addr:         c.Net.IP + c.Net.Port,
			WriteTimeout: time.Duration(c.Timeout.Write) * time.Second,
			ReadTimeout:  time.Duration(c.Timeout.Read) * time.Second,
			IdleTimeout:  time.Duration(c.Timeout.Idle) * time.Second,
			Handler:      router,
		},
		api:     nil,
		fs:      nil,
		address: c.Net.IP + c.Net.Port,
		static:  make(map[string]string),
		config:  c,
		channel: channel{
			signal:  make(chan os.Signal, 1),
			error:   make(chan error, 1),
			started: make(chan struct{}),
		},
		timeout: timeout{
			stop: time.Duration(c.Timeout.Stop) * time.Second,
		},
		listener: listener{
			http: nil,
		},
	}
	return build(s, opts...)
}

// WithAPI inserts the given api client into the server
func WithAPI(a *api.API) Opt {
	return func(s *Server) (*Server, error) {
		s.api = a
		return s, nil
	}
}

// WithFS inserts the given file system client into the server
func WithFS(fs *afs.AFS) Opt {
	return func(s *Server) (*Server, error) {
		s.fs = fs
		return s, nil
	}
}

// WithStatic loads the given static assets into the server
func WithStatic(s *Server) (*Server, error) {
	if err := s.fs.LoadAssets(s.config.Net.Static); err != nil {
		return nil, err
	}
	return s, nil
}

// WithStaticPages loads staic pages as HTTP endpoints
func WithStaticPages(isProd bool) Opt {
	return func(s *Server) (*Server, error) {
		s.HandleFunc("/", s.Route(handlers.Static(s.fs.Use("index")), isProd))
		s.NotFoundHandler = s.Route(handlers.Static(s.fs.Use("404")), isProd)
		logger.Info().Msg("Registered home and 404 html pages to endpoints")
		return s, nil
	}
}
