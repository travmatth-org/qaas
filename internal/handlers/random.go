package handlers

import (
	"net/http"

	"github.com/travmatth-org/qaas/internal/api"
	"github.com/travmatth-org/qaas/internal/logger"
)

// Random endpoint returns a random quote struct from the dynamodb api
func Random(a *api.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		q := a.GetRandomQuote()
		bytes, err := w.Write(q.JSON())
		logger.InfoReq(r).
			Err(err).
			Int("wrote", bytes).
			Object("quote", q).
			Msg("Served Random Quote")
	}
}
