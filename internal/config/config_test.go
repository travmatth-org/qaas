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

	c := Build()

	tests := []struct {
		name   string
		method method
		want   want
		kind   string
	}{
		{
			"TestGetReadTimeout",
			method{c.GetReadTimeout, mockStr, mockBool},
			want{cnv(DefaultReadTimeout), "", false},
			"duration",
		},
		{
			"TestGetWriteTimeout",
			method{c.GetWriteTimeout, mockStr, mockBool},
			want{cnv(DefaultWriteTimeout), "", false},
			"duration",
		},
		{
			"TestGetIdleTimeout",
			method{c.GetIdleTimeout, mockStr, mockBool},
			want{cnv(DefaultIdleTimeout), "", false},
			"duration",
		},
		{
			"TestGetStopTimeout",
			method{c.GetStopTimeout, mockStr, mockBool},
			want{cnv(DefaultStopTimeout), "", false},
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
			want{time.Duration(0), filepath.Join(DefaultRoot, "index.html"), false},
			"string",
		},
		{
			"TestGet404",
			method{mockTime, c.Get404, mockBool},
			want{time.Duration(0), filepath.Join(DefaultRoot, "404.html"), false},
			"string",
		},
		{
			"TestIsDev",
			method{mockTime, mockStr, c.IsDev},
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

func TestDev(t *testing.T) {
	c := New()
	if c.IsDev() != false {
		t.Fatalf("Error: Incorrect default dev val")
	}
}
