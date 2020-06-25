package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Travmatth/faas/internal/config"
	config_utils "github.com/Travmatth/faas/test/utils/config"
)

func TestServer_RecoverHandler(t *testing.T) {
	// Create server
	config_utils.ResetLogger()
	s := New(config.New())

	// Create a crashing handler and wrap it with the recovery handler
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Error")
	})

	// create mock request and response
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("Error creating mock request", err)
	}
	rr := httptest.NewRecorder()

	// trigger panic
	s.RecoverHandler(f).ServeHTTP(rr, req)
}
