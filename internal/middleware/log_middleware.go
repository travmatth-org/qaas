package middleware

import (
	"net/http"

	"github.com/Travmatth/qaas/internal/logger"
)

// Log documents incoming requests before passing to child handler
func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoReq(r).Msg("Received request")
		h.ServeHTTP(w, r)
	})
}
