package handlers

import (
	"io"
	"net/http"

	"github.com/travmatth-org/qaas/internal/logger"
)

// Static prepares and returns a http.Handler serving a single file
func Static(reader io.Reader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		bytes, err := io.Copy(w, reader)
		logger.InfoReq(r).Int64("wrote", bytes).Err(err).Msg("Served file")
	}
}
