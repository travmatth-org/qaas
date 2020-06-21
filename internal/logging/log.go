package logging

import (
	"os"

	"github.com/rs/zerolog"
)

// NewLogger configures and returns a new log isntance
func NewLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return zerolog.New(os.Stdout).With().
		Timestamp().
		Str("role", "faas").
		Logger()
}
