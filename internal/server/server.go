package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/coreos/go-systemd/daemon"
	"github.com/coreos/go-systemd/v22/activation"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/middleware"

	"github.com/NYTimes/gziphandler"
)

type listener struct {
	http net.Listener
}

type channel struct {
	signal  chan os.Signal
	error   chan error
	started chan struct{}
}

type timeout struct {
	stop time.Duration
}

// Server represents the running server with embedded
// *http.Server, *mux.Router, zerolog.Logger instances
// manages the configurations, starting, stopping and routing
// of HTTP server instance
type Server struct {
	*mux.Router
	*http.Server
	address      string
	config       *config.Config
	static       map[string][]byte
	listener listener
	channel channel
	timeout timeout
}

// New configures and returns a server instance struct.
func New(c *config.Config) *Server {
	router := mux.NewRouter()
	s := &Server{
		address: c.Net.IP + c.Net.Port,
		config:  c,
		Router:  router,
		Server: &http.Server{
			Addr:         c.Net.IP + c.Net.Port,
			WriteTimeout: time.Duration(c.Timeout.Write) * time.Second,
			ReadTimeout:  time.Duration(c.Timeout.Read) * time.Second,
			IdleTimeout:  time.Duration(c.Timeout.Idle) * time.Second,
			Handler:      router,
		},
		channel: channel{
			signal:  make(chan os.Signal, 1),
			error:   make(chan error, 1),
			started: make(chan struct{}),
		},
		timeout: timeout{
			stop: time.Duration(c.Timeout.Stop) * time.Second,
		},
		static:       make(map[string][]byte),
		listener: listener{
			http: nil,
		},
	}
	index := filepath.Join(s.config.Net.Static, "index.html")
	missing := filepath.Join(s.config.Net.Static, "404.html")
	s.HandleFunc("/", s.WrapRoute(s.ServeStatic(index)))
	s.NotFoundHandler = s.WrapRoute(s.ServeStatic(missing))
	logger.Info().Msg("Registered home and 404 html pages to endpoints")
	return s
}

// WrapRoute composes endpoints by wrapping destination handler with handler
// pipeline providing tracing with aws x-ray, injecting logging middleware
// with request details into the context, and error recovery middleware,
// and gzipping the response
func (s *Server) WrapRoute(h http.HandlerFunc) http.HandlerFunc {
	gzippedHandler := gziphandler.GzipHandler(h).ServeHTTP
	return xray.Handler(
		xray.NewFixedSegmentNamer("qaas-httpd"),
		alice.New(
			s.RecoverHandler,
			hlog.NewHandler(*logger.GetLogger()),
			hlog.RequestIDHandler("req_id", "Request-Id"),
			hlog.RemoteAddrHandler("ip"),
			hlog.RequestHandler("dest"),
			hlog.RefererHandler("referer"),
			middleware.Log,
		).ThenFunc(gzippedHandler)).ServeHTTP
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
	return net.Listen("tcp", s.config.Net.Port)
}

// AcceptConnections listens on the configured address and ports for http
// traffic. Simultaneously listens for incoming os signals, will return on
// either a server error or a shutdown signal
func (s *Server) AcceptConnections() error {
	// register and intercept shutdown signals
	signal.Notify(s.channel.signal, os.Interrupt)

	switch ln, err := s.OpenListener(); {
	case err != nil:
		logger.Error().Err(err).Msg("Error initializing listener for http server")
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
		logger.Error().Err(err).Msg("Error occurred, shutting down")
		return err
	case sig := <-s.channel.signal:
		ctx, cancel := context.WithTimeout(context.Background(), s.timeout.stop)
		defer cancel()
		err, str := s.Shutdown(ctx), sig.String() 
		logger.Error().Err(err).Str("sig", str).Msg("Received signal, shutting down")
		return err
	}
}

func (s *Server) GetLivenessCheck() (time.Duration, error) {
	if s.config.Env == config.Production {
		return time.Duration(s.config.Net.Liveness) * time.Second, nil
	}
	switch interval, err := daemon.SdWatchdogEnabled(false); {
	case err != nil:
		return time.Duration(0), err
	case interval <= 0:
		message := "Liveness Interval must be greater than 0"
		return time.Duration(0), errors.New(message)
	default:
		return interval, nil
	}
}

// LivenessCheck retrieves home page  to verify the liveness of the server,
// then notifies the systemd daemon to pass the check, systemd will restart server fails health check, 
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
			message := "Not starting readiness checks"
			logger.Warn().Err(err).Dur("duration", dur).Msg(message)
		} else {
			go s.LivenessCheck(dur)
		}
	}
	s.channel.error <- s.Serve(s.listener.http)
}
