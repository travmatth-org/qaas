package types

import (
	"github.com/google/uuid"
)

type Quote struct {
	ID     string   `json:"id"`
	Text   string   `json:"text" dynamoav:",omitempty"`
	Author string   `json:"author" dynamoav:",omitempty"`
	Topics []string `json:"topics" dynamoav:",stringset,omitempty"`
}

func NewQuote() *Quote {
	return &Quote{}
}


func (q *Quote) WithID() *Quote  {
	q.ID = uuid.New().String()
	return q
}

func (q *Quote) WithText(t string) *Quote {
	q.Text = t
	return q
}

func (q *Quote) WithAuthor(a string) *Quote {
	q.Author = a
	return q
}

func (q *Quote) WithTopics(t []string) *Quote {
	q.Topics = t
	return q
}

type Author struct {
	Name    string `json:"name"`
	QuoteID string `json:"quote_id"`
}

func NewAuthor() *Author {
	return &Author{}
}

func (a *Author) WithQuoteID(id string) *Author {
	a.QuoteID = id
	return a
}

func (a *Author) WithName(name string) *Author {
	a.Name = name
	return a
}

// Topic
type Topic struct {
	Name    string `json:"name"`
	QuoteID string `json:"quote_id"`
}

func NewTopic() *Topic {
	return &Topic{}
}

func (t *Topic) WithQuoteID(id string) *Topic {
	t.QuoteID = id
	return t
}

func (t *Topic) WithName(name string) *Topic {
	t.Name = name
	return t
}
