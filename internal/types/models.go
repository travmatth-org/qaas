package types

import (
	"github.com/google/uuid"
)

// Quote models an individual quote by an author, related to a set of topics
type Quote struct {
	ID     string   `validate:"max=0"`
	Text   string   `validate:"min=5,max=100"`
	Author string   `validate:"min=5,max=20,regexp=^[a-zA-Z ]*$"`
	Topics []string `validate:"max=5,topics" dynamodbav:",stringset"`
}

// NewQuote returns a new Quote
func NewQuote() *Quote {
	return &Quote{}
}

// NewID inserts a new UUID into a Quote
func (q *Quote) NewID() *Quote {
	q.ID = uuid.New().String()
	return q
}

// WithID inserts the given ID into a Quote
func (q *Quote) WithID(id string) *Quote {
	q.ID = id
	return q
}

// WithText inserts the given text into a Quote
func (q *Quote) WithText(t string) *Quote {
	q.Text = t
	return q
}

// WithAuthor inserts the given author into a Quote
func (q *Quote) WithAuthor(a string) *Quote {
	q.Author = a
	return q
}

// WithTopics inserts the given topics into a Quote
func (q *Quote) WithTopics(t []string) *Quote {
	q.Topics = t
	return q
}

// Record models a piece of quote metadata and its relation to quotes
type Record struct {
	Name    string `validate:"min=3,max=20,regexp=^[a-zA-Z0-9 ]*$" json:"Name"`
	QuoteID string `validate:"min=3,max=20,regexp=^[a-zA-Z0-9 ]*$" json:"QuoteID"`
}

// NewRecord returns a new Record
func NewRecord() *Record {
	return &Record{}
}

// RecordsFromQuote creates a quotes record metadata from the given quote
func RecordsFromQuote(q *Quote) (*Record, []*Record) {
	var (
		author = NewRecord().WithName(q.Author).WithQuoteID(q.ID)
		topics = []*Record{}
	)
	for _, topic := range q.Topics {
		t := NewRecord().WithName(topic).WithQuoteID(q.ID)
		topics = append(topics, t)
	}
	return author, topics
}

// WithQuoteID inserts a new quote id into the Record
func (r *Record) WithQuoteID(id string) *Record {
	r.QuoteID = id
	return r
}

// WithName inserts a given name into the Record
func (r *Record) WithName(name string) *Record {
	r.Name = name
	return r
}
