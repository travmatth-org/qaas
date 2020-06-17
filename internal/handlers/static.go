package handlers

import (
	"fmt"
	"net/http"

	"github.com/Travmatth/faas/internal/config"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func Static(c *config.Config) http.HandlerFunc {
	dir := http.Dir(c.GetStaticRoot())
	h := http.FileServer(dir)
	return func(w http.ResponseWriter, r *http.Request) {
		srw := &statusResponseWriter{ResponseWriter: w}
		h.ServeHTTP(srw, r)
		if srw.status >= http.StatusBadRequest {
			fmt.Println("error")
		} else {
			fmt.Println("completed request")
		}
	}
}
