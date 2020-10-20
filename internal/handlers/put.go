package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/travmatth-org/qaas/internal/clean"
	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"
)

func (h *Handler) put(reader io.Reader) *types.QuoteRes {
	var (
		quote = types.NewQuote()
		res   = types.NewQuoteRes()
	)

	if err := json.NewDecoder(reader).Decode(&quote); err != nil {
		return res.WithErr(err)
	} else if err := clean.Quote(quote); err != nil {
		return res.WithErr(err)
	}
	quote.GenerateID()
	author, topics := types.RecordsFromQuote(quote)
	err := h.api.Put(
		h.api.PutWithQuote(quote),
		h.api.PutWithAuthor(author),
		h.api.PutWithTopics(topics))
	return res.WithQuote(quote).WithErr(err)
}

// Put reads & validates the given *Quote struct, before saving to DynamoDB
func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	defer r.Body.Close()
	body := h.put(r.Body)
	wr, err := w.Write(body.JSON())
	logger.InfoReq(r).
		Int("bytes", wr).
		Err(err).
		Object("quote", body).
		Msg("Saved Quote")
}
