package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

type output struct {
	Level   string `json:"level"`
	Test    bool   `json:"test"`
	Message string `json:"message"`
}

const incorrectStructMessage = "Incorrect log:\n\thave%+v\n\twant: %+v\n"

func reset() *bytes.Buffer {
	want := new(bytes.Buffer)
	SetDestination(want)
	SetLogger(nil)
	GetLogger()
	return want
}

func TestSetLogger(t *testing.T) {
	var want bytes.Buffer
	var out output

	log := zerolog.New(&want).With().Logger()
	SetLogger(&log)
	log.Error().Bool("test", true).Msg("Succeeded")

	if !reflect.DeepEqual(log, *instance) {
		t.Fatal("SetLogger did not set correctly")
	}
	reference := output{"error", true, "Succeeded"}
	if err := json.Unmarshal(want.Bytes(), &out); err != nil {
		message := "Invalid log object: "
		t.Fatal(message, err, want.String())
	} else if !reflect.DeepEqual(out, reference) {
		t.Fatalf(incorrectStructMessage, out, reference)
	}
}

func TestGetLogger(t *testing.T) {
	if first := GetLogger(); first == nil {
		t.Fatal("SetLogger returning nil")
	} else if !reflect.DeepEqual(first, GetLogger()) {
		t.Fatal("Singleton not returning same instance")
	}
}

func TestLevelLogs(t *testing.T) {
	tests := []struct {
		name  string
		log   func() *zerolog.Event
		level string
	}{
		{"TestError", Error, "error"},
		{"TestInfo", Info, "info"},
		{"TestWarn", Warn, "warn"},
		{"TestDebug", Debug, "debug"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out output

			// reset log singleton
			want := reset()

			// trigger output to log
			tt.log().Bool("test", true).Msg("Succeeded")

			// verify logger prints correct output
			reference := output{tt.level, true, "Succeeded"}
			if err := json.Unmarshal(want.Bytes(), &out); err != nil {
				t.Fatal("Invalid log object: ", err, want.String())
			} else if !reflect.DeepEqual(out, reference) {
				t.Fatalf(incorrectStructMessage, out, reference)
			}
		})
	}
}

func TestLogsFromRequest(t *testing.T) {
	tests := []struct {
		name  string
		log   func(r *http.Request) *zerolog.Event
		level string
	}{
		{"TestErrorWithRequest", ErrorReq, "error"},
		{"TestInfoWithRequest", InfoReq, "info"},
		{"TestWarnWithRequest", WarnReq, "warn"},
		{"TestDebugWithRequest", DebugReq, "debug"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out output

			// reset log singleton
			want := reset()

			// handlerfunc that outputs to log
			f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tt.log(r).Bool("test", true).Msg("Succeeded")
			})

			// error log is created and passed into handler through middleware
			handler := hlog.NewHandler(*GetLogger())(f)

			// execute handler, passing request and responserecorder
			req, _ := http.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			// verify logger from context prints correct output
			ref := output{tt.level, true, "Succeeded"}
			if err := json.Unmarshal(want.Bytes(), &out); err != nil {
				t.Fatal(err, ". Got: ", want.String())
			} else if !reflect.DeepEqual(out, ref) {
				t.Fatalf(incorrectStructMessage, out, ref)
			}
		})
	}
}
