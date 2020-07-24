package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_ServeStatic(t *testing.T) {
	_, s := configureServer(t)
	html := "../../web/www/static/index.html"
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("Error in test while creating request", err)
	}
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	ref, err := ioutil.ReadFile(html)
	if err != nil {
		t.Fatal("Error in test while reading response", err)
	} else if !bytes.Equal(ref, rr.Body.Bytes()) {
		t.Errorf("Expected %s But got %s", ref, rr.Body.String())
	}
}
