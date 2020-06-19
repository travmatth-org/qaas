package handlerfuncs

import (
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Error"))
}
