package config_utils

import (
	"bytes"

	"github.com/Travmatth/faas/internal/logger"
)

// ResetLogger instantiates a new buffer with a
// byte buffer backend, returns buffer
func ResetLogger() *bytes.Buffer {
	want := new(bytes.Buffer)
	logger.Destination, logger.Instance = want, nil
	logger.GetLogger()
	return want
}
