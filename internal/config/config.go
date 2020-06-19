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
	name               = "faas"
)

type Config struct {
	root         string
	ip           string
	port         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	stopTimeout  time.Duration
	idleTimeout  time.Duration
}

func Build() *Config {
	message := "Directory static assets served from"
	static := flag.String("static", DefaultRoot, message)
	message = "ip server should listen on"
	ip := flag.String("ip", DefaultIP, message)
	message = "Port server should listen on"
	port := flag.String("port", DefaultPort, message)
	message = "Default timeout period for HTTP responses"
	readTimeout := flag.Int("read-timeout", DefaultReadTimeout, message)
	message = "Default timeout period for HTTP responses"
	writeTimeout := flag.Int("write-timeout", DefaultWriteTimeout, message)
	message = "Default idle period for HTTP responses"
	idleTimeout := flag.Int("idle-timeout", DefaultIdleTimeout, message)
	message = "Default timeout for server to wait for existing connections to close"
	stopTimeout := flag.Int("stop-timeout", DefaultStopTimeout, message)
	flag.Parse()

	return &Config{
		*static,
		*ip,
		*port,
		time.Duration(*readTimeout) * time.Second,
		time.Duration(*writeTimeout) * time.Second,
		time.Duration(*idleTimeout) * time.Second,
		time.Duration(*stopTimeout) * time.Second,
	}
}

func (c Config) GetReadTimeout() time.Duration {
	return c.readTimeout
}

func (c Config) GetWriteTimeout() time.Duration {
	return c.writeTimeout
}

func (c Config) GetIdleTimeout() time.Duration {
	return c.idleTimeout
}

func (c Config) GetStopTimeout() time.Duration {
	return c.stopTimeout
}

func (c Config) GetStaticRoot() string {
	return c.root
}

func (c Config) GetAddress() string {
	return c.ip + ":" + c.port
}

func (c Config) GetIndexHtml() string {
	return filepath.Join(c.GetStaticRoot(), "index.html")
}

func (c Config) GetName() string {
	return name
}
