package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/travmatth-org/qaas/internal/clean"
	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"
)

func (h *Handler) get(reader io.Reader) *types.MultiQuoteRes {
	var (
		req = types.NewGetBatch()
		res = types.NewMultiQuoteRes()
	)

	if err := json.NewDecoder(reader).Decode(req); err != nil {
		return res.WithErr(err)
	} else if err = clean.Query(req); err != nil {
		return res.WithErr(err)
	} else if req.Name == "author" {
		return h.api.Get(h.tables.Author, req.Value, req.Start)
	}
	return h.api.Get(h.tables.Topic, req.Value, req.Start)
}

// Get returns a http endpoint listening for Get quote by attribute requests
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	defer r.Body.Close()
	body := h.get(r.Body)
	wr, err := w.Write(body.JSON())
	logger.InfoReq(r).
		Int("bytes", wr).
		Err(err).
		Object("quotes", body).
		Msg("Saved Quote")
}
