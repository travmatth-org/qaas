package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Travmatth/faas/internal/config"
	"github.com/Travmatth/faas/internal/logger"
	config_utils "github.com/Travmatth/faas/test/utils/config"
	"github.com/gorilla/mux"
)

type middlewareRef struct {
	Level   string `json:"level"`
	Role    string `json:"role"`
	ID      string `json:"req_id"`
	Dest    string `json:"dest"`
	Test    bool   `json:"test"`
	Time    string `json:"time"`
	Caller  string `json:"caller"`
	Message string `json:"message"`
}

func TestServer_configureMiddlewareConfiguresLogs(t *testing.T) {
	// configure a server endpoint, mocking out logs for a buffer
	logged := config_utils.ResetLogger()
	c := config.New()
	c.Port = "8080"
	srv := New(c)
	f := func(w http.ResponseWriter, r *http.Request) {
		logger.InfoReq(r).Bool("test", true).Msg("Succeeded")
	}
	h := srv.configureMiddleware().ThenFunc(http.HandlerFunc(f))

	// perform a request to server, triggering logging middleware
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	// middleware should test server logs correctly
	var got middlewareRef
	if err := json.Unmarshal(logged.Bytes(), &got); err != nil {
		t.Fatal("Error unmarshaling log object: ", err)
	}
}

func TestServer_Routes(t *testing.T) {
	type fields struct {
		Config      *config.Config
		Router      *mux.Router
		Server      *http.Server
		stopTimeout time.Duration
		static      map[string][]byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Config:      tt.fields.Config,
				Router:      tt.fields.Router,
				Server:      tt.fields.Server,
				stopTimeout: tt.fields.stopTimeout,
				static:      tt.fields.static,
			}
			if err := s.RegisterHandlers(); (err != nil) != tt.wantErr {
				t.Errorf("Server.RegisterHandlers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
