package confighelpers

import (
	"bytes"

	"github.com/travmatth-org/qaas/internal/logger"
)

// ResetLogger instantiates a new buffer with a
// byte buffer backend, returns buffer
func ResetLogger() *bytes.Buffer {
	want := new(bytes.Buffer)
	logger.SetDestination(want)
	logger.SetLogger(nil)
	logger.GetLogger()
	return want
}
