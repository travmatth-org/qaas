package config

import (
	"flag"
	"path/filepath"
	"time"
)

const (
	defaultRoot         = "/srv/www/static"
	defaultIP           = "0.0.0.0"
	defaultPort         = "80"
	defaultReadTimeout  = 10
	defaultWriteTimeout = 10
	defaultStopTimeout  = 10
	defaultIdleTimeout  = 10
	index               = "index.html"
	notFound            = "404.html"
	name                = "faas"
)

// Config manages the configuration options of the program.
// All members are unexported, accessed solely through member methods
type Config struct {
	static       *string
	ip           *string
	port         *string
	readTimeout  *int
	writeTimeout *int
	stopTimeout  *int
	idleTimeout  *int
	dev          *bool
}

// Build uses `flag` package to build and return config struct.
func Build() *Config {
	c := &Config{}

	message := "Directory static assets served from"
	c.static = flag.String("static", defaultRoot, message)
	message = "ip server should listen on"
	c.ip = flag.String("ip", defaultIP, message)
	message = "Port server should listen on"
	c.port = flag.String("port", defaultPort, message)
	message = "Default timeout period for HTTP responses"
	c.readTimeout = flag.Int("read-timeout", defaultReadTimeout, message)
	message = "Default timeout period for HTTP responses"
	c.writeTimeout = flag.Int("write-timeout", defaultWriteTimeout, message)
	message = "Default idle period for HTTP responses"
	c.idleTimeout = flag.Int("idle-timeout", defaultIdleTimeout, message)
	message = "Default timeout for server to wait for existing connections to close"
	c.stopTimeout = flag.Int("stop-timeout", defaultStopTimeout, message)
	message = "Set execution for development environment"
	c.dev = flag.Bool("dev", false, message)

	flag.Parse()
	return c
}

// GetReadTimeout returns the time.Duration of the read timeout
func (c Config) GetReadTimeout() time.Duration {
	return time.Duration(*c.readTimeout) * time.Second
}

// GetWriteTimeout returns the time.Duration of the write timeout
func (c Config) GetWriteTimeout() time.Duration {
	return time.Duration(*c.writeTimeout) * time.Second
}

// GetIdleTimeout returns the time.Duration of the idle timeout
func (c Config) GetIdleTimeout() time.Duration {
	return time.Duration(*c.idleTimeout) * time.Second
}

// GetStopTimeout returns the time.Duration of the stop timeout
func (c Config) GetStopTimeout() time.Duration {
	return time.Duration(*c.stopTimeout) * time.Second
}

// GetStaticRoot returns the directory of static assets
func (c Config) GetStaticRoot() string {
	return *c.static
}

// GetAddress returns the address:port of the server and port to listen on
func (c Config) GetAddress() string {
	return *c.ip + ":" + *c.port
}

// GetIndexHTML returns the filename of the html page
func (c Config) GetIndexHTML() string {
	return filepath.Join(c.GetStaticRoot(), index)
}

// Get404 returns the filename of the 404 page
func (c Config) Get404() string {
	return filepath.Join(c.GetStaticRoot(), notFound)
}

// GetName returns the name of the program
func (c Config) GetName() string {
	return name
}

// IsDev returns bool representing whether program executing in dev mode
func (c Config) IsDev() bool {
	return *c.dev
}
