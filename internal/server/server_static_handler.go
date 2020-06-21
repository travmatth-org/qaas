package server

import (
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/hlog"
)

func (s *Server) loadFileIntoMemory(key, name string) error {
	f, err := ioutil.ReadFile(name)
	if err == nil {
		s.static[key] = f
	}
	return err
}

// ServeStatic prepares and returns a http.Handler serving a single
// file located in the map of the server
func (s *Server) ServeStatic(key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		size, err := w.Write(s.static[key])
		if err != nil {
			hlog.FromRequest(r).Error().
				Err(err).
				Str("file", key).
				Msg("Error serving static file from memory")
			return
		}
		hlog.FromRequest(r).Info().
			Str("file", key).
			Int("file_size", size).
			Msg("Served static file from memory")
	}
}
