package server

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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
	address  string
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
	return build(&Server{
		Router: router,
		Server: &http.Server{
			Addr:         c.Net.IP + c.Net.Port,
			WriteTimeout: time.Duration(c.Timeout.Write) * time.Second,
			ReadTimeout:  time.Duration(c.Timeout.Read) * time.Second,
			IdleTimeout:  time.Duration(c.Timeout.Idle) * time.Second,
			Handler:      router,
		},
		address: c.Net.IP + c.Net.Port,
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
	}, opts...)
}

// WithHandlers loads HTTP handlers on endpoints
func WithHandlers(routes map[string]http.HandlerFunc, isProd bool) Opt {
	return func(s *Server) (*Server, error) {
		for endpoint, handler := range routes {
			switch endpoint {
			case "/404":
				s.NotFoundHandler = handlers.Route(handler, isProd)
			default:
				s.HandleFunc(endpoint, handlers.Route(handler, isProd))
			}
		}
		logger.Info().Msg("Registered endpoints")
		return s, nil
	}
}
