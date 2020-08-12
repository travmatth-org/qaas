package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/Travmatth/faas/internal/config"
	"github.com/Travmatth/faas/internal/logger"
	confighelpers "github.com/Travmatth/faas/test/utils/config"
)

const (
	LoopbackTestPort = "127.0.0.1:8080"
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

func configureServer(t *testing.T) (*bytes.Buffer, *Server) {
	logged := confighelpers.ResetLogger()
	c := config.New()
	// listen on loopback interface only
	c.Port = LoopbackTestPort
	c.Static = "../../web/www/static"
	s := New(c)
	s.RegisterHandlers()
	return logged, s
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

func TestServer_configureMiddlewareConfiguresLogs(t *testing.T) {
	// configure a server endpoint, mocking out logs for a buffer
	logged, s := configureServer(t)
	f := func(w http.ResponseWriter, r *http.Request) {
		logger.InfoReq(r).Bool("test", true).Msg("Succeeded")
	}
	h := s.WrapRoute(http.HandlerFunc(f))

	// perform a request to server, triggering logging middleware
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	// middleware should test server logs correctly
	var got middlewareRef
	first := strings.Split(logged.String(), "\n")[0]
	if err := json.Unmarshal([]byte(first), &got); err != nil {
		t.Fatal("Error unmarshaling log object: ", err, logged.String())
	}
}

func TestServer_HomeAnd404Routes(t *testing.T) {
	// configure a server endpoint, mocking out logs for a buffer
	_, s := configureServer(t)

	tests := []struct {
		name     string
		endpoint string
		file     string
	}{
		{"TestHomeRoute", "/", "../../web/www/static/index.html"},
		{"Test404Route", "/foobar", "../../web/www/static/404.html"},
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
			} else if !bytes.Equal(got, want) {
				t.Fatal("Error incorrect body returned for: ", tt.endpoint, tt.file, "\n", string(got), "\n", string(want))
			}
		})
	}
}

func startServerTest(t *testing.T) (*Server, chan int) {
	// configure a server endpoint, mocking out logs for a buffer
	_, s := configureServer(t)
	ch := make(chan int, 1)
	go func() {
		ch <- s.AcceptConnections()
	}()
	<-s.startedChannel

	// mock an http request
	res, err := http.Get("http://" + LoopbackTestPort)
	if err != nil {
		t.Fatal("Error in test while mocking request: ", err)
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("Error in test while decoding mock response: ", err)
	}
	return s, ch
}

func TestServer_SignalShutdown(t *testing.T) {
	_, ch := startServerTest(t)
	// send shutdown signal
	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	if err != nil {
		t.Fatal("Error in killing test server: ", err)
	}
	// block until shutdown received, or timeout exceeded
	select {
	case status := <-ch:
		if status != ok {
			t.Fatal("Error: incorrect shutdown val: ", status)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Error: timeout exceeded")
	}
}

func TestServer_ErrorShutdown(t *testing.T) {
}

func TestServer_ChecksListenerNotNil(t *testing.T) {
	s := New(config.New())
	go s.startServing()
	select {
	case err := <-s.errorChannel:
		if err == nil {
			t.Fatal("Error: StartServing() should err on nil http listener")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Error: timeout exceeded")
	}
}
