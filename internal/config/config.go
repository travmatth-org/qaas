package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"
	"gopkg.in/yaml.v2"
)

const (
	// Development is string literal for dev environments
	Development = "DEVELOPMENT"
	// Test is string literal for test environments
	Test = "TEST"
	// Production is string literal for prod environments
	Production = "PRODUCTION"
	parseError = "Parsing Error: All opts must be passed in --flag val format"
)

// Tables represent the DynamoDB tables of the QAAS service
type Tables struct {
	Quote  string `yaml:"quote" env:"QUOTE_TABLE" cli:"quote-table"`
	Author string `yaml:"author" env:"AUTHOR_TABLE" cli:"author-table"`
	Topic  string `yaml:"topic" env:"TOPIC_TABLE" cli:"topic-table"`
}

// Config manages the configuration options of the program.
// Program options are configured first attempting to locate and open
// `httpd.yml`, first in the path specified by `QAAS_CONFIG` then in /etc/qaas
// Default values given in config are overriden by env var specified under the
// `env:"<VAL>"` attribute tags, and then by cli flags specified under the
// `cli:"<VAL"` attribute tag
type Config struct {
	Env string `yaml:"env" env:"ENV" cli:"env"`

	Net struct {
		Static   string `yaml:"static" env:"STATIC" cli:"static"`
		IP       string `yaml:"ip" env:"IP" cli:"ip"`
		Port     string `yaml:"port" env:"PORT" cli:"port"`
		Liveness int    `yaml:"liveness" env:"liveness" cli:"liveness"`
	}

	Timeout struct {
		Read  int `yaml:"read" env:"READ_TIMEOUT" cli:"read-timeout"`
		Write int `yaml:"write" env:"WRITE_TIMEOUT" cli:"write-timeout"`
		Stop  int `yaml:"stop" env:"STOP_TIMEOUT" cli:"stop-timeout"`
		Idle  int `yaml:"idle" env:"IDLE_TIMEOUT" cli:"idle-timeout"`
	}

	AWS struct {
		Region   string `yaml:"region" env:"AWS_REGION" cli:"region"`
		DynamoDB struct {
			Endpoint        string `yaml:"endpoint" env:"DB_ENDPOINT" cli:"db-endpoint"`
			PaginationLimit int64  `yaml:"pagination_limit" env:"PAGINATION" cli:"db-pagination"`
			Table           Tables
		}
	}
}

// Opts is the type signature for optional configuration functions
type Opts func(c *Config) (*Config, error)

// New constructs and returns a configuration with the specified options
func New(opt ...Opts) (*Config, error) {
	var (
		err error = nil
		c         = &Config{}
	)
	for _, fn := range opt {
		if c, err = fn(c); err != nil {
			return nil, err
		}
	}
	return c, err
}

// find path to config, first under QAAS_CONFIG var, then /etc/qaas/httpd.yml
func choosePath() (string, error) {
	var (
		err  error = nil
		path       = os.Getenv("QAAS_CONFIG")
	)
	if path == "" {
		path, err = filepath.Abs(filepath.Join("etc", "qaas", "httpd.yml"))
		if err != nil {
			logger.Error().Err(err).Msg("Error locating config file")
		}
	}
	return path, err
}

// WithFile locates and parses the config file into the Config struct.
// Prefers file path given by `QAAS_CONFIG` environment variable, defaults to
// /etc/qaas/httpd.yml
func WithFile(open func(string) (types.AFSFile, error)) Opts {
	return func(c *Config) (*Config, error) {
		path, err := choosePath()
		if err != nil {
			return nil, err
		}
		file, err := open(path)
		if err != nil {
			logger.Error().Err(err).Msg("Error opening config file")
			return nil, err
		}
		defer file.Close()
		d := yaml.NewDecoder(file)
		d.SetStrict(true)
		err = d.Decode(c)
		if err != nil {
			logger.Error().Err(err).Msg("Error creating config")
		}
		return c, err
	}
}

// Update accepts the os.Args array and overrides the Config,
// first with available env vars then with cli options
func Update(options []string) Opts {
	return func(c *Config) (*Config, error) {
		n := len(options)
		if n%2 != 0 {
			return nil, errors.New(parseError)
		}
		m := make(map[string]string, n/2)
		for i := 0; i < n; i += 2 {
			m[options[i][2:]] = options[i+1]
		}
		err := ParseOverrides(c, m)
		return c, err
	}
}

// IsProd indicates whether server under production configuration
func IsProd(c *Config) bool {
	return c.Env == Production
}
