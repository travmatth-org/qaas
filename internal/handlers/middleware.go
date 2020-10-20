package handlers

import (
	"io"
	"net/http"

	"github.com/travmatth-org/qaas/internal/logger"
)

// Log documents incoming requests before passing to child handler
func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoReq(r).Msg("Received request")
		h.ServeHTTP(w, r)
	})
}

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

// Static prepares and returns a http.Handler serving a single file
func Static(reader io.Reader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		bytes, err := io.Copy(w, reader)
		logger.InfoReq(r).Int64("wrote", bytes).Err(err).Msg("Served file")
	}
}
