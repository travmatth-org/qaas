package server

import (
	"io"
	"net/http"

	"github.com/travmatth-org/qaas/internal/logger"
)

// ServeStatic prepares and returns a http.Handler serving a single file
func (s *Server) ServeStatic(key string) http.HandlerFunc {
	reader := s.fs.Use(key)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		bytes, err := io.Copy(w, reader)
		logger.InfoReq(r).
			Int64("bytes written", bytes).
			Err(err).
			Msg("Served file")
	}
}
