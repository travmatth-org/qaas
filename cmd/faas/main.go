package main

import (
	"os"

	"github.com/Travmatth/faas/internal/config"
	"github.com/Travmatth/faas/internal/logging"
	"github.com/Travmatth/faas/internal/server"
)

func main() {
	// Config vals of server
	c := config.Build()

	// Configure log
	log := logging.NewLogger()

	// Create server
	srv := server.New(c, log)
	if err := srv.RegisterHandlers(); err != nil {
		log.Fatal().Err(err).Msg("Launch aborted")
		os.Exit(1)
	}

	// run server
	os.Exit(srv.AcceptConnections())
}
