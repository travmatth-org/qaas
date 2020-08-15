package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func mockTime() time.Duration { return time.Duration(0) }
func mockStr() string         { return "" }
func mockBool() bool          { return false }

type method struct {
	duration func() time.Duration
	str      func() string
	boolean  func() bool
}

type want struct {
	duration time.Duration
	str      string
	boolean  bool
}

var cnv = func(i int) time.Duration { return time.Second * time.Duration(i) }

func TestConfig(t *testing.T) {
	old := os.Args
	os.Args = []string{"faas"}
	defer func() { os.Args = old }()
	cwd, _ := os.Getwd()
	cwd = filepath.Join(cwd, "web", "www", "static")

	c := Build()
	if c == nil {
		t.Fatal("Error building configuration")
	}

	tests := []struct {
		name   string
		method method
		want   want
		kind   string
	}{
		{
			"TestGetReadTimeout",
			method{c.GetReadTimeout, mockStr, mockBool},
			want{cnv(defaultReadTimeout), "", false},
			"duration",
		},
		{
			"TestGetWriteTimeout",
			method{c.GetWriteTimeout, mockStr, mockBool},
			want{cnv(defaultWriteTimeout), "", false},
			"duration",
		},
		{
			"TestGetIdleTimeout",
			method{c.GetIdleTimeout, mockStr, mockBool},
			want{cnv(defaultIdleTimeout), "", false},
			"duration",
		},
		{
			"TestGetStopTimeout",
			method{c.GetStopTimeout, mockStr, mockBool},
			want{cnv(defaultStopTimeout), "", false},
			"duration",
		},
		{
			"TestGetAddress",
			method{mockTime, c.GetAddress, mockBool},
			want{time.Duration(0), "0.0.0.0:80", false},
			"string",
		},
		{
			"TestGetIndexHTML",
			method{mockTime, c.GetIndexHTML, mockBool},
			want{time.Duration(0), filepath.Join(cwd, "index.html"), false},
			"string",
		},
		{
			"TestGet404",
			method{mockTime, c.Get404, mockBool},
			want{time.Duration(0), filepath.Join(cwd, "404.html"), false},
			"string",
		},
		{
			"TestIsProd",
			method{mockTime, mockStr, c.IsProd},
			want{time.Duration(0), "", false},
			"string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.kind == "duration" && tt.method.duration() != tt.want.duration {
				t.Errorf("Want %s : Got %s", tt.want.str, tt.method.duration())
			} else if tt.kind == "string" && tt.method.str() != tt.want.str {
				t.Errorf("Want %s : Got %s", tt.want.str, tt.method.str())
			} else if tt.kind == "boolean" && tt.method.boolean() != tt.want.boolean {
				t.Errorf("Want %s : Got %t", tt.want.str, tt.method.boolean())
			}
		})
	}
}

func TestProd(t *testing.T) {
	c := New()
	if c.IsProd() != false {
		t.Fatalf("Error: Incorrect default prod val")
	}
}
