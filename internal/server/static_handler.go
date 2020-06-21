package server

import (
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/hlog"
)

func (s *Server) LoadFileIntoMemory(key, name string) error {
	f, err := ioutil.ReadFile(name)
	if err == nil {
		s.static[key] = f
	}
	return err
}

func (s *Server) ServeStatic(key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if size, err := w.Write(s.static[key]); err != nil {
			hlog.FromRequest(r).Error().
				Err(err).
				Str("file", key).
				Msg("Error serving static file from memory")
		} else {
			hlog.FromRequest(r).Info().
				Str("file", key).
				Int("file_size", size).
				Msg("Served static file from memory")
		}
	}
}
