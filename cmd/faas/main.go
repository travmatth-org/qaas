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

	// Create server
	srv := server.New(c)
	if err := srv.RegisterHandlers(); err != nil {
		logger.Error().Err(err).Msg("Launch aborted")
		os.Exit(1)
	}

	// run server
	os.Exit(srv.AcceptConnections())
}
