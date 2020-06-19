package config

import "time"

const (
	DefaultRoot    = "/srv/www/static"
	DefaultPort    = "80"
	DefaultTimeout = 10
	name           = "faas"
)

type Config struct {
	root    string
	port    string
	timeout time.Duration
}

func NewConfig(root, port string, timeout time.Duration) (c *Config) {
	return &Config{root, port, timeout}
}

func (c Config) GetTimeout() time.Duration {
	return c.timeout
}

func (c Config) GetStaticnAme() string {
	return name
}

func (c Config) GetStaticRoot() string {
	return c.root
}

func (c Config) GetPort() string {
	return ":" + c.port
}
