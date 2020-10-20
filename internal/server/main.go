package server

import (
	"os"

	"github.com/travmatth-org/qaas/internal/afs"
	"github.com/travmatth-org/qaas/internal/api"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/handlers"
	"github.com/travmatth-org/qaas/internal/logger"
)

// Main configures and starts a server with specified
// configuration, file system, and aws service clients.
// Returns int representing shutdown status.
func Main() int {
	// init filesystem
	fs := afs.New().WithCachedFs()
	// Load config options
	c, err := config.New(config.WithFile(fs.Open), config.Update(os.Args[1:]))
	if err != nil {
		logger.Error().Err(err).Msg("Error configuring server")
		return 1
	}

	// Create API
	a, err := api.New(
		api.WithPaginationLimit(c.AWS.DynamoDB.PaginationLimit),
		api.WithTables(c.AWS.DynamoDB.Table),
		api.WithEC2(config.IsProd(c)),
		api.WithXray(config.IsProd(c)),
		api.WithDynamoDBService(c))
	if err != nil {
		logger.Error().Err(err).Msg("Error configuring API")
		return 1
	}

	// Create Handlers
	home := c.Net.Static
	h, err := handlers.New(handlers.WithFS(fs, home), handlers.WithAPI(a))
	if err != nil {
		logger.Error().Err(err).Msg("Error configuring http handlers")
		return 1
	}
	defer fs.CloseAll()

	// Create Server
	s, err := New(c, WithHandlers(h.RouteMap(), config.IsProd(c)))
	if err != nil {
		logger.Error().Err(err).Msg("Error initializing Server")
		return 1
	}

	// Listen for incoming connections
	if err := s.AcceptConnections(); err != nil {
		logger.Error().Err(err).Msg("Error shutting down server")
		return 1
	}

	return 0
}
