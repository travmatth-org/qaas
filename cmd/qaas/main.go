package main

import (
	"os"

	"github.com/travmatth-org/qaas/internal/api"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/server"
)

func main() {
	// Load config options
	c, err := config.New(
		config.WithConfigFile(config.LocateConfig),
		config.WithOverrides(os.Args),
	)
	if err != nil {
		logger.Error().Err(err).Msg("Error configuring server")
		os.Exit(1)
	}

	// Create API
	a, err := api.New(
		api.WithRegion(c.AWS.Region),
		api.WithSession,
		api.WithEC2(config.IsProd(c)),
		api.WithXray(config.IsProd(c)),
		api.WithNewDynamoDBClient(c),
	)
	if err != nil {
		logger.Error().Err(err).Msg("Error configuring API")
		os.Exit(1)
	}

	// Create Server
	s, err := server.New(c,
		server.WithStatic,
		server.WithAPI(a),
		server.WithRegisterStaticPages)
	if err != nil {
		logger.Error().Err(err).Msg("Error initializing Server")
		os.Exit(1)
	}

	// Listen for incoming connections
	if err := s.AcceptConnections(); err != nil {
		logger.Error().Err(err).Msg("Shutting down server")
		os.Exit(1)
	}
}
