package types

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/aws/session" 
	"github.com/google/uuid"
)

// internal/api/api.go
type AWSSession = session.Session
type AWSConfig = aws.Config

// internal/api/dynamodb/actions
type AWSDynamoDB = dynamodb.DynamoDB
type AWSQueryInput = dynamodb.QueryInput
type AWSAttrVal = dynamodb.AttributeValue
type AWSConditionBuilder = dynamodb.ConditionBuilder
type AWSAttrValMap = []map[string]*dynamodb.AttributeValue
type AWSGetItemInput = dynamodb.GetItemInput

type Quote struct {
	ID     string   `json:"id"`
	Text   string   `json:"text" dynamoav:",omitempty"`
	Author string   `json:"author" dynamoav:",omitempty"`
	Topics []string `json:"topics" dynamoav:",stringset,omitempty"`
}

func NewQuote() *Quote {
	return &Quote{}
}

func (q *Quote) WithID() *Quote {
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

type QuoteResponse struct {
	Quote *Quote
	Err   error
}

func NewQuoteResponse() *QuoteResponse {
	return &QuoteResponse{}
}

func (q *QuoteResponse) WithQuote(quote *Quote) *QuoteResponse {
	q.Quote = quote
	return q
}

func (q *QuoteResponse) WithErr(err error) *QuoteResponse {
	q.Err = err
	return q
}

type MultiQuoteResponse struct {
	Quotes []*Quote
	Next   string
	Err    error
}

func NewMultiQuoteResponse() *MultiQuoteResponse {
	return &MultiQuoteResponse{}
}

func (m *MultiQuoteResponse) WithQuotes(q []*Quote) *MultiQuoteResponse {
	m.Quotes = q
	return m
}

func (m *MultiQuoteResponse) WithNext(n string) *MultiQuoteResponse {
	m.Next = n
	return m
}

func (m *MultiQuoteResponse) WithErr(e error) *MultiQuoteResponse {
	m.Err = e
	return m
}
