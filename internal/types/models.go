package types

import (
	"github.com/google/uuid"
)

// Quote models an individual quote by an author, related to a set of topics
type Quote struct {
	ID     string   `json:"id"`
	Text   string   `json:"text" dynamoav:",omitempty"`
	Author string   `json:"author" dynamoav:",omitempty"`
	Topics []string `json:"topics" dynamoav:",stringset,omitempty"`
}

// NewQuote returns a new Quote
func NewQuote() *Quote {
	return &Quote{}
}

// WithID inserts a new UUID into a Quote
func (q *Quote) WithID() *Quote {
	q.ID = uuid.New().String()
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

// Author models an individual Author and their attributed quote's ID
type Author struct {
	Name    string `json:"name"`
	QuoteID string `json:"quote_id"`
}

// NewAuthor returns a new Author
func NewAuthor() *Author {
	return &Author{}
}

// WithQuoteID inserts a new quote id into the Author
func (a *Author) WithQuoteID(id string) *Author {
	a.QuoteID = id
	return a
}

// WithName inserts a given name into the Author
func (a *Author) WithName(name string) *Author {
	a.Name = name
	return a
}

// Topic models a given topic, and the quote associated with it
type Topic struct {
	Name    string `json:"name"`
	QuoteID string `json:"quote_id"`
}

// NewTopic return a new Topic
func NewTopic() *Topic {
	return &Topic{}
}

// WithQuoteID inserts a given quote id into the Topic
func (t *Topic) WithQuoteID(id string) *Topic {
	t.QuoteID = id
	return t
}

// WithName inserts a given name into the Topic
func (t *Topic) WithName(name string) *Topic {
	t.Name = name
	return t
}
