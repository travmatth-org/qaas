package server

import (
	"io/ioutil"
	"net/http"

	"github.com/Travmatth/faas/internal/logger"
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
			logger.ErrorReq(r).Err(err).Str("file", key).
				Msg("Error serving static file from memory")
			return
		}
		logger.InfoReq(r).Int("file_size", size).Str("file", key).
			Msg("Served static file from memory")
	}
}
