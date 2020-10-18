package types

import (
	"fmt"
	"regexp"
	"sync"

	"gopkg.in/validator.v2"
)

type singleton struct {
	once      sync.Once
	validator *validator.Validator
}

var single = &singleton{}

var regex = regexp.MustCompile("^[a-zA-Z0-9 ]{3,20}$")

func (s *singleton) validateTopics(t interface{}, param string) error {
	switch t := t.(type) {
	case []string:
		for _, val := range t {
			if len(val) < 3 || len(val) > 20 {
				return fmt.Errorf("Topics cannot exceed 20 characters, have %d", len(val))
			} else if !regex.MatchString(val) {
				return fmt.Errorf("Invalid topic: %s", val)
			}
		}
		return nil
	default:
		return fmt.Errorf("Invalid type for topics: %s", t)
	}
}

func (s *singleton) performValidation(i interface{}) error {
	if s.validator == nil {
		s.once.Do(func() {
			s.validator = validator.NewValidator()
			_ = s.validator.SetValidationFunc("topics", s.validateTopics)
		})
	}
	return s.validator.Validate(i)
}

// ValidateStruct ensures struct passed is compliant with `validate` tag
func ValidateStruct(v interface{}) error {
	return single.performValidation(v)
}
