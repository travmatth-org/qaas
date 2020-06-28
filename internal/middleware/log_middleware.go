package middleware

import (
	"net/http"

	"github.com/Travmatth/faas/internal/logger"
)

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoReq(r).Msg("Received request")
		h.ServeHTTP(w, r)
	})
}
