package config

import (
	"reflect"
	"os"
	"strconv"
	"errors"
)

const (
	ENVTag = "env"
	CLITag = "cli"
)

type parser struct {
	prefix string
	opts map[string]string
}

func (p *parser) unknownOpts() string {
	s := ""
	for key, value := range p.opts {
		s += key + "=" + value + "\n"
	}
	return s
}

func (p *parser) getTagVal(v reflect.StructField) string {
	if tag, ok := p.opts[v.Tag.Get(CLITag)]; ok {
		delete(p.opts, v.Tag.Get(CLITag))
		return tag
	} else if tag = v.Tag.Get(ENVTag); tag == "" {
		return ""
	} else {
		return os.Getenv(p.prefix + tag)
	}
}

func (p *parser) walkFields(v reflect.Value, tag string) error {
	for i := 0; i < v.NumField(); i++ {
		sf := v.Type().Field(i)
		switch k := sf.Type.Kind(); {
		case k == reflect.Struct:
			tag = p.getTagVal(sf)
			if err := p.walkFields(v.Field(i), tag); err != nil {
				return err
			}
		case k == reflect.String:
			v.Field(i).SetString(tag)
		case k != reflect.Int && k != reflect.Int64:
			return nil
		}
		switch env, err := strconv.Atoi(tag); {
		case err != nil:
			return err
		case v.OverflowInt(int64(env)):
			return errors.New("Config Parse Error: " + tag + " Overflows Int64")
		default:
			v.Field(i).SetInt(int64(env))
		}
	}
	return nil
}

func ParseOverrides(c *Config, cliOpts map[string]string) (*Config, error) {
	p := &parser{"QAAS_", cliOpts}
	switch err := p.walkFields(reflect.ValueOf(c), ""); {
	case err != nil:
		return nil, err
	case len(p.opts) > 0:
		return nil, errors.New("Malformed/Unknown opts: " + p.unknownOpts())
	default:
		return c, nil
	}
}
