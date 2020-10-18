package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const (
	envTag = "env"
	cliTag = "cli"
)

// parser is a helper struct for parsing configuration struct and options
type parser struct {
	prefix string
	opts   map[string]string
}

// getTagVal gets tag value from the struct field and returns cli or env val
func (p *parser) getTagVal(v reflect.StructField) string {
	if tag, ok := p.opts[v.Tag.Get(cliTag)]; ok {
		delete(p.opts, v.Tag.Get(cliTag))
		return tag
	} else if tag = v.Tag.Get(envTag); tag == "" {
		return ""
	} else {
		return os.Getenv(p.prefix + tag)
	}
}

// walkFields walks a reflected value, recursing on structs and setting
// fields with values returned by getTagVal
func (p *parser) walkFields(v reflect.Value, tag string) error {
	switch k := v.Kind(); {
	case k == reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f, sf := v.Field(i), v.Type().Field(i)
			next := p.getTagVal(sf)
			if err := p.walkFields(f, next); err != nil {
				return err
			}
		}
	case k == reflect.String && tag != "":
		v.SetString(tag)
	case (k == reflect.Int || k == reflect.Int64) && tag != "":
		switch num, err := strconv.ParseInt(tag, 10, 64); {
		case err != nil:
			return err
		case v.OverflowInt(num):
			return fmt.Errorf("Value %s overflows int64", tag)
		default:
			v.SetInt(num)
		}
	}
	return nil
}

// ParseOverrides overrides values in the config struct with
// environment variables and values passed from os.Args through opts
func ParseOverrides(c interface{}, opts map[string]string) error {
	var (
		p   = &parser{"QAAS_", opts}
		val = reflect.ValueOf(c).Elem()
	)
	switch err := p.walkFields(val, ""); {
	case err != nil:
		return err
	case len(p.opts) > 0:
		message := ""
		for key, value := range p.opts {
			message += key + "=" + value + "\n"
		}
		return errors.New("Malformed/Unknown opts: " + message)
	default:
		return nil
	}
}
