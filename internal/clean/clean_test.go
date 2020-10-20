package clean

import (
	"testing"

	"github.com/travmatth-org/qaas/internal/types"
	confighelpers "github.com/travmatth-org/qaas/test"
)

const (
	xss     = "Hello <STYLE>.XSS{background-image:url(\"javascript:alert('XSS')\");}</STYLE><A CLASS=XSS></A>World"
	illegal = "{}\\/[]"
	short   = "a"
	long    = `
abcdefghijklmnopqrstuvwxyz
abcdefghijklmnopqrstuvwxyz
abcdefghijklmnopqrstuvwxyz
abcdefghijklmnopqrstuvwxyz`
)

var (
	longArr      = []string{"foo", "foo", "foo", "foo", "foo", "foo"}
	longElemArr  = []string{long}
	shortElemArr = []string{short}
	illegalArr   = []string{"func (int) {}", "<div>alert(\"foo\")</div>"}
)

func goodQuote(text, author string) *types.Quote {
	if text == "" {
		text = "This is a good quote"
	}
	if author == "" {
		author = "Author"
	}
	return types.NewQuote().
		WithText(text).
		WithAuthor(author).
		WithTopics([]string{"foo", "bar"})
}

func TestValidateQuote(t *testing.T) {
	_ = confighelpers.ResetLogger()
	for _, tt := range []struct {
		quote *types.Quote
		valid bool
	}{
		// good
		{goodQuote("", ""), true},
		// id
		{goodQuote("", "").WithID(short), false},
		{goodQuote("", "").WithID(illegal), false},
		// text
		{goodQuote("", "").WithText(""), false},
		{goodQuote("<div></div>", ""), false},
		{goodQuote(short, ""), false},
		{goodQuote(long, ""), false},
		// author
		{goodQuote("", "").WithAuthor(""), false},
		{goodQuote("", short), false},
		{goodQuote("", long), false},
		{goodQuote("", illegal), false},
		// topics
		{goodQuote("", "").WithTopics([]string{short}), false},
		{goodQuote("", "").WithTopics([]string{long}), false},
		{goodQuote("", "").WithTopics([]string{illegal}), false},
		{goodQuote("", "").WithTopics(longArr), false},
		{goodQuote("", "").WithTopics(illegalArr), false},
		{goodQuote("", "").WithTopics(longElemArr), false},
		{goodQuote("", "").WithTopics(shortElemArr), false},
	} {
		if err := Quote(tt.quote); err != nil && tt.valid {
			t.Fatalf("Error validating %+v: %s", tt.quote, err)
		} else if err == nil && !tt.valid {
			t.Fatalf("Error validating %+v: shouldn't validate", tt.quote)
		}
	}
}

func TestSanitize(t *testing.T) {
	for _, tt := range []struct {
		quote    *types.Quote
		expected string
	}{
		{goodQuote("<div></div>", ""), ""},
		{goodQuote(xss, ""), "Hello World"},
		{goodQuote("<div></div>", ""), ""},
		{goodQuote("should encode < >", ""), "should encode &lt; &gt;"},
		{goodQuote("foo    \n    \t     foo", ""), "foo foo"},
	} {
		if got := singleton.sanitize(tt.quote.Text); got != tt.expected {
			t.Fatalf("Error validating string: %s != %s", got, tt.expected)
		}
	}
}

func newBatch(name, val string, start *types.Record) *types.GetBatch {
	if name == "" {
		name = "author"
	}
	if val == "" {
		val = "A good value"
	}
	return &types.GetBatch{
		Name:  name,
		Value: val,
		Start: start,
	}
}

func TestValidateGetBatch(t *testing.T) {
	for _, tt := range []struct {
		get *types.GetBatch
		ok  bool
	}{
		{newBatch("", "", nil), true},
		// name
		{newBatch("topic", "", nil), true},
		{newBatch(long, "", nil), false},
		// val
		{newBatch("", short, nil), false},
		{newBatch("", long, nil), false},
		{newBatch("", illegal, nil), false},
		// record
		{newBatch("", "", types.NewRecord().WithQuoteID(short)), false},
		{newBatch("", "", types.NewRecord().WithQuoteID(long)), false},
		{newBatch("", "", types.NewRecord().WithQuoteID(illegal)), false},
		{newBatch("", "", types.NewRecord().WithName(short)), false},
		{newBatch("", "", types.NewRecord().WithName(long)), false},
		{newBatch("", "", types.NewRecord().WithName(illegal)), false},
	} {
		if err := Query(tt.get); err != nil && tt.ok {
			t.Fatalf("Error validating %+v: %s", tt.get, err)
		} else if err == nil && !tt.ok {
			t.Fatalf("Error validating %+v: shouldn't validate", tt.get)
		}
	}
}
