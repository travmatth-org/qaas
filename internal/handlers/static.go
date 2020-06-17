package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/Travmatth/faas/internal/config"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func Handle404(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Error"))
}

func Static(c *config.Config) http.HandlerFunc {
	dir := http.Dir(c.GetStaticRoot())
	fs := http.FileServer(dir)
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if !strings.HasPrefix(p, "/") {
			p = "/" + r.URL.Path
		}
		if f, err := dir.Open(path.Clean(p)); err != nil {
			if os.IsNotExist(err) {
				Handle404(w, r)
				fmt.Println("Handled nonexistend request")
			} else {
				fmt.Println(err)
			}
		} else {
			f.Close()
			fs.ServeHTTP(w, r)
			fmt.Println("Completed request")
		}
	}
}
