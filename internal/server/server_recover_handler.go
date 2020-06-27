package server

import (
	"net/http"

	"github.com/Travmatth/faas/internal/logger"
)

// RecoverHandler catches panics in downstream handlers,
// sends an error to gracefully alert the client
func (s *Server) RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error().Interface("panic", err).Msg("Handler crashed")
				status := http.StatusInternalServerError
				http.Error(w, http.StatusText(status), status)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
