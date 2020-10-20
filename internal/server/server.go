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

	"github.com/coreos/go-systemd/v22/activation"
	"github.com/coreos/go-systemd/v22/daemon"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/logger"
)

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
			logger.Error().Err(err).Msg("Error notifying systemd of readiness")
		} else if dur, err := s.GetLivenessCheck(); err != nil {
			logger.Error().Err(err).Dur("duration", dur).Msg("Not starting readiness checks")
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
