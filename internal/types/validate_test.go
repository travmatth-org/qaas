package types

import (
	"testing"
)

func goodQuote() *Quote {
	return NewQuote().
		WithText("This is a good quote").
		WithAuthor("Author").
		WithTopics([]string{"foo", "bar"})
}

var illegal = "{}\\/[]"
var short = "a"
var long = `
abcdefghijklmnopqrstuvwxyz
abcdefghijklmnopqrstuvwxyz
abcdefghijklmnopqrstuvwxyz
abcdefghijklmnopqrstuvwxyz`

var longArr = []string{"foo", "foo", "foo", "foo", "foo", "foo"}
var longElemArr = []string{long}
var shortElemArr = []string{short}
var illegalArr = []string{"func (int) {}", "<div>alert(\"foo\")</div>"}

func TestValidateQuote(t *testing.T) {
	for _, tt := range []struct {
		quote *Quote
		valid bool
	}{
		// good
		{goodQuote(), true},
		// id
		{goodQuote().WithID(short), false},
		{goodQuote().WithID(illegal), false},
		// text
		{goodQuote().WithText(""), false},
		{goodQuote().WithText(short), false},
		{goodQuote().WithText(long), false},
		// author
		{goodQuote().WithAuthor(""), false},
		{goodQuote().WithAuthor(short), false},
		{goodQuote().WithAuthor(long), false},
		{goodQuote().WithAuthor(illegal), false},
		// topics
		{goodQuote().WithTopics([]string{short}), false},
		{goodQuote().WithTopics([]string{long}), false},
		{goodQuote().WithTopics([]string{illegal}), false},
		{goodQuote().WithTopics(longArr), false},
		{goodQuote().WithTopics(illegalArr), false},
		{goodQuote().WithTopics(longElemArr), false},
		{goodQuote().WithTopics(shortElemArr), false},
	} {
		if err := ValidateStruct(tt.quote); err != nil && tt.valid {
			t.Fatalf("Error validating %+v: %s", tt.quote, err)
		} else if err == nil && !tt.valid {
			t.Fatalf("Error validating %+v: shouldn't validate", tt.quote)
		}
	}
}

func goodBatch() *GetBatch {
	return &GetBatch{
		Name:  "author",
		Value: "A good value",
		Start: nil,
	}
}

func (g *GetBatch) WithName(s string) *GetBatch {
	g.Name = s
	return g
}

func (g *GetBatch) WithValue(s string) *GetBatch {
	g.Value = s
	return g
}

func (g *GetBatch) WithRecord(r *Record) *GetBatch {
	g.Start = r
	return g
}

func TestValidateGetBatch(t *testing.T) {
	for _, tt := range []struct {
		get *GetBatch
		ok  bool
	}{
		// {goodBatch(), true},
		// name
		// {goodBatch().WithName("topic"), true},
		// {goodBatch().WithName(long), false},
		// val
		// {goodBatch().WithValue(short), false},
		// {goodBatch().WithValue(long), false},
		{goodBatch().WithValue(illegal), false},
		// record
		{goodBatch().WithRecord(NewRecord().WithQuoteID(short)), false},
		{goodBatch().WithRecord(NewRecord().WithQuoteID(long)), false},
		{goodBatch().WithRecord(NewRecord().WithQuoteID(illegal)), false},
		{goodBatch().WithRecord(NewRecord().WithName(short)), false},
		{goodBatch().WithRecord(NewRecord().WithName(long)), false},
		{goodBatch().WithRecord(NewRecord().WithName(illegal)), false},
	} {
		if err := ValidateStruct(tt.get); err != nil && tt.ok {
			t.Fatalf("Error validating %+v: %s", tt.get, err)
		} else if err == nil && !tt.ok {
			t.Fatalf("Error validating %+v: shouldn't validate", tt.get)
		}
	}
}
