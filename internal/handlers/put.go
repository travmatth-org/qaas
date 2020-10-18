package handlers

import (
	"net/http"

	"github.com/travmatth-org/qaas/internal/api"
	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"
)

// Put reads & validates the given *Quote struct, before saving to DynamoDB
func Put(a *api.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			q            = types.NewQuoteRes()
			quote        = types.NewQuote()
			author       = types.NewRecord()
			topics       = []*types.Record{}
			err    error = nil
		)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		quote, err = types.ValidateQuoteFrom(r.Body)
		if err != nil {
			goto FAIL
		}
		author, topics = types.RecordsFromQuote(quote)
		err = a.PutNewQuote(quote, author, topics)
		q.WithQuote(quote).WithErr(err)
		if err != nil {
			goto FAIL
		}
		w.Write(q.JSON())
		logger.InfoReq(r).Object("quote", q).Msg("Saved Quote")
		return
	FAIL:
		http.Error(w, q.Error(), http.StatusBadRequest)
	}
}
