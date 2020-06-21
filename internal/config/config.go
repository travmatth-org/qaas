package config

import (
	"flag"
	"path/filepath"
	"time"
)

const (
	DefaultRoot = "/srv/www/static"
	DefaultIP   = "0.0.0.0"
	DefaultPort = "80"
	// DefaultReadTimeout covers the time from when the connection is accepted
	// to when the request body is fully read
	DefaultReadTimeout = 10
	// DefaultWriteTimeout normally covers the time from the end of the request
	// header read to the end of the response write
	// (a.k.a. the lifetime of the ServeHTTP)
	DefaultWriteTimeout = 10
	// Default timeout for server to wait for existing connections to close
	DefaultStopTimeout = 10
	// DefaultIdleTimeout limits server-side the amount of time a Keep-Alive
	// connection will be kept idle before being reused.
	DefaultIdleTimeout = 10
	Index              = "index.html"
	NotFound           = "404.html"
	name               = "faas"
)

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

func Build() *Config {
	c := &Config{}
	message := "Directory static assets served from"
	c.static = flag.String("static", DefaultRoot, message)

	message = "ip server should listen on"
	c.ip = flag.String("ip", DefaultIP, message)

	message = "Port server should listen on"
	c.port = flag.String("port", DefaultPort, message)

	message = "Default timeout period for HTTP responses"
	c.readTimeout = flag.Int("read-timeout", DefaultReadTimeout, message)

	message = "Default timeout period for HTTP responses"
	c.writeTimeout = flag.Int("write-timeout", DefaultWriteTimeout, message)

	message = "Default idle period for HTTP responses"
	c.idleTimeout = flag.Int("idle-timeout", DefaultIdleTimeout, message)

	message = "Default timeout for server to wait for existing connections to close"
	c.stopTimeout = flag.Int("stop-timeout", DefaultStopTimeout, message)

	message = "Set execution for development environment"
	c.dev = flag.Bool("dev", false, message)

	flag.Parse()

	return c
}

func (c Config) GetReadTimeout() time.Duration {
	return time.Duration(*c.readTimeout) * time.Second
}

func (c Config) GetWriteTimeout() time.Duration {
	return time.Duration(*c.writeTimeout) * time.Second
}

func (c Config) GetIdleTimeout() time.Duration {
	return time.Duration(*c.idleTimeout) * time.Second
}

func (c Config) GetStopTimeout() time.Duration {
	return time.Duration(*c.stopTimeout) * time.Second
}

func (c Config) GetStaticRoot() string {
	return *c.static
}

func (c Config) GetAddress() string {
	return *c.ip + ":" + *c.port
}

func (c Config) GetIndexHtml() string {
	return filepath.Join(c.GetStaticRoot(), Index)
}

func (c Config) Get404() string {
	return filepath.Join(c.GetStaticRoot(), NotFound)
}

func (c Config) GetName() string {
	return name
}

func (c Config) IsDev() bool {
	return *c.dev
}
