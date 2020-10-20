package clean

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/travmatth-org/qaas/internal/logger"
	"github.com/travmatth-org/qaas/internal/types"

	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/validator.v2"
)

type cleaner struct {
	once      sync.Once
	validator *validator.Validator
	sanitizer *bluemonday.Policy
}

const (
	matchAlNumString = "^[a-zA-Z0-9 ]{3,20}$"
)

var (
	regex     = regexp.MustCompile(matchAlNumString)
	singleton = cleaner{}
)

func init() {
	singleton.sanitizer = bluemonday.NewPolicy()
	singleton.validator = validator.NewValidator()
	_ = singleton.validator.SetValidationFunc("topics", singleton.validateTopics)
}

func (c *cleaner) validateTopics(t interface{}, param string) error {
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

// sanitize cleans untrusted user input according to OWASP Go Language Guide
// https://owasp.org/www-project-go-secure-coding-practices-guide/
func (c *cleaner) sanitize(raw string) string {
	// Strip all tags, convert single less-than characters < to entity
	new := c.sanitizer.Sanitize(raw)
	// remove line breaks, tabs and extra white space
	return strings.Join(strings.Fields(new), " ")
}

// Quote validates to ensure struct passed is compliant with `validate` tag,
// sanitizes quote to strip html tags and html encode html entitities
func Quote(quote *types.Quote) error {
	if err := singleton.validator.Validate(quote); err != nil {
		return err
	} else if clean := singleton.sanitize(quote.Text); clean == "" {
		logger.Warn().Str("text", quote.Text).Msg("No text after sanitization")
		return errors.New("Error: Quote cannot be empty")
	} else {
		quote.Text = clean
	}
	return nil
}

// Query decodes and validates the client batch get request
func Query(get *types.GetBatch) error {
	return singleton.validator.Validate(get)
}
