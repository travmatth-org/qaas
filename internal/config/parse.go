package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const (
	ENVTag = "env"
	CLITag = "cli"
)

type parser struct {
	prefix string
	opts   map[string]string
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
	switch k := v.Kind(); {
	case k == reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			next := p.getTagVal(v.Type().Field(i))
			if err := p.walkFields(f, next); err != nil {
				return err
			}
		}
	case k == reflect.String:
		v.SetString(tag)
	case k == reflect.Int:
		switch num, err := strconv.ParseInt(tag, 10, 64); {
		case err != nil:
			return err
		case v.OverflowInt(num):
			return errors.New(fmt.Sprintf("Value %s overflows int64", tag))
		default:
			v.SetInt(num)
		}
	}
	return nil
}

func ParseOverrides(c interface{}, opts map[string]string) error {
	p := &parser{"QAAS_", opts}
	val := reflect.ValueOf(c).Elem()
	switch err := p.walkFields(val, ""); {
	case err != nil:
		return err
	case len(p.opts) > 0:
		return errors.New("Malformed/Unknown opts: " + p.unknownOpts())
	default:
		return nil
	}
}
