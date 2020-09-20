package server

import (
	"net/http"

	"github.com/Travmatth/qaas/internal/logger"
)

// ServeStatic prepares and returns a http.Handler serving a single file
func (s *Server) ServeStatic(key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, key)
		logger.InfoReq(r).Str("file", key).Msg("Served file")
	}
}
