package main

import (
	"os"

	"github.com/Travmatth/faas/internal/config"
	"github.com/Travmatth/faas/internal/logger"
	"github.com/Travmatth/faas/internal/server"
)

func main() {
	// Config vals of server
	c := config.Build()
	if c == nil {
		logger.Error().Msg("Error configuring server")
		os.Exit(1)
	}

	// Create server
	s := server.New(c)
	s.RegisterHandlers()

	// run server
	os.Exit(s.AcceptConnections())
}
