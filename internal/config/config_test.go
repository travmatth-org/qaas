package config

import (
	"errors"
	"os"
	"testing"

	"github.com/travmatth-org/qaas/internal/fs"
)

var config = []byte(`env: PRODUCTION
net:
  static: /web/www/static
  ip: 0.0.0.0
  port: :80
  liveness: 10
timeouts:
  read: 5
  write: 5
  stop: 5
  idle: 5
aws:
  region: us-west-1
  dynamodb:
    endpoint: http://localhost:8000
    pagination_limit: 5
    tables:
      quote: qaas-quote-table
      author: qaas-author-table
      topic: qaas-topic-table`)

func TestNew(t *testing.T) {
	for _, tt := range []struct {
		f   opts
		err error
	}{{
		func(c *Config) (*Config, error) {
			c.Env = Test
			return c, nil
		}, nil,
	}, {
		func(c *Config) (*Config, error) {
			return nil, errors.New("")
		}, errors.New(""),
	}} {
		c, err := New(tt.f)
		if tt.err == nil && (err != nil || c.Env != Test) {
			t.Fatalf("Error creating new config from opt")
		} else if tt.err != nil && err == nil {
			t.Fatalf("Error creating new config from opt")
		}
	}
}

func TestWithConfigFile(t *testing.T) {
	env, location := "QAAS_CONFIG", "/etc/qaas/httpd.yml"
	fileSystem := fs.New().WithMemFs()
	err := fileSystem.Write(location, config, fs.ReadAllWriteUser)
	if err != nil {
		t.Errorf("Error configuring file system: %+v", err)
	}
	os.Setenv(env, location)
	defer os.Unsetenv(env)

	c, err := WithConfigFile(fileSystem.Locate("QAAS_CONFIG"))(&Config{})
	if err != nil || c.Env != Production {
		t.Fatalf("Error unmarshaling file to config struct")
	}
}

func TestWithUpdates(t *testing.T) {
	os.Setenv("QAAS_ENV", "DEVELOPMENT")
	defer os.Unsetenv("QAAS_ENV")
	c := &Config{}
	if c, err := WithUpdates([]string{"--static", "foo"})(c); err != nil {
		t.Fatalf("Error overriding values in config struct: %s", err)
	} else if c.Env != Development || c.Net.Static != "foo" {
		t.Fatalf("Error: Overriding values in config struct failed")
	}
}

func TestIsProd(t *testing.T) {
	c := Config{}
	c.Env = Production
	if IsProd(&c) != true {
		t.Fatalf("Error: IsProd should return True when Env == Production")
	}
	c.Env = Development
	if IsProd(&c) != false {
		t.Fatalf("Error: IsProd should return False when Env != Production")
	}
}
