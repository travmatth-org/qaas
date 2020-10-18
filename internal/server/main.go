package server

import (
	"os"

	"github.com/travmatth-org/qaas/internal/afs"
	"github.com/travmatth-org/qaas/internal/api"
	"github.com/travmatth-org/qaas/internal/config"
	"github.com/travmatth-org/qaas/internal/logger"
)

// Main configures and starts a server with
// specified configuration, file system, and aws service clients.
// Returns 1 on server error, or 0 on graceful shutdown.
func Main() int {
	// init filesystem
	fs := afs.New().WithCachedFs()
	// Load config options
	c, err := config.New(
		config.WithConfigFile(fs.Open),
		config.WithUpdates(os.Args[1:]))
	if err != nil {
		logger.Error().Msg("Error configuring server")
		return 1
	}

	// Create API
	a, err := api.New(
		api.WithRegion(c.AWS.Region),
		api.WithSession,
		api.WithEC2(config.IsProd(c)),
		api.WithXray(config.IsProd(c)),
		api.WithNewDDB(c))
	if err != nil {
		logger.Error().Msg("Error configuring API")
		return 1
	}

	// Create Server
	s, err := New(c,
		WithFS(fs),
		WithStatic,
		WithAPI(a),
		WithStaticPages(config.IsProd(c)))
	if err != nil {
		logger.Error().Msg("Error initializing Server")
		return 1
	}
	defer fs.CloseAll()

	// Listen for incoming connections
	if err := s.AcceptConnections(); err != nil {
		logger.Error().Msg("Error shutting down server")
		return 1
	}
	return 0
}
