package handlerfuncs

import (
	"net/http"
	"os"
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

func Static(c *config.Config) http.HandlerFunc {
	dir := http.Dir(c.GetStaticRoot())
	fs := http.FileServer(dir)
	return func(w http.ResponseWriter, r *http.Request) {
		p := extractAndCleanPath(r)
		switch f, err := dir.Open(p); {
		case os.IsNotExist(err):
			NotFoundHandler(w, r)
			hlog.FromRequest(r).Info().Msg("File not found")
		case err != nil:
			hlog.FromRequest(r).Info().Err(err).Msg("Completed Request")
		default:
			f.Close()
			fs.ServeHTTP(w, r)
			hlog.FromRequest(r).Info().Msg("Completed Request")
		}
	}
}
