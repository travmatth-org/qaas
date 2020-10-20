package handlers

import (
	"net/http"

	"github.com/travmatth-org/qaas/internal/logger"
)

// Random endpoint returns a random quote struct from the dynamodb api
func (h *Handler) Random(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	q := h.api.Random()
	bytes, err := w.Write(q.JSON())
	logger.InfoReq(r).
		Err(err).
		Int("wrote", bytes).
		Object("quote", q).
		Msg("Served Random Quote")
}
