package handlers

import (
	"net/http"

	"github.com/travmatth-org/qaas/internal/logger"
)

// Recover catches panics in downstream handlers, sends err to alert client
func Recover(next http.Handler) http.Handler {
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
