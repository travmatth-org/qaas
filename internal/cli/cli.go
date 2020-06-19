package cli

import (
	"flag"
	"time"

	"github.com/Travmatth/faas/internal/config"
)

// RegisterOpts initializes command line opts
func RegisterOpts() (string, string, time.Duration) {
	message := "Directory static assets served from"
	static := flag.String("static", config.DefaultRoot, message)
	message = "Port server should listen on"
	port := flag.String("port", config.DefaultPort, message)
	message = "Default timeout period for HTTP responses"
	timeout := flag.Int("timeout", config.DefaultTimeout, message)
	flag.Parse()
	return *static, *port, time.Duration(*timeout) * time.Second
}
