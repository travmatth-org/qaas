package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Travmatth/faas/internal/config"
	"github.com/Travmatth/faas/internal/logger"
	config_utils "github.com/Travmatth/faas/test/utils/config"
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

func configureServer() (*bytes.Buffer, *Server) {
	logged := config_utils.ResetLogger()
	c := config.New()
	c.Port = "8080"
	c.Static = "../../web"
	s := New(c)
	return logged, s
}

func TestServer_configureMiddlewareConfiguresLogs(t *testing.T) {
	// configure a server endpoint, mocking out logs for a buffer
	logged, s := configureServer()
	f := func(w http.ResponseWriter, r *http.Request) {
		logger.InfoReq(r).Bool("test", true).Msg("Succeeded")
	}
	h := s.configureMiddleware().ThenFunc(http.HandlerFunc(f))

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

func makeRequest(t *testing.T, s *Server, endpoint string) *bytes.Buffer {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatal("Error in test while creating request: ", err)
	}
	rr := httptest.NewRecorder()
	s.ServeHTTP(rr, req)
	return rr.Body
}

func TestServer_HomeAnd404Routes(t *testing.T) {
	// configure a server endpoint, mocking out logs for a buffer
	_, s := configureServer()
	if err := s.RegisterHandlers(); err != nil {
		t.Fatal("Error creating server endpoints: ", err)
	}

	tests := []struct {
		name     string
		endpoint string
		file     string
	}{
		{"TestHomeRoute", "/", s.GetIndexHTML()},
		{"Test404Route", "/foobar", s.Get404()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// read original file from filesystem
			want, err := ioutil.ReadFile(tt.file)
			if err != nil {
				t.Fatal("Error in test while reading ref file: ", err)
			}

			// make request to server
			f := makeRequest(t, s, tt.endpoint)
			got, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal("Error in test while reading request to []byte: ", err)
			} else if bytes.Compare(got, want) != 0 {
				t.Fatal("Error incorrect body returned for: ", tt.endpoint, tt.file)
			}
		})
	}
}
