package config

import (
	"errors"

	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"
	"gopkg.in/yaml.v2"
)

const (
	Development = "DEVELOPMENT"
	Test        = "TEST"
	Production  = "PRODUCTION"
	ParseError  = "Parsing Error: All opts must be passed in --flag <val> format"
)

// Config manages the configuration options of the program.
// Program options are configured first attempting to locate and open
// `httpd.yml`, first in the current working directory then in /etc/qaas
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
			Endpoint        string `yaml:"dev_db_endpoint" env:"DB_ENDPOINT" cli:"db-endpoint"`
			PaginationLimit int64  `yaml:"pagination_limit" env:"PAGINATION" cli:"db-pagination"`
			Table           struct {
				Quote  string `yaml:"quote" env:"QUOTE_TABLE" cli:"quote-table"`
				Author string `yaml:"author" env:"AUTHOR_TABLE" cli:"author-table"`
				Topic  string `yaml:"topic" env:"TOPIC_TABLE" cli:"topic-table"`
			}
		}
	}
}

type opts func(c *Config) (*Config, error)

func New(opts ...opts) (*Config, error) {
	var err error
	c := &Config{}
	for _, opt := range opts {
		if c, err = opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// WithConfigFile locates and parses the config file into the *config struct
func WithConfigFile(locate func() (types.AFSFile, error)) opts {
	return func(c *Config) (*Config, error) {
		file, err := locate()
		if err != nil {
			logger.Error().Err(err).Msg("Error locating config file")
			return nil, err
		}
		defer file.Close()
		return c, yaml.NewDecoder(file).Decode(c)
	}
}

func WithUpdates(opts []string) opts {
	return func(c *Config) (*Config, error) {
		n := len(opts)
		if n%2 != 0 {
			return nil, errors.New(ParseError)
		}
		m := make(map[string]string, n/2)
		for i := 0; i < n; i += 2 {
			m[opts[i][2:]] = opts[i+1]
		}
		err := ParseOverrides(c, m)
		return c, err
	}
}

func IsProd(c *Config) bool {
	return c.Env == Production
}
