package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func TestSetLogger(t *testing.T) {
	want := zerolog.New(os.Stderr).With().Logger()
	SetLogger(&want)
	got := *instance
	if !reflect.DeepEqual(want, got) {
		t.Fatal("SetLogger did not set correctly")
	}
}

func TestGetLogger(t *testing.T) {
	first := GetLogger()
	if first == nil {
		t.Fatal("SetLogger returning nil")
	} else if !reflect.DeepEqual(first, GetLogger()) {
		t.Fatal("Singleton not returning same instance")
	}
}

func TestError(t *testing.T) {
	type output struct {
		Level   string `json:"level"`
		Test    bool   `json:"test"`
		Message string `json:"message"`
	}

	var out output
	var want bytes.Buffer
	reference := output{"error", true, "Succeeded"}

	destination, instance = &want, nil
	GetLogger()

	Error().Bool("test", true).Msg("Succeeded")

	if err := json.Unmarshal(want.Bytes(), &out); err != nil {
		t.Fatal(err, ". Got: ", want.String())
	} else if !reflect.DeepEqual(out, reference) {
		t.Fatalf("Incorrect error:\n\t%+v\n\twant: %+v\n", out, reference)
	}
}

type output struct {
	Level   string `json:"level"`
	Test    bool   `json:"test"`
	Message string `json:"message"`
}

func TestErrorReq(t *testing.T) {
	var out output
	var want bytes.Buffer

	// create log output and initialize log
	destination, instance = &want, nil
	GetLogger()

	// handlerfunc that outputs to log
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ErrorReq(r).Bool("test", true).Msg("Succeeded")
	})

	// error log is created and passed into handler through middleware
	handler := hlog.NewHandler(*GetLogger())(f)

	// execute handler, passing request and responserecorder
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// verify logger from context prints correct output
	ref := output{"error", true, "Succeeded"}
	if err := json.Unmarshal(want.Bytes(), &out); err != nil {
		t.Fatal(err, ". Got: ", want.String())
	} else if !reflect.DeepEqual(out, ref) {
		t.Fatalf("Incorrect error:\n\t got: %+v\n\twant: %+v\n", out, ref)
	}
}

func TestInfo(t *testing.T) {
	tests := []struct {
		name string
		want *zerolog.Event
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Info(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Info() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfoReq(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want *zerolog.Event
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InfoReq(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InfoReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWarn(t *testing.T) {
	tests := []struct {
		name string
		want *zerolog.Event
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Warn(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Warn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWarnReq(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want *zerolog.Event
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WarnReq(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WarnReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDebug(t *testing.T) {
	tests := []struct {
		name string
		want *zerolog.Event
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Debug(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Debug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDebugReq(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want *zerolog.Event
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DebugReq(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DebugReq() = %v, want %v", got, tt.want)
			}
		})
	}
}
