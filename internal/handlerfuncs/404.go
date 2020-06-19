package handlerfuncs

import (
	"net/http"

	"github.com/rs/zerolog/hlog"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Error"))
	hlog.FromRequest(r).Info().
		Msg("Path not found")
}
