package types

import (
	"github.com/google/uuid"
)

// Quote
type Quote struct {
	ID     string   `json:"id"`
	Text   string   `json:"text" dynamoav:",omitempty"`
	Author string   `json:"author" dynamoav:",omitempty"`
	Topics []string `json:"topics" dynamoav:",stringset,omitempty"`
}

type quoteOpt func(q *Quote) *Quote

func NewQuote(opts ...quoteOpt) *Quote {
	q := &Quote{}
	for _, opt := range opts {
		q = opt(q)
	}
	return q
}


func (q *Quote) WithID() *quoteOpt {
	return func(q *Quote) *Quote {
		q.ID = uuid.New().String()
		return q
	}
}

func (q *Quote) WithText(t string) *quoteOpt {
	return func(q *Quote) *Quote {
		q.Text = t
		return q
	}
}

func (q *Quote) WithAuthor(a string) *quoteOpt {
	return func(q *Quote) *Quote {
		q.Author = a
		return q
	}
}

func (q *Quote) WithTopics(t []string) *quoteOpt {
	return func(q *Quote) *Quote {
		q.Topics = t
		return q
	}
}

// Author
type Author struct {
	Name    string `json:"name"`
	QuoteID string `json:"quote_id"`
}

type authorOpt func(a *Author) *Author

func NewAuthor(opts ...authorOpt) *Author {
	a := &Author{}
	for _, opt := range opts {
		a = opt(a)
	}
	return a
}

func (a *Author) WithQuoteID(id string) *authorOpt {
	return func(a *Author) *Author {
		a.QuoteID = id
		return a
	}
}

func (a *Author) WithName(name string) *authorOpt {
	return func(a *Author) *Author {
		a.Name = name
		return a
	}
}

// Topic
type Topic struct {
	Name    string `json:"name"`
	QuoteID string `json:"quote_id"`
}

type topicOpt func(t *Topic) *Topic

func NewTopic(opts ...topicOpt) *Topic {
	t := &Topic{}
	for _, opt := range opts {
		t = opt(t)
	}
	return t
}

func (t *Topic) WithQuoteID(id string) *topicOpt {
	return func(t *Topic) *Topic {
		t.QuoteID = id
		return t
	}
}

func (t *Topic) WithName(name string) *topicOpt {
	return func(t *Topic) *Topic {
		t.Name = name
		return t
	}
}
