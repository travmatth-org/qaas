package handlers

import (
	"net/http"

	"github.com/travmatth-org/qaas/internal/api"
	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"
)

// Get returns a http endpoint listening for Get quote by attribute requests
func Get(a *api.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			m = types.NewMultiQuoteRes()
		)
		req, err := types.ValidateRequestFrom(r.Body)
		if err != nil {
			m.WithErr(err)
			goto FAIL
		}
		if req.Name == "author" {
			m = a.GetQuotesByAuthor(req.Name, req.Start)
		} else { // == "topic"
			m = a.GetQuotesByTopic(req.Name, req.Start)
		}
		if m.Err != nil {
			goto FAIL
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(m.JSON())
		logger.InfoReq(r).Object("quotes", m).Msg("Saved Quote")
	FAIL:
		http.Error(w, m.Error(), http.StatusBadRequest)
	}
}
