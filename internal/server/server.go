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

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/coreos/go-systemd/daemon"
	"github.com/coreos/go-systemd/v22/activation"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"
	"github.com/travmatth-org/qaas/internal/api"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/fs"
	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/middleware"

	"github.com/NYTimes/gziphandler"
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
	fs       *fs.FS
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
func WithFS(fs *fs.FS) Opt {
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
		s.HandleFunc("/", s.WrapRoute(s.ServeStatic("index"), isProd))
		s.NotFoundHandler = s.WrapRoute(s.ServeStatic("404"), isProd)
		logger.Info().Msg("Registered home and 404 html pages to endpoints")
		return s, nil
	}
}

// WrapRoute composes endpoints by wrapping destination handler with handler
// pipeline providing tracing with aws x-ray, injecting logging middleware
// with request details into the context, and error recovery middleware,
// and gzipping the response
func (s *Server) WrapRoute(h http.HandlerFunc, isProd bool) http.HandlerFunc {
	handler := alice.New(
		s.RecoverHandler,
		hlog.NewHandler(*logger.GetLogger()),
		hlog.RequestIDHandler("req_id", "Request-Id"),
		hlog.RemoteAddrHandler("ip"),
		hlog.RequestHandler("dest"),
		hlog.RefererHandler("referer"),
		gziphandler.GzipHandler,
		middleware.Log,
	).ThenFunc(h)
	if isProd {
		namer := xray.NewFixedSegmentNamer("qaas-httpd")
		return xray.Handler(namer, handler).ServeHTTP
	}
	return handler.ServeHTTP
}

// OpenListener returns a listener for the server to receive traffic on, or err
// Will prefer using a systemd activated socket if `LISTEN_PID` defined in env
func (s *Server) OpenListener() (net.Listener, error) {
	// when systemd starts a process using socket-based activation it sets
	// `LISTEN_PID` & `LISTEN_FDS`. To check if socket based activation is
	// check to see if they are set
	if os.Getenv("LISTEN_PID") == strconv.Itoa(os.Getpid()) {
		logger.Info().Msg("Activating systemd socket")
		listeners, err := activation.Listeners()
		if err != nil {
			logger.Error().Err(err).Msg("Error Activating systemd socket")
			return listeners[0], err
		} else if n := len(listeners); n != 1 {
			err = fmt.Errorf("Systemd socket err: too many listeners: %d", n)
			logger.Error().Err(err).Msg("Activating non-systemd socket")
		}
		return listeners[0], err
	}
	logger.Info().Msg("Activating non-systemd socket")
	return net.Listen("tcp", s.config.Net.Port)
}

// GetLivenessCheck retrieves the liveness check interval from systemd
// if running in production mode (i.e., with systemd), else the interval
// specified in the configuration
func (s *Server) GetLivenessCheck() (time.Duration, error) {
	if s.config.Env == config.Production {
		return time.Duration(s.config.Net.Liveness) * time.Second, nil
	}
	switch interval, err := daemon.SdWatchdogEnabled(false); {
	case err != nil:
		logger.Error().Err(err).Msg("Error initializing liveness checks")
		return time.Duration(0), err
	case interval <= 0:
		err := errors.New("Liveness Interval must be greater than 0")
		logger.Error().Err(err).Msg("Error initializing liveness checks")
		return time.Duration(0), err
	default:
		return interval, nil
	}
}

// LivenessCheck retrieves home page  to verify the liveness of the server,
// then notifies the systemd daemon to pass the check.
// systemd will restart server on failed health check
func (s *Server) LivenessCheck(interval time.Duration) {
	for {
		_, err := http.Get(s.address)
		if err != nil {
			logger.Error().Err(err).Msg("Liveness check failed")
			return
		}
		_, err = daemon.SdNotify(false, daemon.SdNotifyWatchdog)
		if err != nil {
			logger.Error().Err(err).Msg("Error in systemd health check")
			return
		}
		time.Sleep(interval)
	}
}

// serve http on given listener, or return if no listener
func (s *Server) start() {
	if s.listener.http == nil {
		s.channel.error <- errors.New("Not listening on port")
		return
	}

	static := s.config.Net.Static
	logger.Info().Str("addr", s.address).Str("static", static).Msg("Started")

	// drop permissions before serving
	_ = syscall.Umask(0022)

	// notify systemd daemon server is ready
	if s.config.Env == config.Production {
		if _, err := daemon.SdNotify(false, daemon.SdNotifyReady); err != nil {
			message := "Error notifying systemd of readiness"
			logger.Error().Err(err).Msg(message)
		} else if dur, err := s.GetLivenessCheck(); err != nil {
			logger.Error().Err(err).
				Dur("duration", dur).
				Msg("Not starting readiness checks")
			s.channel.error <- err
			return
		} else {
			go s.LivenessCheck(dur)
		}
	}
	s.channel.error <- s.Serve(s.listener.http)
}

// AcceptConnections listens on the configured address and ports for http
// traffic. Simultaneously listens for incoming os signals, will return on
// either a server error or a shutdown signal
func (s *Server) AcceptConnections() error {
	// register and intercept shutdown signals
	signal.Notify(s.channel.signal, os.Interrupt)

	switch ln, err := s.OpenListener(); {
	case err != nil:
		logger.Error().Err(err).Msg("Error initializing listener")
		return err
	default:
		s.listener.http = ln
		close(s.channel.started)
	}

	// process incoming requests
	go s.start()

	// close on err or force shutdown on signal
	select {
	case err := <-s.channel.error:
		return err
	case sig := <-s.channel.signal:
		logger.Info().Msg("Received signal: " + sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), s.timeout.stop)
		defer cancel()
		return s.Shutdown(ctx)
	}
}
