package config

import (
	"flag"
	"path/filepath"
	"time"
)

const (
	defaultRoot         = "/srv/www/static"
	defaultIP           = "0.0.0.0"
	defaultPort         = ":80"
	defaultReadTimeout  = 5
	defaultWriteTimeout = 5
	defaultStopTimeout  = 5
	defaultIdleTimeout  = 5
	index               = "index.html"
	notFound            = "404.html"
	name                = "faas"
)

// Config manages the configuration options of the program.
// All members are unexported, accessed solely through member methods
type Config struct {
	Static       string
	IP           string
	Port         string
	ReadTimeout  int
	WriteTimeout int
	StopTimeout  int
	IdleTimeout  int
	Dev          bool
}

// New construct and returns a config with default values,
// for use in testing server
func New() *Config {
	return &Config{
		Static:       defaultRoot,
		IP:           defaultIP,
		Port:         defaultPort,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		StopTimeout:  defaultStopTimeout,
		IdleTimeout:  defaultIdleTimeout,
		Dev:          false,
	}
}

// Build uses `flag` package to build and return config struct.
func Build() *Config {
	message := "Directory static assets served from"
	static := flag.String("static", defaultRoot, message)
	message = "ip server should listen on"
	ip := flag.String("ip", defaultIP, message)
	message = "Port server should listen on"
	port := flag.String("port", defaultPort, message)
	message = "Default timeout period for HTTP responses"
	readTimeout := flag.Int("read-timeout", defaultReadTimeout, message)
	message = "Default timeout period for HTTP responses"
	writeTimeout := flag.Int("write-timeout", defaultWriteTimeout, message)
	message = "Default idle period for HTTP responses"
	idleTimeout := flag.Int("idle-timeout", defaultIdleTimeout, message)
	message = "Default timeout for server to wait for existing connections to close"
	stopTimeout := flag.Int("stop-timeout", defaultStopTimeout, message)
	message = "Set execution for development environment"
	dev := flag.Bool("dev", false, message)

	flag.Parse()

	return &Config{
		*static, *ip, *port, *readTimeout, *writeTimeout,
		*stopTimeout, *idleTimeout, *dev,
	}
}

// GetReadTimeout returns the time.Duration of the read timeout
func (c Config) GetReadTimeout() time.Duration {
	return time.Duration(c.ReadTimeout) * time.Second
}

// GetWriteTimeout returns the time.Duration of the write timeout
func (c Config) GetWriteTimeout() time.Duration {
	return time.Duration(c.WriteTimeout) * time.Second
}

// GetIdleTimeout returns the time.Duration of the idle timeout
func (c Config) GetIdleTimeout() time.Duration {
	return time.Duration(c.IdleTimeout) * time.Second
}

// GetStopTimeout returns the time.Duration of the stop timeout
func (c Config) GetStopTimeout() time.Duration {
	return time.Duration(c.StopTimeout) * time.Second
}

// GetAddress returns the address:port of the server and port to listen on
func (c Config) GetAddress() string {
	return c.IP + c.Port
}

// GetIndexHTML returns the filename of the html page
func (c Config) GetIndexHTML() string {
	return filepath.Join(c.Static, index)
}

// Get404 returns the filename of the 404 page
func (c Config) Get404() string {
	return filepath.Join(c.Static, notFound)
}

// IsDev returns bool representing whether program executing in dev mode
func (c Config) IsDev() bool {
	return c.Dev
}
