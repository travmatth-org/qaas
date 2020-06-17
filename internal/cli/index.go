package cli

import (
	"flag"

	"github.com/Travmatth/faas/internal/config"
)

// RegisterOpts initializes command line opts
func RegisterOpts() (string, string) {
	message := "Directory static assets served from"
	static := flag.String("static", config.DefaultRoot, message)
	message = "Port server should listen on"
	port := flag.String("port", config.DefaultPort, message)
	flag.Parse()
	return *static, *port
}
