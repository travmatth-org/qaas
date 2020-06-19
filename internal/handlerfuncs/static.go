package handlerfuncs

import (
	"net/http"
	"path"
	"strings"

	"github.com/Travmatth/faas/internal/config"
	"github.com/rs/zerolog/hlog"
)

func extractAndCleanPath(r *http.Request) string {
	p := r.URL.Path
	if !strings.HasPrefix(p, "/") {
		p = "/" + r.URL.Path
	}
	return path.Clean(p)
}

func Home(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
		hlog.FromRequest(r).Info().
			Msg("Request completed")
	}
}
