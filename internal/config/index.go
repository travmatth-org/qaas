package config

const (
	DefaultRoot = "/srv/www/static"
	DefaultPort = "80"
)

type Config struct {
	root string
	port string
}

func NewConfig(root, port string) (c *Config) {
	return &Config{root, port}
}

func (c Config) GetStaticRoot() string {
	return c.root
}

func (c Config) GetPort() string {
	return ":" + c.port
}
